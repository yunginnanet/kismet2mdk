package data

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
	pragma journal_size_limit = ;
*/

type Pragma string

const (
	PragmaJournalMode      Pragma = "journal_mode"
	PragmaSynchronous      Pragma = "synchronous"
	PragmaJournalSizeLimit Pragma = "journal_size_limit"
)

func (p Pragma) String() string { return string(p) }

func (p Pragma) SetQuery(v string) string {
	return fmt.Sprintf("PRAGMA %s = %s;", p, v)
}

func (kdb *KismetDatabase) RestorePragma(s Pragma, butnot ...string) error {
	var old string
	if old = kdb.checkBackupPragma(s); old == "" {
		return errors.New("no backup pragma found")
	}
	if len(butnot) > 0 {
		for _, not := range butnot {
			if strings.EqualFold(old, not) {
				return fmt.Errorf("stored pragma conflict: %s == %s", s, not)
			}
		}
	}
	_, err := kdb.conn.Exec(s.SetQuery(old))
	if err == nil {
		kdb.clearPragmaBackup(s)
	}
	return err
}

func (kdb *KismetDatabase) EnableWAL(b bool) error {
	if b {
		if stored := kdb.checkBackupPragma(PragmaJournalMode); stored != "" && stored != "WAL" {
			return errors.New("WAL mode (should be) enabled already")
		}
		if err := kdb.backupPragma(PragmaJournalMode); err != nil {
			return err
		}
		_, err := kdb.conn.Exec(PragmaJournalMode.SetQuery("WAL"))
		return err
	}
	var err error
	if err = kdb.RestorePragma(PragmaJournalMode, "WAL"); err != nil {
		_, err = kdb.conn.Exec(PragmaJournalMode.SetQuery("DELETE"))
	}
	return err
}

func (kdb *KismetDatabase) EnableAsync(b bool) error {
	if b {
		if stored := kdb.checkBackupPragma(PragmaSynchronous); stored != "" && stored != "OFF" {
			return errors.New("async mode (should be) enabled already")
		}
		if err := kdb.backupPragma(PragmaSynchronous); err != nil {
			return err
		}
		_, err := kdb.conn.Exec(PragmaSynchronous.SetQuery("OFF"))
		return err
	}
	var err error
	if err = kdb.RestorePragma(PragmaSynchronous, "OFF"); err != nil {
		_, err = kdb.conn.Exec(PragmaSynchronous.SetQuery("NORMAL"))
	}
	return err
}

func (kdb *KismetDatabase) JournalSizeLimit(size int64) error {
	if err := kdb.backupPragma(PragmaJournalSizeLimit); err != nil {
		return err
	}
	if _, err := kdb.conn.Exec(PragmaJournalSizeLimit.SetQuery(strconv.Itoa(int(size)))); err != nil {
		return fmt.Errorf("failed to set journal size limit: %w", err)
	}
	return nil
}

func (kdb *KismetDatabase) SetTmpDir(path string) {
	_ = os.MkdirAll(path, 0755)
	_, tmpDirErr := kdb.conn.Exec("PRAGMA temp_store_directory = '" + path + "';")
	if tmpDirErr != nil {
		println("WARN: unable to set tmp dir", tmpDirErr.Error())
		return
	}
	kdb.newTmpDir = path
}
