package main

import "github.com/nsf/termbox-go"

/*
Style Defines the colors for the terminal display, basically
*/
type Style struct {
	defaultBg termbox.Attribute
	defaultFg termbox.Attribute
	titleFg   termbox.Attribute
	titleBg   termbox.Attribute
	cursorFg  termbox.Attribute
	cursorBg  termbox.Attribute
}

func defaultStyle() Style {
	var style Style
	style.defaultBg = termbox.ColorBlack
	style.defaultFg = termbox.ColorWhite
	style.titleFg = termbox.ColorBlack
	style.titleBg = termbox.ColorGreen
	style.cursorFg = termbox.ColorBlack
	style.cursorBg = termbox.ColorGreen

	return style
}
