package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "github.com/glebarez/go-sqlite"
)

//goland:noinspection SqlNoDataSourceInspection
const kismetSchema = `
CREATE TABLE IF NOT EXISTS KISMET (kismet_version TEXT, db_version INT, db_module TEXT);
CREATE TABLE IF NOT EXISTS devices (first_time INT, last_time INT, devkey TEXT, phyname TEXT, devmac TEXT, strongest_signal INT, min_lat REAL, min_lon REAL, max_lat REAL, max_lon REAL, avg_lat REAL, avg_lon REAL, bytes_data INT, type TEXT, device BLOB, UNIQUE(phyname, devmac) ON CONFLICT REPLACE);
CREATE TABLE IF NOT EXISTS packets (ts_sec INT, ts_usec INT, phyname TEXT, sourcemac TEXT, destmac TEXT, transmac TEXT, frequency REAL, devkey TEXT, lat REAL, lon REAL, alt REAL, speed REAL, heading REAL, packet_len INT, signal INT, datasource TEXT, dlt INT, packet BLOB, error INT, tags TEXT, datarate REAL, hash INT, packetid INT );
CREATE TABLE IF NOT EXISTS data (ts_sec INT, ts_usec INT, phyname TEXT, devmac TEXT, lat REAL, lon REAL, alt REAL, speed REAL, heading REAL, datasource TEXT, type TEXT, json BLOB );
CREATE TABLE IF NOT EXISTS datasources (uuid TEXT, typestring TEXT, definition TEXT, name TEXT, interface TEXT, json BLOB, UNIQUE(uuid) ON CONFLICT REPLACE);
CREATE TABLE IF NOT EXISTS alerts (ts_sec INT, ts_usec INT, phyname TEXT, devmac TEXT, lat REAL, lon REAL, header TEXT, json BLOB );
CREATE TABLE IF NOT EXISTS messages (ts_sec INT, lat REAL, lon REAL, msgtype TEXT, message TEXT );
CREATE TABLE IF NOT EXISTS snapshots (ts_sec INT, ts_usec INT, lat REAL, lon REAL, snaptype TEXT, json BLOB );
`

// intentionally partial list
var kismetTables = []string{"KISMET", "devices", "packets", "data", "alerts", "messages"}

//goland:noinspection SqlNoDataSourceInspection
func tableExistsQuery(name string) string {
	return fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", name)
}

type KismetDatabase struct {
	path string
	conn *sql.DB

	pragma map[Pragma]string
	mu     sync.Mutex

	newTmpDir string
}

func (kdb *KismetDatabase) backupPragma(s Pragma) error {
	var pragma string
	if err := kdb.conn.QueryRow("PRAGMA " + string(s)).Scan(&pragma); err != nil {
		return fmt.Errorf("failed to backup pragma %s: %w", s, err)
	}
	kdb.mu.Lock()
	kdb.pragma[s] = pragma
	kdb.mu.Unlock()
	return nil
}

func (kdb *KismetDatabase) checkBackupPragma(s Pragma) string {
	kdb.mu.Lock()
	r, ok := kdb.pragma[s]
	kdb.mu.Unlock()
	if !ok {
		return ""
	}
	return r
}

func (kdb *KismetDatabase) clearPragmaBackup(s Pragma) {
	kdb.mu.Lock()
	delete(kdb.pragma, s)
	kdb.mu.Unlock()
}

func CheckKismetSchema(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}

	var tErrs = make([]error, 0, len(kismetTables))

	for _, t := range kismetTables {
		if _, err := db.Exec(tableExistsQuery(t)); err != nil {
			tErrs = append(tErrs, fmt.Errorf("missing table %s in kismet database: %w", t, err))
		}
	}

	return errors.Join(tErrs...)
}

func OpenKismetDatabase(path string) (*KismetDatabase, error) {
	stat, err := os.Stat(path)

	switch {
	case err == nil, errors.Is(err, os.ErrNotExist):
		break
	case errors.Is(err, os.ErrPermission):
		return nil, fmt.Errorf("fs perms on target db: %w", err)
	case stat.IsDir():
		return nil, fmt.Errorf("fs: %s is a directory", path)
	default:
		return nil, fmt.Errorf("db path stat err: %w", err)
	}

	kdb := new(KismetDatabase)
	kdb.pragma = make(map[Pragma]string)
	kdb.path = path
	if kdb.conn, err = sql.Open("sqlite", path); err != nil {
		return nil, fmt.Errorf("sql: %w", err)
	}

	if err = kdb.conn.Ping(); err != nil {
		return nil, fmt.Errorf("sql ping: %w", err)
	}

	res, err := kdb.conn.Exec(kismetSchema)
	if err != nil {
		return nil, fmt.Errorf("sql, failed to assure schema: %w", err)
	}

	if affected, _ := res.RowsAffected(); affected > 0 {
		println("wrote schema to db at", path)
	}

	return kdb, nil
}

func (kdb *KismetDatabase) tables() (*sql.Rows, error) {
	return kdb.conn.Query("SELECT name FROM sqlite_master WHERE type='table'")
}

func (kdb *KismetDatabase) Vacuum() error {
	_, err := kdb.conn.Exec("VACUUM")
	return err
}

func (kdb *KismetDatabase) Analyze() error {
	_, err := kdb.conn.Exec("ANALYZE")
	return err
}

func (kdb *KismetDatabase) Close() error {
	if kdb.newTmpDir != "" {
		defer func(td string) {
			_ = os.RemoveAll(td)
		}(kdb.newTmpDir)
	}
	return kdb.conn.Close()
}
