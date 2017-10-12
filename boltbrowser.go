package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
)

/*
ProgramName is the name of the program
*/
var ProgramName = "boltbrowser"

var databaseFiles []string
var db *bolt.DB
var memBolt *BoltDB

var currentFilename string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Usage: %s <filename(s)>\n", ProgramName)
	}
}

func main() {
	var err error

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	style := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)

	databaseFiles := flag.Args()
	for _, databaseFile := range databaseFiles {
		currentFilename = databaseFile
		db, err = bolt.Open(databaseFile, 0600, nil)
		if err != nil {
			if len(databaseFiles) > 1 {
				mainLoop(nil, style)
				continue
			} else {
				termbox.Close()
				fmt.Printf("Error reading file: %q\n", err.Error())
				os.Exit(111)
			}
		}

		// First things first, load the database into memory
		memBolt.refreshDatabase()
		// Kick off the UI loop
		mainLoop(memBolt, style)
		defer db.Close()
	}
}
