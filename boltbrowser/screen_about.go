package main

import (
	"fmt"

	"github.com/br0xen/termbox-util"
	"github.com/nsf/termbox-go"
)

/*
Command is a struct for associating descriptions to keys
*/
type Command struct {
	key         string
	description string
}

/*
AboutScreen is just a basic 'int' type that we can extend to make this screen
*/
type AboutScreen int

func drawCommandsAtPoint(commands []Command, x int, y int, style Style) {
	xPos, yPos := x, y
	for index, cmd := range commands {
		termboxUtil.DrawStringAtPoint(fmt.Sprintf("%6s", cmd.key), xPos, yPos, style.defaultFg, style.defaultBg)
		termboxUtil.DrawStringAtPoint(cmd.description, xPos+8, yPos, style.defaultFg, style.defaultBg)
		yPos++
		if index > 2 && index%2 == 1 {
			yPos++
		}
	}
}

func (screen *AboutScreen) handleKeyEvent(event termbox.Event) int {
	return BrowserScreenIndex
}

func (screen *AboutScreen) performLayout() {}

func (screen *AboutScreen) drawScreen(style Style) {
	defaultFg := style.defaultFg
	defaultBg := style.defaultBg
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
	if width < 100 {
		template = [...]string{
			" ____  ____  _   _____  ____  ____   ____  _     _  ____  ___  ____   ",
			"|  _ ||    || | |     ||  _ ||  _ | |    || | _ | ||    ||   ||  _ |  ",
			"| |_||| _  || | |_   _|| |_||| | || | _  || || || || ___||  _|| | ||  ",
			"|    ||| | || |   | |  |    || |_|| || | ||       |||___ | |_ | |_||_ ",
			"|  _ |||_| || |___| |  |  _ ||  _  |||_| ||       ||__  ||  _||  __  |",
			"| |_|||    ||     | |  | |_||| | | ||    ||   _   | __| || |_ | |  | |",
			"|____||____||_____|_|  |____||_| |_||____||__| |__||____||___||_|  |_|",
		}
	}
	firstLine := template[0]
	startX := (width - len(firstLine)) / 2
	//startX := (width - len(firstLine)) / 2
	startY := ((height - 2*len(template)) / 2) - 2
	xPos := startX
	yPos := startY
	if height <= 20 {
		title := "BoltBrowser"
		startY = 0
		yPos = 0
		termboxUtil.DrawStringAtPoint(title, (width-len(title))/2, startY, style.titleFg, style.titleBg)
	} else {
		if height < 25 {
			startY = 0
			yPos = 0
		}
		for _, line := range template {
			xPos = startX
			for _, runeValue := range line {
				bg := defaultBg
				displayRune := ' '
				if runeValue != ' ' {
					//bg = termbox.Attribute(125)
					displayRune = runeValue
					termbox.SetCell(xPos, yPos, displayRune, defaultFg, bg)
				}
				xPos++
			}
			yPos++
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
		{"D", "delete item"},
		{"x,X", "export as string/json to file"},

		{"?", "this screen"},
		{"q", "quit program"},
	}
	xPos = startX // + 20
	yPos++

	drawCommandsAtPoint(commands1[:], xPos, yPos+1, style)
	drawCommandsAtPoint(commands2[:], xPos+20, yPos+1, style)
	exitTxt := "Press any key to return to browser"
	termboxUtil.DrawStringAtPoint(exitTxt, (width-len(exitTxt))/2, height-1, style.titleFg, style.titleBg)
}
