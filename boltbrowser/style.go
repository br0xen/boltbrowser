package boltbrowser

import "github.com/nsf/termbox-go"

/*
termStyle Defines the colors for the terminal display, basically
*/
type termStyle struct {
	defaultBg termbox.Attribute
	defaultFg termbox.Attribute
	titleFg   termbox.Attribute
	titleBg   termbox.Attribute
	cursorFg  termbox.Attribute
	cursorBg  termbox.Attribute
}

func defaultStyle() termStyle {
	var style termStyle
	style.defaultBg = termbox.ColorBlack
	style.defaultFg = termbox.ColorWhite
	style.titleFg = termbox.ColorBlack
	style.titleBg = termbox.ColorGreen
	style.cursorFg = termbox.ColorBlack
	style.cursorBg = termbox.ColorGreen

	return style
}
