package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

type Command struct {
	key         string
	description string
}

type AboutScreen int

func drawCommandsAtPoint(commands []Command, x int, y int, style Style) {
	x_pos, y_pos := x, y
	for index, cmd := range commands {
		drawStringAtPoint(fmt.Sprintf("%6s", cmd.key), x_pos, y_pos, style.default_fg, style.default_bg)
		drawStringAtPoint(cmd.description, x_pos+8, y_pos, style.default_fg, style.default_bg)
		y_pos++
		if index > 2 && index%2 == 1 {
			y_pos++
		}
	}
}

func (screen *AboutScreen) handleKeyEvent(event termbox.Event) int {
	return BROWSER_SCREEN_INDEX
}

func (screen *AboutScreen) performLayout() {}

func (screen *AboutScreen) drawScreen(style Style) {
	default_fg := style.default_fg
	default_bg := style.default_bg
	width, height := termbox.Size()
	template := [...]string{
		" _______  _______  ___    _______  _______  ______    _______  _     _  _______  _______  ______   ",
		"|  _    ||       ||   |  |       ||  _    ||    _ |  |       || | _ | ||       ||       ||    _ |  ",
		"| |_|   ||   _   ||   |  |_     _|| |_|   ||   | ||  |   _   || || || ||  _____||    ___||   | ||  ",
		"|       ||  | |  ||   |    |   |  |       ||   |_||_ |  | |  ||       || |_____ |   |___ |   |_||_ ",
		"|  _   | |  |_|  ||   |___ |   |  |  _   | |    __  ||  |_|  ||       ||_____  ||    ___||    __  |",
		"| |_|   ||       ||       ||   |  | |_|   ||   |  | ||       ||   _   | _____| ||   |___ |   |  | |",
		"|_______||_______||_______||___|  |_______||___|  |_||_______||__| |__||_______||_______||___|  |_|",
	}
	first_line := template[0]
	start_x := (width - len(first_line)) / 2
	start_y := (height - 2*len(template)) / 2
	x_pos := start_x
	y_pos := start_y
	for _, line := range template {
		x_pos = start_x
		for _, runeValue := range line {
			bg := default_bg
			displayRune := ' '
			if runeValue != ' ' {
				//bg = termbox.Attribute(125)
				displayRune = runeValue
				termbox.SetCell(x_pos, y_pos, displayRune, default_fg, bg)
			}
			x_pos++
		}
		y_pos++
	}

	commands1 := [...]Command{
		{"h", "close parent"},
		{"j", "down"},
		{"k", "up"},
		{"l", "open item"},

		{"g", "goto top"},
		{"G", "goto bottom"},

		{"ctrl-e", "scroll down"},
		{"ctrl-y", "scroll up"},

		{"ctrl-f", "page down"},
		{"ctrl-b", "page up"},
	}

	commands2 := [...]Command{
		{"p", "create pair"},
		{"b", "create bucket"},
		{"d", "delete item"},

		{"?", "this screen"},
		{"q", "quit program"},
	}
	x_pos = start_x + 3
	y_pos++

	drawCommandsAtPoint(commands1[:], x_pos, y_pos+1, style)
	drawCommandsAtPoint(commands2[:], x_pos+20, y_pos+1, style)
}
