package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type table struct {
	Type      sql.NullString `sql:"type"`
	Name      sql.NullString `sql:"name"`
	TableName sql.NullString `sql:"tbl_name"`
	RootPage  sql.NullInt64  `sql:"rootpage"`
}

func attachQuery(file, name string) string {
	return "ATTACH" + " '" + file + "' AS " + name + ";"
}

func detachQuery(alias string) string {
	return "DETACH" + " '" + alias + "';"
}

func gatherSources(sources ...string) ([][]string, error) {
	groupedSources := make([][]string, len(sources)/10)

	groupIndex := 0
	innerIndex := 0

	for _, source := range sources {
		if groupedSources[groupIndex] == nil {
			groupedSources[groupIndex] = make([]string, 0, 10)
		}
		groupedSources[groupIndex] = append(groupedSources[groupIndex], source)
		innerIndex++
		if innerIndex == 10 {
			groupIndex++
			innerIndex = 0
		}
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
		t := new(table)

		if err = rowsOfTables.Scan(&t); err != nil {
			return nil, fmt.Errorf("failed scanning sqlite table: %w", err)
		}
		tables = append(tables, t.Name.String)
	}

	return tables, nil
}

func ingestTenSources(tx *sql.Tx, tableNames []string, sources []string) error {
	if len(sources) > 10 {
		return errors.New("oversized sources slice")
	}

	var (
		runItUp = make(chan string, len(sources))
		errs    = make(chan error, len(sources))
		wg      sync.WaitGroup
		doneCh  = make(chan struct{}, 1)
	)

	for i, source := range sources {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := checkSource(source); err != nil {
				errs <- err
			}
			println("attaching " + source + "...")
			sourceAlias := "db" + strconv.Itoa(i)
			if _, err := tx.Exec(attachQuery(source, sourceAlias)); err != nil {
				errs <- fmt.Errorf("failed to attach %s: %w", source, err)
			}
		}()
	}

	wg2 := sync.WaitGroup{}

	go func() {
		wg.Wait()
		close(runItUp)
		wg2.Wait()
		close(doneCh)
	}()

	for inc := range runItUp {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for _, t := range tableNames {
				wg2.Add(1)
				go func() {
					defer wg2.Done()
					println("inserting " + inc + "...")
					if err := attachedInsert(tx, t, inc); err != nil {
						errs <- err
					}
				}()
			}
			println("detaching " + inc + "...")
			if _, err := tx.Exec(detachQuery(inc)); err != nil {
				errs <- err
			}
		}()
	}

snoozin:
	select {
	case <-doneCh:
		return nil
	case e := <-errs:
		if e == nil {
			goto snoozin
		}
		return e
	}
}

func attachedInsert(tx *sql.Tx, table, alias string) error {
	_, err := tx.Query("INSERT OR IGNORE INTO ? SELECT * FROM ?.?;",
		table, alias, table,
	)
	return err
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

	var tx *sql.Tx

	if tx, err = target.conn.Begin(); err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	for _, group := range grouped {
		if err = ingestTenSources(tx, tableNames, group); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("sql transaction failed: %w", err)
	}

	return nil
}
