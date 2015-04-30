package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
	"os"
	"syscall"
)

const PROGRAM_NAME = "boltbrowser"

var database_file string
var db *bolt.DB
var memBolt *BoltDB

func mainLoop(memBolt *BoltDB, style Style) {
	screens := defaultScreensForData(memBolt)
	display_screen := screens[BROWSER_SCREEN_INDEX]
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
	defer db.Close()
}