package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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

const DefaultDBOpenTimeout = time.Second

var args struct {
	DBOpenTimeout time.Duration
}

func init() {
	flag.DurationVar(&args.DBOpenTimeout, "timeout", DefaultDBOpenTimeout, "DB file open timeout")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <filename(s)>\nOptions:\n", ProgramName)
		flag.PrintDefaults()
	}
}

func main() {
	var err error

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
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
		db, err = bolt.Open(databaseFile, 0600, &bolt.Options{Timeout: args.DBOpenTimeout})
		if err == bolt.ErrTimeout {
			termbox.Close()
			fmt.Printf("File %s is locked. Make sure it's not used by another app and try again\n", databaseFile)
			os.Exit(1)
		} else if err != nil {
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
