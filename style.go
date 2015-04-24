package main

import "github.com/nsf/termbox-go"

type Style struct {
	default_bg termbox.Attribute
	default_fg termbox.Attribute
}

func defaultStyle() Style {
	var style Style
	style.default_bg = termbox.Attribute(1)
	style.default_fg = termbox.Attribute(256)

	return style
}
