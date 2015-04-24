package main

import (
	"github.com/nsf/termbox-go"
)

type ViewPort struct {
	bytes_per_row  int
	number_of_rows int
	first_row      int
}

type BrowserScreen struct {
	memBolt   BoltDB
	cursor    Cursor
	view_port ViewPort
}

func (screen *BrowserScreen) handleKeyEvent(event termbox.Event) int {
	if event.Ch == '?' { // About
		return ABOUT_SCREEN_INDEX
	} else if event.Ch == 'q' || event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
		return EXIT_SCREEN_INDEX
	}

	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) performLayout() {
	/*
		width, height := termbox.Size()
		legend_height := heightOfWidgets()
		line_height := 3
		cursor := screen.cursor
		cursor_row_within_view_port := 0

		var new_view_port ViewPort
		new_view_port.bytes_per_row = (width - 3) / 3
		new_view_port.number_of_rows = (height - 1 - legend_height) / line_height
		new_view_port.first_row = screen.view_port.first_row

		if new_view_port.first_row < 0 {
			new_view_port.first_row = 0
		}

		screen.view_port = new_view_port
	*/
}

func (screen *BrowserScreen) drawScreen(style Style) {
	x, y := 2, 1
	//	x_pad := 2
	//line_height := 1
	//	width, height := termbox.Size()
	//	widget_width, widget_height := drawWidgets(screen.cursor, style)

	//	cursor := screen.cursor
	//view_port := screen.view_port

	//last_y := y + view_port.number_of_rows*line_height - 1
	//last_x := x + view_port.bytes_per_row*3 - 1

	start_x := x
	for _, bkt := range screen.memBolt.buckets {
		bucket_fg := style.default_fg
		bucket_bg := style.default_bg
		termbox.SetCell(x, y, '+', bucket_fg, bucket_bg)
		x = drawStringAtPoint(bkt.name, x+1, y, bucket_fg, bucket_bg)
		y = y + 1
		x = start_x
	}
}
