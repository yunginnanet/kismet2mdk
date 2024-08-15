package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"

	"github.com/l0nax/go-spew/spew"
)

const queryfmt = `SELECT DISTINCT * FROM (SELECT sourcemac FROM packets WHERE destmac = ? UNION SELECT destmac FROM packets WHERE sourcemac = ?);`

func (kdb *KismetDatabase) FindRelatedMacs(mac string) ([]string, error) {
	return kdb.FindRelatedMacsCtx(context.Background(), mac)
}

func debugRows(rows *sql.Rows) {
	cols, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	spew.Dump(cols)
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		panic(err.Error())
	}
	spew.Dump(colTypes)

}

func (kdb *KismetDatabase) FindRelatedMacsCtx(ctx context.Context, mac string) ([]string, error) {
	if _, err := net.ParseMAC(mac); err != nil {
		return nil, err
	}

	//goland:noinspection SqlResolve
	rows, err := kdb.conn.QueryContext(ctx, queryfmt, mac, mac)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var related = make([]string, 0)

	var parseErrs []error

	debugRows(rows)

	for rows.Next() {

		if err = rows.Err(); err != nil {
			return nil, err
		}
		var addr string
		if err = rows.Scan(&addr); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if addr == "" {
			parseErrs = append(parseErrs, errors.New("blank value was present"))
			continue
		}
		if _, validErr := net.ParseMAC(addr); validErr != nil {
			parseErrs = append(parseErrs, fmt.Errorf("ignored seemingly invalid mac: %s", addr))
			continue
		}
		println(addr)
		related = append(related, addr)
	}

	return related, errors.Join(parseErrs...)
}
