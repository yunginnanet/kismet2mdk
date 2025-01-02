package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"git.tcp.direct/kayos/kismet2mdk/pkg/data"
)

func optimize(targetDB *data.KismetDatabase) error {
	var err error

	println("enabling WAL...")
	if err = targetDB.EnableWAL(true); err != nil {
		return err
	}

	println("enabling async...")
	if err = targetDB.EnableAsync(true); err != nil {
		return err
	}

	println("setting journal size limit...")
	err = targetDB.JournalSizeLimit(6144000)

	return err
}

func restorePragma(targetDB *data.KismetDatabase) error {
	var errs = make([]error, 0, 3)

	for _, p := range []data.Pragma{data.PragmaJournalMode, data.PragmaSynchronous, data.PragmaJournalSizeLimit} {
		if err := targetDB.RestorePragma(p); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func main() {
	var target string
	var sources = make([]string, 0, len(os.Args[1:])-1)
	for i, arg := range os.Args[1:] {
		_, err := os.Stat(arg)
		if errors.Is(err, os.ErrNotExist) && i == 0 {
			var f *os.File
			f, err = os.Create(arg)
			if f != nil {
				_ = f.Close()
			}
		}
		if err != nil {
			println("kismet db access failure: ", err.Error())
			os.Exit(1)
		}
		if i == 0 {
			target = arg
			continue
		}
		sources = append(sources, arg)
	}

	targetDB, err := data.OpenKismetDatabase(target)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}

	if err = optimize(targetDB); err != nil {
		print(err.Error())
		os.Exit(1)
	}

	if cwd, _ := os.Getwd(); cwd != "" {
		targetDB.SetTmpDir(filepath.Join(cwd, ".sqlite_tmp"))
	}

	defer func() {
		if err = restorePragma(targetDB); err != nil {
			println(err.Error())
		} else {
			println("pragma restored")
		}

		println("closing " + targetDB.String())

		for err = targetDB.Close(); err != nil; err = targetDB.Close() {
			println(targetDB.String() + ": " + err.Error())
			time.Sleep(1 * time.Second)
		}
		println("db closed")

		println("fin.")

		os.Exit(0)
	}()

	if err = data.MergeKismetDatabases(targetDB, sources...); err != nil {
		print(err.Error())
		os.Exit(1)
	}
}
