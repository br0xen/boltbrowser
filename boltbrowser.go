package main

import (
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

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filename(s)>\n", ProgramName)
		os.Exit(1)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	style := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)

	databaseFiles := os.Args[1:]
	for _, databaseFile := range databaseFiles {
		currentFilename = databaseFile
		db, err = bolt.Open(databaseFile, 0600, nil)
		if err != nil {
			mainLoop(nil, style)
			continue
		}

		// First things first, load the database into memory
		memBolt.refreshDatabase()
		mainLoop(memBolt, style)
		defer db.Close()
	}
}
