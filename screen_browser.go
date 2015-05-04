package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"gogs.bullercodeworks.com/brian/termbox-util"
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
	mode           BrowserMode
	input_modal    *termbox_util.InputModal
	confirm_modal  *termbox_util.ConfirmModal
}

type BrowserMode int

const (
	MODE_BROWSE          = 16  // 0001 0000
	MODE_CHANGE_VAL      = 32  // 0010 0000
	MODE_INSERT_BUCKET   = 48  // 0011 0000
	MODE_INSERT_PAIR     = 64  // 0100 0000
	MODE_INSERT_PAIR_KEY = 65  // 0100 0001
	MODE_INSERT_PAIR_VAL = 66  // 0100 0010
	MODE_DELETE          = 128 // 1000 0000
)

type BoltType int

const (
	TYPE_BUCKET = iota
	TYPE_PAIR
)

func (screen *BrowserScreen) handleKeyEvent(event termbox.Event) int {
	if screen.mode == 0 {
		screen.mode = MODE_BROWSE
	}
	if screen.mode == MODE_BROWSE {
		return screen.handleBrowseKeyEvent(event)
	} else if screen.mode == MODE_CHANGE_VAL {
		return screen.handleInputKeyEvent(event)
	} else if screen.mode == MODE_INSERT_BUCKET {
		return screen.handleInsertKeyEvent(event)
	} else if screen.mode == MODE_DELETE {
		return screen.handleDeleteKeyEvent(event)
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) handleBrowseKeyEvent(event termbox.Event) int {
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
		screen.jumpCursorDown(half)

	} else if event.Key == termbox.KeyCtrlB {
		_, h := termbox.Size()
		half := h / 2
		screen.jumpCursorUp(half)

	} else if event.Ch == 'j' || event.Key == termbox.KeyArrowDown {
		screen.moveCursorDown()

	} else if event.Ch == 'k' || event.Key == termbox.KeyArrowUp {
		screen.moveCursorUp()

	} else if event.Ch == 'p' {
		screen.startInsertItem(TYPE_PAIR)

	} else if event.Ch == 'b' {
		screen.startInsertItem(TYPE_BUCKET)

	} else if event.Ch == 'e' {
		b, p, _ := screen.db.getGenericFromPath(screen.current_path)
		if b != nil {
			screen.message = "Cannot edit a bucket yet"
		} else if p != nil {
			screen.startEditItem()
		}

	} else if event.Key == termbox.KeyEnter {
		b, p, _ := screen.db.getGenericFromPath(screen.current_path)
		if b != nil {
			screen.db.toggleOpenBucket(screen.current_path)
		} else if p != nil {
			screen.startEditItem()
		}

	} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight {
		b, p, _ := screen.db.getGenericFromPath(screen.current_path)
		// Select the current item
		if b != nil {
			screen.db.toggleOpenBucket(screen.current_path)
		} else if p != nil {
			screen.startEditItem()
		} else {
			screen.message = "Not sure what to do here..."
		}

	} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft {
		// If we are _on_ a bucket that's open, close it
		b, _, e := screen.db.getGenericFromPath(screen.current_path)
		if e == nil && b != nil && b.expanded {
			screen.db.closeBucket(screen.current_path)
		} else {
			if len(screen.current_path) > 1 {
				parent_bucket, err := screen.db.getBucketFromPath(screen.current_path[:len(screen.current_path)-1])
				if err == nil {
					screen.db.closeBucket(parent_bucket.path)
					// Figure out how far up we need to move the cursor
					screen.current_path = parent_bucket.path
				}
			} else {
				screen.db.closeBucket(screen.current_path)
			}
		}

	} else if event.Ch == 'D' {
		screen.startDeleteItem()
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) handleInputKeyEvent(event termbox.Event) int {
	if event.Key == termbox.KeyEsc {
		screen.mode = MODE_BROWSE
		screen.input_modal.Clear()
	} else {
		screen.input_modal.HandleKeyPress(event)
		if screen.input_modal.IsDone() {
			new_val := screen.input_modal.GetValue()
			_, p, _ := screen.db.getGenericFromPath(screen.current_path)
			if p != nil {
				if updatePairValue(screen.current_path, new_val) != nil {
					screen.message = "Error occurred updating Pair."
				} else {
					p.val = new_val
					screen.message = "Pair updated!"
				}
			}
			screen.mode = MODE_BROWSE
			screen.input_modal.Clear()
		}
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) handleDeleteKeyEvent(event termbox.Event) int {
	screen.confirm_modal.HandleKeyPress(event)
	if screen.confirm_modal.IsDone() {
		if screen.confirm_modal.IsAccepted() {
			hold_next_path := screen.db.getNextVisiblePath(screen.current_path)
			hold_prev_path := screen.db.getPrevVisiblePath(screen.current_path)
			if deleteKey(screen.current_path) == nil {
				shadow_db := screen.db
				screen.db = refreshDatabase()
				screen.db.syncOpenBuckets(shadow_db)
				// Move the current path endpoint appropriately
				//found_new_path := false
				if hold_next_path != nil {
					if len(hold_next_path) > 2 {
						if hold_next_path[len(hold_next_path)-2] == screen.current_path[len(screen.current_path)-2] {
							screen.current_path = hold_next_path
						} else if hold_prev_path != nil {
							screen.current_path = hold_prev_path
						} else {
							// Otherwise, go to the parent
							screen.current_path = screen.current_path[:(len(hold_next_path) - 2)]
						}
					} else {
						// Root bucket deleted, set to next
						screen.current_path = hold_next_path
					}
				} else if hold_prev_path != nil {
					screen.current_path = hold_prev_path
				} else {
					screen.current_path = screen.current_path[:0]
				}
			}
		}
		screen.mode = MODE_BROWSE
		screen.confirm_modal.Clear()
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) handleInsertKeyEvent(event termbox.Event) int {
	if event.Key == termbox.KeyEsc {
		screen.mode = MODE_BROWSE
		screen.input_modal.Clear()
	} else {
		/*
			} else if event.Key == termbox.KeyEnter {
				b, p, e := screen.db.getGenericFromPath(screen.current_path)
				if e == nil {
					if b != nil {
						if err := insertBucket(screen.current_path, screen.input_modal.curr_val); err != nil {
							screen.message = fmt.Sprint(err)
						} else {
							if b.parent != nil {
								//b.parent.buckets = append(b.parent.buckets, BoltBucket{name: screen.input_modal.curr_val
							} else {
								//screen.db.buckets = append(screen.db.buckets, BoltBucket{
							}
						}
						screen.mode = MODE_BROWSE
					} else if p != nil {
						if err := updatePairValue(screen.current_path, screen.input_modal.curr_val); err != nil {
							screen.message = fmt.Sprint(err)
						} else {
							p.val = screen.input_modal.curr_val
						}
						screen.mode = MODE_BROWSE
					}
				}
			} else if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 || event.Key == termbox.KeyDelete {
				screen.input_modal.curr_val = screen.input_modal.curr_val[:len(screen.input_modal.curr_val)-1]
			} else {
				screen.input_modal.curr_val += string(event.Ch)
		*/
	}
	return BROWSER_SCREEN_INDEX
}

func (screen *BrowserScreen) jumpCursorUp(distance int) bool {
	// Jump up 'distance' lines
	vis_paths, err := screen.db.buildVisiblePathSlice(nil)
	if err == nil {
		find_path := strings.Join(screen.current_path, "/")
		start_jump := false
		for i := range vis_paths {
			if vis_paths[len(vis_paths)-1-i] == find_path {
				start_jump = true
			}
			if start_jump {
				distance -= 1
				if distance == 0 {
					screen.current_path = strings.Split(vis_paths[len(vis_paths)-1-i], "/")
					break
				}
			}
		}
		if strings.Join(screen.current_path, "/") == find_path {
			screen.current_path = screen.db.getNextVisiblePath(nil)
		}
	}
	return true
}
func (screen *BrowserScreen) jumpCursorDown(distance int) bool {
	vis_paths, err := screen.db.buildVisiblePathSlice(nil)
	if err == nil {
		find_path := strings.Join(screen.current_path, "/")
		start_jump := false
		for i := range vis_paths {
			if vis_paths[i] == find_path {
				start_jump = true
			}
			if start_jump {
				distance -= 1
				if distance == 0 {
					screen.current_path = strings.Split(vis_paths[i], "/")
					break
				}
			}
		}
		if strings.Join(screen.current_path, "/") == find_path {
			screen.current_path = screen.db.getPrevVisiblePath(nil)
		}
	}
	return true
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
	if screen.mode == MODE_CHANGE_VAL || screen.mode == MODE_INSERT_BUCKET || screen.mode&MODE_INSERT_PAIR == MODE_INSERT_PAIR {
		screen.input_modal.Draw()
	}
	if screen.mode == MODE_DELETE {
		screen.confirm_modal.Draw()
	}
}

func (screen *BrowserScreen) drawHeader(style Style) {
	width, _ := termbox.Size()
	spaces := strings.Repeat(" ", (width / 2))
	termbox_util.DrawStringAtPoint(fmt.Sprintf("%s%s%s", spaces, PROGRAM_NAME, spaces), 0, 0, style.title_fg, style.title_bg)
}
func (screen *BrowserScreen) drawFooter(style Style) {
	_, height := termbox.Size()
	termbox_util.DrawStringAtPoint(fmt.Sprintf("%s(%d) - %s", screen.current_path, screen.current_type, screen.message), 0, height-1, style.default_fg, style.default_bg)
}

func (screen *BrowserScreen) drawLeftPane(style Style) {
	w, h := termbox.Size()
	if w >= 80 {
		w = w / 2
	}
	screen.view_port.number_of_rows = h - 2

	termbox_util.FillWithChar('=', 0, 1, w, 1, style.default_fg, style.default_bg)
	y := 2
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
	max_cursor := screen.view_port.number_of_rows * 2 / 3
	if cur_path_spot > max_cursor {
		tree_offset = cur_path_spot - max_cursor
	}

	for i := range screen.db.buckets {
		// The drawBucket function returns how many lines it took up
		bkt_h := screen.drawBucket(&screen.db.buckets[i], style, (y - tree_offset))
		y += bkt_h
	}
}

func (screen *BrowserScreen) drawRightPane(style Style) {
	w, h := termbox.Size()
	if w >= 80 {
		// Screen is wide enough, split it
		termbox_util.FillWithChar('=', 0, 1, w, 1, style.default_fg, style.default_bg)
		termbox_util.FillWithChar('|', (w / 2), screen.view_port.first_row-1, (w / 2), h, style.default_fg, style.default_bg)

		b, p, err := screen.db.getGenericFromPath(screen.current_path)
		if err == nil {
			start_x := (w / 2) + 2
			start_y := 2
			if b != nil {
				termbox_util.DrawStringAtPoint(fmt.Sprintf("Path: %s", strings.Join(b.path, "/")), start_x, start_y, style.default_fg, style.default_bg)
				termbox_util.DrawStringAtPoint(fmt.Sprintf("Buckets: %d", len(b.buckets)), start_x, start_y+1, style.default_fg, style.default_bg)
				termbox_util.DrawStringAtPoint(fmt.Sprintf("Pairs: %d", len(b.pairs)), start_x, start_y+2, style.default_fg, style.default_bg)
			} else if p != nil {
				termbox_util.DrawStringAtPoint(fmt.Sprintf("Path: %s", strings.Join(p.path, "/")), start_x, start_y, style.default_fg, style.default_bg)
				termbox_util.DrawStringAtPoint(fmt.Sprintf("Key: %s", p.key), start_x, start_y+1, style.default_fg, style.default_bg)
				termbox_util.DrawStringAtPoint(fmt.Sprintf("Value: %s", p.val), start_x, start_y+2, style.default_fg, style.default_bg)
			}
		}
	}
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

		termbox_util.DrawStringAtPoint(bkt_string, 0, (y + used_lines), bucket_fg, bucket_bg)
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
		termbox_util.DrawStringAtPoint(bkt_string, 0, (y + used_lines), bucket_fg, bucket_bg)
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
	termbox_util.DrawStringAtPoint(pair_string, 0, y, bucket_fg, bucket_bg)
	return 1
}

func (screen *BrowserScreen) startDeleteItem() bool {
	b, p, e := screen.db.getGenericFromPath(screen.current_path)
	if e == nil {
		w, h := termbox.Size()
		inp_w, inp_h := (w / 2), 6
		inp_x, inp_y := ((w / 2) - (inp_w / 2)), ((h / 2) - inp_h)
		mod := termbox_util.CreateConfirmModal("", inp_x, inp_y, inp_w, inp_h, termbox.ColorWhite, termbox.ColorBlack)
		if b != nil {
			mod.SetTitle(termbox_util.AlignText(fmt.Sprintf("Delete Bucket '%s'?", b.name), inp_w, termbox_util.ALIGN_CENTER))
		} else if p != nil {
			mod.SetTitle(termbox_util.AlignText(fmt.Sprintf("Delete Pair '%s'?", p.key), inp_w, termbox_util.ALIGN_CENTER))
		}
		mod.SetText("This cannot be undone!")
		screen.confirm_modal = mod
		screen.mode = MODE_DELETE
		return true
	}
	return false
}

func (screen *BrowserScreen) startEditItem() bool {
	b, p, e := screen.db.getGenericFromPath(screen.current_path)
	if e == nil {
		w, h := termbox.Size()
		inp_w, inp_h := (w / 2), 6
		inp_x, inp_y := ((w / 2) - (inp_w / 2)), ((h / 2) - inp_h)
		mod := termbox_util.CreateInputModal("", inp_x, inp_y, inp_w, inp_h, termbox.ColorWhite, termbox.ColorBlack)
		if b != nil {
			mod.SetTitle(termbox_util.AlignText(fmt.Sprintf("Rename Bucket '%s' to:", b.name), inp_w, termbox_util.ALIGN_CENTER))
			mod.SetValue(b.name)
		} else if p != nil {
			mod.SetTitle(termbox_util.AlignText(fmt.Sprintf("Input new value for '%s'", p.key), inp_w, termbox_util.ALIGN_CENTER))
			mod.SetValue(p.val)
		}
		screen.input_modal = mod
		screen.mode = MODE_CHANGE_VAL
		return true
	}
	return false
}

func (screen *BrowserScreen) startInsertItem(tp BoltType) bool {
	if tp == TYPE_BUCKET {
		screen.mode = MODE_INSERT_BUCKET

	} else if tp == TYPE_PAIR {
		screen.mode = MODE_INSERT_PAIR

	}
	return false
}

func (screen *BrowserScreen) insertBucket() bool {
	//b, _, e := screen.db.getGenericFromPath(screen.current_path)
	//if e == nil {
	w, h := termbox.Size()
	inp_w, inp_h := (w / 2), 6
	inp_x, inp_y := ((w / 2) - (inp_w / 2)), ((h / 2) - inp_h)
	mod := termbox_util.CreateInputModal("", inp_x, inp_y, inp_w, inp_h, termbox.ColorWhite, termbox.ColorBlack)
	mod.SetTitle(termbox_util.AlignText("New Bucket Name:", inp_w, termbox_util.ALIGN_CENTER))
	mod.SetValue("")
	screen.input_modal = mod
	screen.mode = MODE_INSERT_BUCKET
	return true
	//}
	return false
}

func (screen *BrowserScreen) insertPair() bool {
	return false
}

func comparePaths(p1, p2 []string) bool {
	return strings.Join(p1, "/") == strings.Join(p2, "/")
}
