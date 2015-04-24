package main

import "github.com/nsf/termbox-go"

type Screen interface {
	handleKeyEvent(event termbox.Event) int
	performLayout()
	drawScreen(style Style)
}

const (
	BROWSER_SCREEN_INDEX = iota
	ABOUT_SCREEN_INDEX
	EXIT_SCREEN_INDEX
)

func defaultScreensForData(memBolt *BoltDB) []Screen {
	var view_port ViewPort
	var cursor Cursor

	browser_screen := BrowserScreen{*memBolt, cursor, view_port}
	about_screen := AboutScreen(0)
	screens := [...]Screen{
		&browser_screen,
		&about_screen,
	}

	return screens[:]
}

func drawBackground(bg termbox.Attribute) {
	termbox.Clear(0, bg)
}

func layoutAndDrawScreen(screen Screen, style Style) {
	screen.performLayout()
	drawBackground(style.default_bg)
	screen.drawScreen(style)
	termbox.Flush()
}
