package boltbrowser

import "github.com/nsf/termbox-go"

// screen is a basic structure for all of the applications screens
type screen interface {
	handleKeyEvent(event termbox.Event) int
	performLayout()
	drawScreen(style termStyle)
}

const (
	// browserScreenIndex is the index
	browserScreenIndex = iota
	// aboutScreenIndex The idx number for the 'About' Screen
	aboutScreenIndex
	// exitScreenIndex The idx number for Exiting
	exitScreenIndex
)

func defaultScreensForData(db *boltDB) []screen {
	var vp viewPort

	browserScr := browserScreen{db: db, viewPort: vp}
	aboutScr := aboutScreen(0)
	screens := [...]screen{
		&browserScr,
		&aboutScr,
	}

	return screens[:]
}

func drawBackground(bg termbox.Attribute) {
	termbox.Clear(0, bg)
}

func layoutAndDrawScreen(scr screen, style termStyle) {
	scr.performLayout()
	drawBackground(style.defaultBg)
	scr.drawScreen(style)
	termbox.Flush()
}
