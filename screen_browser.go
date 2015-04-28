package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
)

type ViewPort struct {
	bytes_per_row  int
	number_of_rows int
	first_row      int
}

type BrowserScreen struct {
	memBolt        BoltDB
	cursor         Cursor
	view_port      ViewPort
	queued_command string
}

func (screen *BrowserScreen) handleKeyEvent(event termbox.Event) int {
	if event.Ch == '?' { // About
		return ABOUT_SCREEN_INDEX
	} else if event.Ch == 'q' || event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
		return EXIT_SCREEN_INDEX
	} else if event.Ch == 'j' {
		// Move the cursor down
		if screen.cursor.y < screen.view_port.number_of_rows-1 {
			screen.cursor.y++
		}
	} else if event.Ch == 'k' {
		// Move the cursor up
		if screen.cursor.y > screen.view_port.first_row {
			screen.cursor.y--
		}
	} else if event.Ch == 'e' {
		screen.queued_command = "edit"
	} else if event.Key == termbox.KeyEnter {
		// Select the current item
		screen.queued_command = "select"
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) performLayout() {}

func (screen *BrowserScreen) drawScreen(style Style) {
	width, _ := termbox.Size()
	spaces := strings.Repeat(" ", (width/2)-(len(PROGRAM_NAME)-2))
	drawStringAtPoint(fmt.Sprintf("%s%s%s", spaces, PROGRAM_NAME, spaces), 0, 0, style.title_fg, style.title_bg)

	x, y := 2, 2
	screen.view_port.first_row = y
	if screen.cursor.y == 0 {
		screen.cursor.y = y
	}
	for _, bkt := range screen.memBolt.buckets {
		_, y = screen.drawBucket(&bkt, style, x, y)
	}
	screen.view_port.number_of_rows = y
}

func (screen *BrowserScreen) drawBucket(b *BoltBucket, style Style, x, y int) (int, int) {
	bkt := b
	bucket_fg := style.default_fg
	bucket_bg := style.default_bg
	if y == screen.cursor.y {
		bucket_fg = style.cursor_fg
		bucket_bg = style.cursor_bg
		if screen.queued_command == "select" {
			// Expand/Collapse the bucket
			bkt.expanded = !bkt.expanded
			screen.queued_command = ""
		}
	}
	bkt_string := " "
	start_x := x
	if bkt.expanded {
		bkt_string = bkt_string + "- " + bkt.name + " "
		x = drawStringAtPoint(bkt_string, x, y, bucket_fg, bucket_bg)
		y = y + 1

		for _, ib := range bkt.buckets {
			_, y = screen.drawBucket(&ib, style, start_x+2, y)
		}
		for _, ip := range bkt.pairs {
			_, y = screen.drawPair(ip, style, x, y)
		}
	} else {
		bkt_string = bkt_string + "+ " + bkt.name + " "
		x = drawStringAtPoint(bkt_string, x, y, bucket_fg, bucket_bg)
		y = y + 1
	}

	return x, y
}

func (screen *BrowserScreen) drawPair(bp BoltPair, style Style, x, y int) (int, int) {
	bucket_fg := style.default_fg
	bucket_bg := style.default_bg
	if y == screen.cursor.y {
		bucket_fg = style.cursor_fg
		bucket_bg = style.cursor_bg
	}

	pair_string := fmt.Sprintf("%s: %s", bp.key, bp.val)
	x = drawStringAtPoint(pair_string, x, y, bucket_fg, bucket_bg)
	y = y + 1
	return x, y
}
