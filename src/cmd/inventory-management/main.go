package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ucok-man/go-tui-inventory-management/src/internal/data"
	"github.com/ucok-man/go-tui-inventory-management/src/internal/tui"
)

func main() {
	const defaultfile = "inventory/db.json"
	dbfile := flag.String("file", defaultfile, "Path to the storage file in json")
	flag.Parse()

	if *dbfile == "" || *dbfile != defaultfile {
		extension := filepath.Ext(*dbfile)
		if extension != "json" {
			fmt.Fprintf(flag.CommandLine.Output(), "error: file extension must be in json format")
			flag.Usage()
			os.Exit(1)
		}
	}

	model, err := data.NewInventoryModel(*dbfile)
	if err != nil {
		log.Fatal(err)
	}
	app := tui.NewTUI(model)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
