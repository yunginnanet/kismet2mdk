package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"git.tcp.direct/kayos/common/entropy"
)

func attachQuery(file, name string) string {
	return "ATTACH" + " '" + file + "' AS " + name + ";"
}

func detachQuery(alias string) string {
	return "DETACH" + " '" + alias + "';"
}

func gatherSources(sources ...string) ([][]string, error) {
	groupedSources := make([][]string, 0, (len(sources)/10)+1)

	groupIndex := 0
	innerIndex := 0

	for _, source := range sources {
		if innerIndex == 10 {
			groupIndex++
			innerIndex = 0
		}
		if len(groupedSources) <= groupIndex {
			groupedSources = append(groupedSources, make([]string, 0, 10))
		}
		groupedSources[groupIndex] = append(groupedSources[groupIndex], source)
		innerIndex++
	}

	return groupedSources, nil
}

func checkSource(source string) error {
	c, err := OpenKismetDatabase(source)
	if err != nil {
		return fmt.Errorf("failed to open kismet database %s: %w", source, err)
	}

	if err = CheckKismetSchema(c.conn); err != nil {
		return fmt.Errorf("%s does not appear to be a valid kismet database: %w", source, err)
	}

	if err = c.Close(); err != nil {
		return fmt.Errorf("failed to close kismet database %s: %w", source, err)
	}

	return nil
}

func (kdb *KismetDatabase) Tables() ([]string, error) {
	var tables = make([]string, 0)
	rowsOfTables, err := kdb.tables()
	if err != nil {
		return nil, fmt.Errorf("failed getting target DB tables: %w", err)
	}

	for rowsOfTables.Next() {
		if errors.Is(rowsOfTables.Err(), sql.ErrNoRows) {
			return nil, fmt.Errorf("failed getting tables, no rows: %w", err)
		}
		var t string
		if err = rowsOfTables.Scan(&t); err != nil {
			return nil, fmt.Errorf("failed scanning sqlite table: %w", err)
		}
		tables = append(tables, t)
	}

	return tables, nil
}

type mergeTx struct {
	alias string
	tx    *sql.Tx
}

func newMergeTx(source, alias string, target *KismetDatabase) (*mergeTx, error) {
	tx, err := target.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	if _, err = tx.Exec(attachQuery(source, alias)); err != nil {
		return nil, fmt.Errorf("failed to attach %s: %w", source, err)
	}
	return &mergeTx{alias: alias, tx: tx}, nil
}

func ingestTenSources(target *KismetDatabase, tableNames []string, sources []string) error {
	if len(sources) > 10 {
		return errors.New("oversized sources slice")
	}

	var (
		merges  = make(chan *mergeTx, len(sources))
		errs    = make(chan error, len(sources))
		wg      sync.WaitGroup
		doneCh1 = make(chan struct{}, 1)
	)

	for i, source := range sources {
		wg.Add(1)
		go func(ii int) {
			defer wg.Done()
			if err := checkSource(source); err != nil {
				errs <- err
			}
			println("attaching " + source + "...")
			sourceAlias := "db" + strconv.Itoa(ii)
			tx, err := newMergeTx(source, sourceAlias, target)
			if err != nil {
				errs <- err
			}
			merges <- tx
		}(i)
	}

	// wg2 := sync.WaitGroup{}

	go func() {
		wg.Wait()
		close(merges)
		// wg2.Wait()
		close(doneCh1)
	}()

	var waitExp = &atomic.Int64{}
	waitExp.Add(1)

	for incomingMergeTx := range merges {
		// wg2.Add(1)

		mtx := incomingMergeTx

		// go func(mtx *mergeTx) {
		wg3 := sync.WaitGroup{}

		for _, t := range tableNames {
			wg3.Add(1)
			go func(mtxInner *mergeTx, tn string) {
				defer wg3.Done()

				if mtxInner == nil {
					return
				}

				println("inserting values from", mtxInner.alias, "for table", tn)

				var err = &SQLiteError{}

				ins := func() *SQLiteError {
					return mtxInner.attachedInsert(tn, mtxInner.alias)
				}

				tries := 0

				//goland:noinspection GoDirectComparisonOfErrors
				for {
					err = ins()
					if (err != nil && !err.IsBusy()) || err == nil {
						break
					}
					tries++
					if tries%2 == 1 {
						waitExp.Add(1)
					}
					println(mtxInner.alias + "." + tn + ": database busy (" + strconv.Itoa(int(waitExp.Load())) + "), waiting...")
					entropy.RandSleepMS(100 * int(waitExp.Load()))
				}
				if err != nil {
					errs <- fmt.Errorf("failed to insert values from %s for table %s: %w", mtxInner.alias, tn, err.e)
				}

			}(mtx, t)
		}

		wg3.Wait()

		if _, err := mtx.tx.Prepare(detachQuery(mtx.alias)); err != nil {
			errs <- fmt.Errorf("failed to detatch %s: %w", mtx.alias, err)
		}

		println("committing " + mtx.alias + "...")
		if err := mtx.tx.Commit(); err != nil {
			errs <- fmt.Errorf("failed to commit transaction: %w", err)
		}

		// wg2.Done()

		// }(incomingMergeTx)
	}

snoozin:
	select {
	case <-doneCh1:
		return nil
	case e := <-errs:
		if e == nil {
			goto snoozin
		}
		return e
	}
}

func (mtx *mergeTx) attachedInsert(table, alias string) *SQLiteError {
	if table == "" {
		return NewSQLiteError(errors.New("blank table during attempted merge from " + alias))
	}
	if alias == "" {
		return NewSQLiteError(errors.New("blank alias during attempted merge of " + table))
	}

	_, err := mtx.tx.Exec(fmt.Sprintf("INSERT OR IGNORE INTO '%s' SELECT * FROM %s.%s", table, alias, table))

	return NewSQLiteError(err)
}

func tidyUp(target *KismetDatabase) error {
	var err error

	print("\nvacuuming...")
	if err = target.Vacuum(); err != nil {
		return fmt.Errorf("failed to vacuum during merge: %w", err)
	}
	print("done\n")

	print("analyzing...")
	if err = target.Analyze(); err != nil {
		return fmt.Errorf("failed to analyze during merge: %w", err)
	}
	print("done\n\n")

	return nil
}

func MergeKismetDatabases(target *KismetDatabase, sources ...string) error {
	grouped, err := gatherSources(sources...)
	if err != nil {
		return err
	}

	tableNames, err := target.Tables()
	if err != nil {
		return err
	}

	for _, group := range grouped {
		if err = ingestTenSources(target, tableNames, group); err != nil {
			return err
		}
		if err = tidyUp(target); err != nil {
			return err
		}
	}

	println("committing...")

	return nil
}
