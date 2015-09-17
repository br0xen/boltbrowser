package main

import "github.com/nsf/termbox-go"

// Screen is a basic structure for all of the applications screens
type Screen interface {
	handleKeyEvent(event termbox.Event) int
	performLayout()
	drawScreen(style Style)
}

const (
	// BrowserScreenIndex is the index
	BrowserScreenIndex = iota
	// AboutScreenIndex The idx number for the 'About' Screen
	AboutScreenIndex
	// ExitScreenIndex The idx number for Exiting
	ExitScreenIndex
)

func defaultScreensForData(db *BoltDB) []Screen {
	var viewPort ViewPort

	browserScreen := BrowserScreen{db: db, viewPort: viewPort}
	aboutScreen := AboutScreen(0)
	screens := [...]Screen{
		&browserScreen,
		&aboutScreen,
	}

	return screens[:]
}

func drawBackground(bg termbox.Attribute) {
	termbox.Clear(0, bg)
}

func layoutAndDrawScreen(screen Screen, style Style) {
	screen.performLayout()
	drawBackground(style.defaultBg)
	screen.drawScreen(style)
	termbox.Flush()
}
