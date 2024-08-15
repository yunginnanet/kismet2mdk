package main

import (
	"errors"
	"os"

	"git.tcp.direct/kayos/kismet2mdk/pkg/data"
)

// wat is a dee oth cy purr ess?///?///
// how does computer work??/////////////

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

	if err = data.MergeKismetDatabases(targetDB, sources...); err != nil {
		print(err.Error())
		os.Exit(1)
	}

	println("fin.")
	os.Exit(0)
}
