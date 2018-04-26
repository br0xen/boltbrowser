package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
)

var ProgramName = "boltbrowser"

var databaseFiles []string
var db *bolt.DB
var memBolt *BoltDB

var currentFilename string

const DefaultDBOpenTimeout = time.Second

var AppArgs struct {
	DBOpenTimeout time.Duration
	ReadOnly      bool
}

func init() {
	AppArgs.DBOpenTimeout = DefaultDBOpenTimeout
	AppArgs.ReadOnly = false
}

func parseArgs() {
	var err error
	if len(os.Args) == 1 {
		printUsage(nil)
	}
	parms := os.Args[1:]
	for i := range parms {
		// All 'option' arguments start with "-"
		if !strings.HasPrefix(parms[i], "-") {
			databaseFiles = append(databaseFiles, parms[i])
			continue
		}
		if strings.Contains(parms[i], "=") {
			// Key/Value pair Arguments
			pts := strings.Split(parms[i], "=")
			key, val := pts[0], pts[1]
			switch key {
			case "-timeout":
				AppArgs.DBOpenTimeout, err = time.ParseDuration(val)
				if err != nil {
					// See if we can successfully parse by adding a 's'
					AppArgs.DBOpenTimeout, err = time.ParseDuration(val + "s")
				}
				// If err is still not nil, print usage
				if err != nil {
					printUsage(err)
				}
			case "-readonly", "-ro":
				if val == "true" {
					AppArgs.ReadOnly = true
				}
			case "-help":
				printUsage(nil)
			default:
				printUsage(errors.New("Invalid option"))
			}
		} else {
			// Single-word arguments
			switch parms[i] {
			case "-readonly", "-ro":
				AppArgs.ReadOnly = true
			case "-help":
				printUsage(nil)
			default:
				printUsage(errors.New("Invalid option"))
			}
		}
	}
}

func printUsage(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <filename(s)>\nOptions:\n", ProgramName)
	fmt.Fprintf(os.Stderr, "  -timeout=duration\n        DB file open timeout (default 1s)\n")
	fmt.Fprintf(os.Stderr, "  -ro, -readonly   \n        Open the DB in read-only mode\n")
}

func main() {
	var err error

	parseArgs()

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	style := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)

	for _, databaseFile := range databaseFiles {
		currentFilename = databaseFile
		db, err = bolt.Open(databaseFile, 0600, &bolt.Options{Timeout: AppArgs.DBOpenTimeout})
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
				os.Exit(1)
			}
		}

		// First things first, load the database into memory
		memBolt.refreshDatabase()
		if AppArgs.ReadOnly {
			// If we're opening it in readonly mode, close it now
			db.Close()
		}

		// Kick off the UI loop
		mainLoop(memBolt, style)
		defer db.Close()
	}
}
