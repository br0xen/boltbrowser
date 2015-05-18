package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"gogs.bullercodeworks.com/brian/termbox-util"
)

type Command struct {
	key         string
	description string
}

type AboutScreen int

func drawCommandsAtPoint(commands []Command, x int, y int, style Style) {
	x_pos, y_pos := x, y
	for index, cmd := range commands {
		termbox_util.DrawStringAtPoint(fmt.Sprintf("%6s", cmd.key), x_pos, y_pos, style.default_fg, style.default_bg)
		termbox_util.DrawStringAtPoint(cmd.description, x_pos+8, y_pos, style.default_fg, style.default_bg)
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
	start_y := ((height - 2*len(template)) / 2) - 2
	x_pos := start_x
	y_pos := start_y
	if height <= 20 {
		title := "BoltBrowser"
		start_y = 0
		y_pos = 0
		termbox_util.DrawStringAtPoint(title, (width-len(title))/2, start_y, style.title_fg, style.title_bg)
	} else {
		if height < 25 {
			start_y = 0
			y_pos = 0
		}
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
	}

	commands1 := [...]Command{
		{"h,←", "close parent"},
		{"j,↓", "down"},
		{"k,↑", "up"},
		{"l,→", "open item"},

		{"g", "goto top"},
		{"G", "goto bottom"},
		{"ctrl+f", "jump down"},
		{"ctrl+b", "jump up"},
	}

	commands2 := [...]Command{
		{"p,P", "create pair/at parent"},
		{"b,B", "create bucket/at parent"},
		{"e", "edit value of pair"},
		{"r", "rename pair/bucket"},
		{"d", "delete item"},

		{"?", "this screen"},
		{"q", "quit program"},
	}
	x_pos = start_x + 20
	y_pos++

	drawCommandsAtPoint(commands1[:], x_pos, y_pos+1, style)
	drawCommandsAtPoint(commands2[:], x_pos+20, y_pos+1, style)
	exit_txt := "Press any key to return to browser"
	termbox_util.DrawStringAtPoint(exit_txt, (width-len(exit_txt))/2, height-1, style.title_fg, style.title_bg)
}
