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
	db             *BoltDB
	cursor         Cursor
	view_port      ViewPort
	queued_command string
	current_path   []string
	current_type   int
	message        string
}

type BoltType int

const (
	TYPE_BUCKET = iota
	TYPE_PAIR
)

func (screen *BrowserScreen) handleKeyEvent(event termbox.Event) int {
	if event.Ch == '?' {
		// About
		return ABOUT_SCREEN_INDEX
	} else if event.Ch == 'q' || event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
		// Quit
		return EXIT_SCREEN_INDEX
	} else if event.Ch == 'g' {
		// Jump to Beginning
		screen.current_path = screen.db.getNextVisiblePath(nil)
	} else if event.Ch == 'G' {
		// Jump to End
		screen.current_path = screen.db.getPrevVisiblePath(nil)
	} else if event.Key == termbox.KeyCtrlF {
		// Jump forward half a screen
		_, h := termbox.Size()
		half := h / 2
		vis_paths, err := screen.db.buildVisiblePathSlice(nil)
		if err == nil {
			find_path := strings.Join(screen.current_path, "/")
			start_jump := false
			for i := range vis_paths {
				if vis_paths[i] == find_path {
					start_jump = true
				}
				if start_jump {
					half -= 1
					if half == 0 {
						screen.current_path = strings.Split(vis_paths[i], "/")
						break
					}
				}
			}
			if strings.Join(screen.current_path, "/") == find_path {
				screen.current_path = screen.db.getPrevVisiblePath(nil)
			}
		}
	} else if event.Key == termbox.KeyCtrlB {
		// Jump back half a screen
		_, h := termbox.Size()
		half := h / 2
		vis_paths, err := screen.db.buildVisiblePathSlice(nil)
		if err == nil {
			find_path := strings.Join(screen.current_path, "/")
			start_jump := false
			for i := range vis_paths {
				if vis_paths[len(vis_paths)-1-i] == find_path {
					start_jump = true
				}
				if start_jump {
					half -= 1
					if half == 0 {
						screen.current_path = strings.Split(vis_paths[len(vis_paths)-1-i], "/")
						break
					}
				}
			}
			if strings.Join(screen.current_path, "/") == find_path {
				screen.current_path = screen.db.getNextVisiblePath(nil)
			}
		}
	} else if event.Ch == 'j' || event.Key == termbox.KeyArrowDown {
		screen.moveCursorDown()
	} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp {
		screen.moveCursorUp()
	} else if event.Ch == 'e' {
		screen.queued_command = "edit"
	} else if event.Key == termbox.KeyEnter {
		b, p, _ := screen.db.getGenericFromPath(screen.current_path)
		if b != nil {
			toggleOpenBucket(screen.current_path)
		} else if p != nil {
		} else {
			screen.message = "Not sure what to do here..."
		}
	} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight {
		b, p, _ := screen.db.getGenericFromPath(screen.current_path)
		// Select the current item
		if b != nil {
			toggleOpenBucket(screen.current_path)
		} else if p != nil {
		} else {
			screen.message = "Not sure what to do here..."
		}
	} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft {
		// If we are _on_ a bucket that's open, close it
		b, _, e := screen.db.getGenericFromPath(screen.current_path)
		if e == nil && b != nil && b.expanded {
			closeBucket(screen.current_path)
		} else {
			if len(screen.current_path) > 1 {
				parent_bucket, err := screen.db.getBucketFromPath(screen.current_path[:len(screen.current_path)-1])
				if err == nil {
					closeBucket(parent_bucket.path)
					// Figure out how far up we need to move the cursor
					screen.current_path = parent_bucket.path
				}
			} else {
				closeBucket(screen.current_path)
			}
		}
	} else if event.Ch == 'D' {
		deleteKey(screen.current_path)
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) moveCursorUp() bool {
	new_path := screen.db.getPrevVisiblePath(screen.current_path)
	if new_path != nil {
		screen.current_path = new_path
		return true
	}
	return false
}
func (screen *BrowserScreen) moveCursorDown() bool {
	new_path := screen.db.getNextVisiblePath(screen.current_path)
	if new_path != nil {
		screen.current_path = new_path
		return true
	}
	return false
}

func (screen *BrowserScreen) performLayout() {}

func (screen *BrowserScreen) drawScreen(style Style) {
	screen.drawLeftPane(style)
	screen.drawRightPane(style)
	screen.drawHeader(style)
	screen.drawFooter(style)
}

func (screen *BrowserScreen) drawHeader(style Style) {
	width, _ := termbox.Size()
	spaces := strings.Repeat(" ", (width / 2))
	drawStringAtPoint(fmt.Sprintf("%s%s%s", spaces, PROGRAM_NAME, spaces), 0, 0, style.title_fg, style.title_bg)
}
func (screen *BrowserScreen) drawFooter(style Style) {
	_, height := termbox.Size()
	drawStringAtPoint(fmt.Sprintf("%s(%d) - %s", screen.current_path, screen.current_type, screen.message), 0, height-1, style.default_fg, style.default_bg)
}

func (screen *BrowserScreen) drawLeftPane(style Style) {
	w, h := termbox.Size()
	if w >= 80 {
		w = w / 2
	}
	screen.view_port.number_of_rows = h - 2
	_, y := 1, 2
	screen.view_port.first_row = y
	if len(screen.current_path) == 0 {
		screen.current_path = screen.db.getNextVisiblePath(nil)
	}

	// So we know how much of the tree _wants_ to be visible
	// we only have screen.view_port.number_of_rows of space though
	cur_path_spot := 0
	vis_slice, err := screen.db.buildVisiblePathSlice(nil)
	if err == nil {
		for i := range vis_slice {
			if strings.Join(screen.current_path, "/") == vis_slice[i] {
				cur_path_spot = i
			}
		}
	}

	tree_offset := 0
	half_screen := screen.view_port.number_of_rows / 2
	if cur_path_spot > half_screen {
		tree_offset = cur_path_spot - half_screen
	}

	screen.message = fmt.Sprintf("Offset: %d", tree_offset)
	for i := range screen.db.buckets {
		// The drawBucket function returns how many lines it took up
		bkt_h := screen.drawBucket(&screen.db.buckets[i], style, (y - tree_offset))
		y += bkt_h
	}
}

func (screen *BrowserScreen) drawRightPane(style Style) {
	w, _ := termbox.Size()
	vis_slice, err := screen.db.buildVisiblePathSlice(nil)
	if err == nil {
		for i := range vis_slice {
			if strings.Join(screen.current_path, "/") == vis_slice[i] {
				drawStringAtPoint(vis_slice[i], (w/2)+2, i+1, style.title_fg, style.title_bg)
			} else {
				drawStringAtPoint(vis_slice[i], (w/2)+2, i+1, style.default_fg, style.default_bg)
			}
		}
	}
	/*
		w, h := termbox.Size()
		if w >= 80 {
			// Screen is wide enough, split it
			fillWithChar('|', (w / 2), screen.view_port.first_row-1, (w / 2), h, style.default_fg, style.default_bg)

			b, p, err := screen.db.getGenericFromPath(screen.current_path)
			if err == nil {
				start_x := (w / 2) + 1
				parent_str := "/"
				if b != nil {
					if b.parent != nil {
						parent_str = b.parent.name
					}
					drawStringAtPoint(fmt.Sprintf("Parent: %s", parent_str), start_x, 1, style.default_fg, style.default_bg)
					drawStringAtPoint(fmt.Sprintf("Buckets: %d", len(b.buckets)), start_x, 2, style.default_fg, style.default_bg)
					drawStringAtPoint(fmt.Sprintf("Pairs: %d", len(b.pairs)), start_x, 3, style.default_fg, style.default_bg)
					drawStringAtPoint(fmt.Sprintf("Path: %s", strings.Join(b.path, "/")), start_x, 4, style.default_fg, style.default_bg)
				} else if p != nil {
					if p.parent != nil {
						parent_str = p.parent.name
					}
					drawStringAtPoint(fmt.Sprintf("Parent: %s", parent_str), start_x, 1, style.default_fg, style.default_bg)
					drawStringAtPoint(fmt.Sprintf("Key: %s", p.key), start_x, 2, style.default_fg, style.default_bg)
					drawStringAtPoint(fmt.Sprintf("Value: %s", p.val), start_x, 3, style.default_fg, style.default_bg)
					drawStringAtPoint(fmt.Sprintf("Path: %s", strings.Join(p.path, "/")), start_x, 4, style.default_fg, style.default_bg)
				}
			}
		}
	*/
}

/* drawBucket
 * @bkt *BoltBucket - The bucket to draw
 * @style Style - The style to use
 * @w int - The Width of the lines
 * @y int - The Y position to start drawing
 * return - The number of lines used
 */
func (screen *BrowserScreen) drawBucket(bkt *BoltBucket, style Style, y int) int {
	w, _ := termbox.Size()
	if w >= 80 {
		w = w / 2
	}
	used_lines := 0
	bucket_fg := style.default_fg
	bucket_bg := style.default_bg
	if comparePaths(screen.current_path, bkt.path) {
		bucket_fg = style.cursor_fg
		bucket_bg = style.cursor_bg
	}

	bkt_string := strings.Repeat(" ", screen.db.getDepthFromPath(bkt.path)*2)
	if bkt.expanded {
		bkt_string = bkt_string + "- " + bkt.name + " "
		bkt_string = fmt.Sprintf("%s%s", bkt_string, strings.Repeat(" ", (w-len(bkt_string))))

		drawStringAtPoint(bkt_string, 0, (y + used_lines), bucket_fg, bucket_bg)
		used_lines += 1

		for i := range bkt.buckets {
			used_lines += screen.drawBucket(&bkt.buckets[i], style, y+used_lines)
		}
		for i := range bkt.pairs {
			used_lines += screen.drawPair(&bkt.pairs[i], style, y+used_lines)
		}
	} else {
		bkt_string = bkt_string + "+ " + bkt.name
		bkt_string = fmt.Sprintf("%s%s", bkt_string, strings.Repeat(" ", (w-len(bkt_string))))
		drawStringAtPoint(bkt_string, 0, (y + used_lines), bucket_fg, bucket_bg)
		used_lines += 1
	}
	return used_lines
}

func (screen *BrowserScreen) drawPair(bp *BoltPair, style Style, y int) int {
	w, _ := termbox.Size()
	if w >= 80 {
		w = w / 2
	}
	bucket_fg := style.default_fg
	bucket_bg := style.default_bg
	if comparePaths(screen.current_path, bp.path) {
		bucket_fg = style.cursor_fg
		bucket_bg = style.cursor_bg
	}

	pair_string := strings.Repeat(" ", screen.db.getDepthFromPath(bp.path)*2)
	pair_string = fmt.Sprintf("%s%s: %s", pair_string, bp.key, bp.val)
	pair_string = fmt.Sprintf("%s%s", pair_string, strings.Repeat(" ", (w-len(pair_string))))
	drawStringAtPoint(pair_string, 0, y, bucket_fg, bucket_bg)
	return 1
}

func comparePaths(p1, p2 []string) bool {
	return strings.Join(p1, "/") == strings.Join(p2, "/")
}
