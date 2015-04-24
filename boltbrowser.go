package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
	"os"
	"syscall"
)

type BoltBucket struct {
	name    string
	pairs   map[string]string
	buckets []BoltBucket
}

// A Database, basically a collection of buckets
type BoltDB struct {
	buckets []BoltBucket
}

const PROGRAM_NAME = "boltbrowser"

var database_file string
var db *bolt.DB
var memBolt *BoltDB

func mainLoop(memBolt *BoltDB, style Style) {
	screens := defaultScreensForData(memBolt)
	display_screen := screens[ABOUT_SCREEN_INDEX]
	layoutAndDrawScreen(display_screen, style)
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyCtrlZ {
				process, _ := os.FindProcess(os.Getpid())
				termbox.Close()
				process.Signal(syscall.SIGSTOP)
				termbox.Init()
			}
			new_screen_index := display_screen.handleKeyEvent(event)
			if new_screen_index < len(screens) {
				display_screen = screens[new_screen_index]
				layoutAndDrawScreen(display_screen, style)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize {
			layoutAndDrawScreen(display_screen, style)
		}
	}
}

func main() {
	var err error
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", PROGRAM_NAME)
		os.Exit(1)
	}
	database_file := os.Args[1]

	db, err = bolt.Open(database_file, 0600, nil)
	if err != nil {
		fmt.Printf("Error reading file: %q\n", err.Error())
		os.Exit(1)
	}
	//defer db.Close()

	// First things first, load the database into memory
	refreshDatabase()

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	style := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)

	mainLoop(memBolt, style)
}

func refreshDatabase() {
	// Reload the database into memBolt
	memBolt = new(BoltDB)
	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(nm []byte, b *bolt.Bucket) error {
			bb := readBucket(b)
			bb.name = string(nm)
			memBolt.buckets = append(memBolt.buckets, *bb)
			return nil
		})
		return nil
	})
}

func readBucket(b *bolt.Bucket) *BoltBucket {
	bb := new(BoltBucket)
	b.ForEach(func(k, v []byte) error {
		if v == nil {
			tb := readBucket(b.Bucket(k))
			tb.name = string(k)
			bb.buckets = append(bb.buckets, *tb)
		} else {
			bb.pairs[string(k)] = string(v)
		}
		return nil
	})
	return bb
}