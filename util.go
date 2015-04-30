package main

import "github.com/nsf/termbox-go"

func drawStringAtPoint(str string, x int, y int, fg termbox.Attribute, bg termbox.Attribute) int {
	x_pos := x
	for _, runeValue := range str {
		termbox.SetCell(x_pos, y, runeValue, fg, bg)
		x_pos++
	}
	return x_pos - x
}

func fillWithChar(r rune, x1, y1, x2, y2 int, fg termbox.Attribute, bg termbox.Attribute) {
	for xx := x1; xx <= x2; xx++ {
		for yx := y1; yx <= y2; yx++ {
			termbox.SetCell(xx, yx, r, fg, bg)
		}
	}
}
