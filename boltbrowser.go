package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
)

/*
ProgramName is the name of the program
*/
const ProgramName = "boltbrowser"

var databaseFile string
var db *bolt.DB
var memBolt *BoltDB

func mainLoop(memBolt *BoltDB, style Style) {
	screens := defaultScreensForData(memBolt)
	displayScreen := screens[BrowserScreenIndex]
	layoutAndDrawScreen(displayScreen, style)
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyCtrlZ {
				process, _ := os.FindProcess(os.Getpid())
				termbox.Close()
				process.Signal(syscall.SIGSTOP)
				termbox.Init()
			}
			newScreenIndex := displayScreen.handleKeyEvent(event)
			if newScreenIndex < len(screens) {
				displayScreen = screens[newScreenIndex]
				layoutAndDrawScreen(displayScreen, style)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize {
			layoutAndDrawScreen(displayScreen, style)
		}
	}
}

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>\n", ProgramName)
		os.Exit(1)
	}

	databaseFile := os.Args[1]
	db, err = bolt.Open(databaseFile, 0600, nil)
	if err != nil {
		fmt.Printf("Error reading file: %q\n", err.Error())
		os.Exit(1)
	}

	// First things first, load the database into memory
	memBolt.refreshDatabase()

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