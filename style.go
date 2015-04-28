package main

import "github.com/nsf/termbox-go"

type Style struct {
	default_bg termbox.Attribute
	default_fg termbox.Attribute
	title_fg   termbox.Attribute
	title_bg   termbox.Attribute
	cursor_fg  termbox.Attribute
	cursor_bg  termbox.Attribute
}

func defaultStyle() Style {
	var style Style
	style.default_bg = termbox.ColorBlack
	style.default_fg = termbox.ColorWhite
	style.title_fg = termbox.ColorBlack
	style.title_bg = termbox.ColorRed
	style.cursor_fg = termbox.ColorBlack
	style.cursor_bg = termbox.ColorRed

	return style
}
