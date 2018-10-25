package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/br0xen/boltbrowser/boltbrowser"
)

var ProgramName = "boltbrowser"

var databaseFiles []string

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
	parseArgs()
	for _, databaseFile := range databaseFiles {
		db, err := bolt.Open(databaseFile, 0600, &bolt.Options{Timeout: AppArgs.DBOpenTimeout})
		if err == bolt.ErrTimeout {
			fmt.Printf("File %s is locked. Make sure it's not used by another app and try again\n", databaseFile)
			os.Exit(1)
		} else if err != nil {
			fmt.Printf("Error reading file: %q\n", err.Error())
			os.Exit(1)
		}
		if AppArgs.ReadOnly {
			db.Close()
		}
		boltbrowser.Browse(db, AppArgs.ReadOnly)
		db.Close()
	}
}
