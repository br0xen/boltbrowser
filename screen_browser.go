package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	termboxUtil "github.com/br0xen/termbox-util"
	"github.com/nsf/termbox-go"
)

/*
ViewPort helps keep track of what's being displayed on the screen
*/
type ViewPort struct {
	bytesPerRow  int
	numberOfRows int
	firstRow     int
	scrollRow    int
}

/*
BrowserScreen holds all that's going on :D
*/
type BrowserScreen struct {
	db             *BoltDB
	leftViewPort   ViewPort
	rightViewPort  ViewPort
	queuedCommand  string
	currentPath    []string
	currentType    int
	message        string
	filter         string
	mode           BrowserMode
	inputModal     *termboxUtil.InputModal
	confirmModal   *termboxUtil.ConfirmModal
	messageTimeout time.Duration
	messageTime    time.Time

	leftPaneBuffer  []Line
	rightPaneBuffer []Line
}

/*
BrowserMode is just for designating the mode that we're in
*/
type BrowserMode int

const (
	modeBrowse        = 16  // 0000 0001 0000
	modeChange        = 32  // 0000 0010 0000
	modeChangeKey     = 33  // 0000 0010 0001
	modeChangeVal     = 34  // 0000 0010 0010
	modeFilter        = 35  // 0100 0010 0011
	modeInsert        = 64  // 0000 0100 0000
	modeInsertBucket  = 65  // 0000 0100 0001
	modeInsertPair    = 68  // 0000 0100 0100
	modeInsertPairKey = 69  // 0000 0100 0101
	modeInsertPairVal = 70  // 0000 0100 0110
	modeDelete        = 256 // 0001 0000 0000
	modeModToParent   = 8   // 0000 0000 1000
	modeExport        = 512 // 0010 0000 0000
	modeExportValue   = 513 // 0010 0000 0001
	modeExportJSON    = 514 // 0010 0000 0010
)

/*
BoltType is just for tracking what type of db item we're looking at
*/
type BoltType int

const (
	typeBucket = iota
	typePair
)

func (screen *BrowserScreen) handleKeyEvent(event termbox.Event) int {
	if screen.mode == 0 {
		screen.mode = modeBrowse
	}
	if screen.mode == modeBrowse {
		return screen.handleBrowseKeyEvent(event)
	} else if screen.mode&modeChange == modeChange {
		return screen.handleInputKeyEvent(event)
	} else if screen.mode&modeInsert == modeInsert {
		return screen.handleInsertKeyEvent(event)
	} else if screen.mode == modeDelete {
		return screen.handleDeleteKeyEvent(event)
	} else if screen.mode&modeExport == modeExport {
		return screen.handleExportKeyEvent(event)
	}
	return BrowserScreenIndex
}

func (screen *BrowserScreen) handleBrowseKeyEvent(event termbox.Event) int {
	if event.Ch == '?' {
		// About
		return AboutScreenIndex

	} else if event.Ch == 'q' || event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC {
		// Quit
		return ExitScreenIndex

	} else if event.Ch == 'g' {
		// Jump to Beginning
		screen.currentPath = screen.db.getNextVisiblePath(nil, screen.filter)

	} else if event.Ch == 'G' {
		// Jump to End
		screen.currentPath = screen.db.getPrevVisiblePath(nil, screen.filter)

	} else if event.Key == termbox.KeyCtrlR {
		screen.refreshDatabase()

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

	} else if event.Ch == 'J' {
		screen.moveRightPaneDown()

	} else if event.Ch == 'K' {
		screen.moveRightPaneUp()

	} else if event.Ch == 'p' {
		// p creates a new pair at the current level
		screen.startInsertItem(typePair)
	} else if event.Ch == 'P' {
		// P creates a new pair at the parent level
		screen.startInsertItemAtParent(typePair)

	} else if event.Ch == 'b' {
		// b creates a new bucket at the current level
		screen.startInsertItem(typeBucket)
	} else if event.Ch == 'B' {
		// B creates a new bucket at the parent level
		screen.startInsertItemAtParent(typeBucket)

	} else if event.Ch == 'e' {
		b, p, _ := screen.db.getGenericFromPath(screen.currentPath)
		if b != nil {
			screen.setMessage("Cannot edit a bucket, did you mean to (r)ename?")
		} else if p != nil {
			screen.startEditItem()
		}
	} else if event.Ch == '/' {
		screen.startFilter()

	} else if event.Ch == 'r' {
		screen.startRenameItem()

	} else if event.Key == termbox.KeyEnter {
		b, p, _ := screen.db.getGenericFromPath(screen.currentPath)
		if b != nil {
			screen.db.toggleOpenBucket(screen.currentPath)
		} else if p != nil {
			screen.startEditItem()
		}

	} else if event.Ch == 'l' || event.Key == termbox.KeyArrowRight {
		b, p, _ := screen.db.getGenericFromPath(screen.currentPath)
		// Select the current item
		if b != nil {
			screen.db.toggleOpenBucket(screen.currentPath)
		} else if p != nil {
			screen.startEditItem()
		} else {
			screen.setMessage("Not sure what to do here...")
		}

	} else if event.Ch == 'h' || event.Key == termbox.KeyArrowLeft {
		// If we are _on_ a bucket that's open, close it
		b, _, e := screen.db.getGenericFromPath(screen.currentPath)
		if e == nil && b != nil && b.expanded {
			screen.db.closeBucket(screen.currentPath)
		} else {
			if len(screen.currentPath) > 1 {
				parentBucket, err := screen.db.getBucketFromPath(screen.currentPath[:len(screen.currentPath)-1])
				if err == nil {
					screen.db.closeBucket(parentBucket.GetPath())
					// Figure out how far up we need to move the cursor
					screen.currentPath = parentBucket.GetPath()
				}
			} else {
				screen.db.closeBucket(screen.currentPath)
			}
		}

	} else if event.Ch == 'D' {
		screen.startDeleteItem()
	} else if event.Ch == 'x' {
		// Export Value
		screen.startExportValue()
	} else if event.Ch == 'X' {
		// Export Key/Value (or Bucket) as JSON
		screen.startExportJSON()
	}
	return BrowserScreenIndex
}

func (screen *BrowserScreen) handleInputKeyEvent(event termbox.Event) int {
	if event.Key == termbox.KeyEsc {
		screen.mode = modeBrowse
		screen.inputModal.Clear()
	} else {
		screen.inputModal.HandleEvent(event)
		if screen.inputModal.IsDone() {
			b, p, _ := screen.db.getGenericFromPath(screen.currentPath)
			if screen.mode == modeFilter {
				screen.filter = screen.inputModal.GetValue()
				if !screen.db.isVisiblePath(screen.currentPath, screen.filter) {
					screen.currentPath = screen.currentPath[:len(screen.currentPath)-1]
				}
			}
			if b != nil {
				if screen.mode == modeChangeKey {
					newName := screen.inputModal.GetValue()
					if renameBucket(screen.currentPath, newName) != nil {
						screen.setMessage("Error renaming bucket.")
					} else {
						b.name = newName
						screen.currentPath[len(screen.currentPath)-1] = b.name
						screen.setMessage("Bucket Renamed!")
						screen.refreshDatabase()
					}
				}
			} else if p != nil {
				if screen.mode == modeChangeKey {
					newKey := screen.inputModal.GetValue()
					if updatePairKey(screen.currentPath, newKey) != nil {
						screen.setMessage("Error occurred updating Pair.")
					} else {
						p.key = newKey
						screen.currentPath[len(screen.currentPath)-1] = p.key
						screen.setMessage("Pair updated!")
						screen.refreshDatabase()
					}
				} else if screen.mode == modeChangeVal {
					newVal := screen.inputModal.GetValue()
					if updatePairValue(screen.currentPath, newVal) != nil {
						screen.setMessage("Error occurred updating Pair.")
					} else {
						p.val = newVal
						screen.setMessage("Pair updated!")
						screen.refreshDatabase()
					}
				}
			}
			screen.mode = modeBrowse
			screen.inputModal.Clear()
		}
	}
	return BrowserScreenIndex
}

func (screen *BrowserScreen) handleDeleteKeyEvent(event termbox.Event) int {
	screen.confirmModal.HandleEvent(event)
	if screen.confirmModal.IsDone() {
		if screen.confirmModal.IsAccepted() {
			holdNextPath := screen.db.getNextVisiblePath(screen.currentPath, screen.filter)
			holdPrevPath := screen.db.getPrevVisiblePath(screen.currentPath, screen.filter)
			if deleteKey(screen.currentPath) == nil {
				screen.refreshDatabase()
				// Move the current path endpoint appropriately
				//found_new_path := false
				if holdNextPath != nil {
					if len(holdNextPath) > 2 {
						if holdNextPath[len(holdNextPath)-2] == screen.currentPath[len(screen.currentPath)-2] {
							screen.currentPath = holdNextPath
						} else if holdPrevPath != nil {
							screen.currentPath = holdPrevPath
						} else {
							// Otherwise, go to the parent
							screen.currentPath = screen.currentPath[:(len(holdNextPath) - 2)]
						}
					} else {
						// Root bucket deleted, set to next
						screen.currentPath = holdNextPath
					}
				} else if holdPrevPath != nil {
					screen.currentPath = holdPrevPath
				} else {
					screen.currentPath = screen.currentPath[:0]
				}
			}
		}
		screen.mode = modeBrowse
		screen.confirmModal.Clear()
	}
	return BrowserScreenIndex
}

func (screen *BrowserScreen) handleInsertKeyEvent(event termbox.Event) int {
	if event.Key == termbox.KeyEsc {
		if len(screen.db.buckets) == 0 {
			return ExitScreenIndex
		}
		screen.mode = modeBrowse
		screen.inputModal.Clear()
	} else {
		screen.inputModal.HandleEvent(event)
		if screen.inputModal.IsDone() {
			newVal := screen.inputModal.GetValue()
			screen.inputModal.Clear()
			var insertPath []string
			if len(screen.currentPath) > 0 {
				_, p, e := screen.db.getGenericFromPath(screen.currentPath)
				if e != nil {
					screen.setMessage("Error Inserting new item. Invalid Path.")
				}
				insertPath = screen.currentPath
				// where are we inserting?
				if p != nil {
					// If we're sitting on a pair, we have to go to it's parent
					screen.mode = screen.mode | modeModToParent
				}
				if screen.mode&modeModToParent == modeModToParent {
					if len(screen.currentPath) > 1 {
						insertPath = screen.currentPath[:len(screen.currentPath)-1]
					} else {
						insertPath = make([]string, 0)
					}
				}
			}

			parentB, _, _ := screen.db.getGenericFromPath(insertPath)
			if screen.mode&modeInsertBucket == modeInsertBucket {
				err := insertBucket(insertPath, newVal)
				if err != nil {
					screen.setMessage(fmt.Sprintf("%s => %s", err, insertPath))
				} else {
					if parentB != nil {
						parentB.expanded = true
					}
				}
				screen.currentPath = append(insertPath, newVal)

				screen.refreshDatabase()
				screen.mode = modeBrowse
				screen.inputModal.Clear()
			} else if screen.mode&modeInsertPair == modeInsertPair {
				err := insertPair(insertPath, newVal, "")
				if err != nil {
					screen.setMessage(fmt.Sprintf("%s => %s", err, insertPath))
					screen.refreshDatabase()
					screen.mode = modeBrowse
					screen.inputModal.Clear()
				} else {
					if parentB != nil {
						parentB.expanded = true
					}
					screen.currentPath = append(insertPath, newVal)
					screen.refreshDatabase()
					screen.startEditItem()
				}
			}
		}
	}
	return BrowserScreenIndex
}

func (screen *BrowserScreen) handleExportKeyEvent(event termbox.Event) int {
	if event.Key == termbox.KeyEsc {
		screen.mode = modeBrowse
		screen.inputModal.Clear()
	} else {
		screen.inputModal.HandleEvent(event)
		if screen.inputModal.IsDone() {
			b, p, _ := screen.db.getGenericFromPath(screen.currentPath)
			fileName := screen.inputModal.GetValue()
			if screen.mode&modeExportValue == modeExportValue {
				// Exporting the value
				if p != nil {
					if err := exportValue(screen.currentPath, fileName); err != nil {
						//screen.setMessage("Error Exporting to file " + fileName + ".")
						screen.setMessage(err.Error())
					} else {
						screen.setMessage("Value exported to file: " + fileName)
					}
				}
			} else if screen.mode&modeExportJSON == modeExportJSON {
				if b != nil || p != nil {
					if exportJSON(screen.currentPath, fileName) != nil {
						screen.setMessage("Error Exporting to file " + fileName + ".")
					} else {
						screen.setMessage("Value exported to file: " + fileName)
					}
				}
			}
			screen.mode = modeBrowse
			screen.inputModal.Clear()
		}
	}
	return BrowserScreenIndex
}

func (screen *BrowserScreen) jumpCursorUp(distance int) bool {
	// Jump up 'distance' lines
	visPaths, err := screen.db.buildVisiblePathSlice(screen.filter)
	if err == nil {
		findPath := screen.currentPath
		for idx, pth := range visPaths {
			startJump := true
			for i := range pth {
				if len(screen.currentPath) > i && pth[i] != screen.currentPath[i] {
					startJump = false
				}
			}
			if startJump {
				distance--
				if distance == 0 {
					screen.currentPath = visPaths[len(visPaths)-1-idx]
					break
				}
			}
		}
		isCurPath := true
		for i := range screen.currentPath {
			if screen.currentPath[i] != findPath[i] {
				isCurPath = false
				break
			}
		}
		if isCurPath {
			screen.currentPath = screen.db.getNextVisiblePath(nil, screen.filter)
		}
	}
	return true
}
func (screen *BrowserScreen) jumpCursorDown(distance int) bool {
	visPaths, err := screen.db.buildVisiblePathSlice(screen.filter)
	if err == nil {
		findPath := screen.currentPath
		for idx, pth := range visPaths {
			startJump := true

			for i := range pth {
				if len(screen.currentPath) > i && pth[i] != screen.currentPath[i] {
					startJump = false
				}
			}
			if startJump {
				distance--
				if distance == 0 {
					screen.currentPath = visPaths[idx]
					break
				}
			}
		}
		isCurPath := true
		for i := range screen.currentPath {
			if screen.currentPath[i] != findPath[i] {
				isCurPath = false
				break
			}
		}
		if isCurPath {
			screen.currentPath = screen.db.getNextVisiblePath(nil, screen.filter)
		}
	}
	return true
}

func (screen *BrowserScreen) moveCursorUp() bool {
	newPath := screen.db.getPrevVisiblePath(screen.currentPath, screen.filter)
	if newPath != nil {
		screen.currentPath = newPath
		return true
	}
	return false
}
func (screen *BrowserScreen) moveCursorDown() bool {
	newPath := screen.db.getNextVisiblePath(screen.currentPath, screen.filter)
	if newPath != nil {
		screen.currentPath = newPath
		return true
	}
	return false
}
func (screen *BrowserScreen) moveRightPaneUp() bool {
	if screen.rightViewPort.scrollRow > 0 {
		screen.rightViewPort.scrollRow--
		return true
	}
	return false
}
func (screen *BrowserScreen) moveRightPaneDown() bool {
	screen.rightViewPort.scrollRow++
	return true
}

func (screen *BrowserScreen) performLayout() {}

func (screen *BrowserScreen) drawScreen(style Style) {
	if screen.db == nil {
		screen.drawHeader(style)
		screen.setMessage("Invalid DB. Press 'q' to quit, '?' for help")
		screen.drawFooter(style)
		return
	}
	if len(screen.db.buckets) == 0 && screen.mode&modeInsertBucket != modeInsertBucket {
		// Force a bucket insert
		screen.startInsertItemAtParent(typeBucket)
	}
	if screen.message == "" {
		screen.setMessageWithTimeout("Press '?' for help", -1)
	}
	screen.drawLeftPane(style)
	screen.drawRightPane(style)
	screen.drawHeader(style)
	screen.drawFooter(style)

	if screen.inputModal != nil {
		screen.inputModal.Draw()
	}
	if screen.mode == modeDelete {
		screen.confirmModal.Draw()
	}
}

func (screen *BrowserScreen) drawHeader(style Style) {
	width, _ := termbox.Size()
	headerStringLen := func(fileName string) int {
		return len(ProgramName) + len(fileName) + 1
	}
	headerFileName := currentFilename
	if headerStringLen(headerFileName) > width {
		headerFileName = filepath.Base(headerFileName)
	}
	headerString := ProgramName + ": " + headerFileName
	count := ((width-len(headerString))/2)+1
	if count < 0 {
		count = 0
	}
	spaces := strings.Repeat(" ", count)
	termboxUtil.DrawStringAtPoint(fmt.Sprintf("%s%s%s", spaces, headerString, spaces), 0, 0, style.titleFg, style.titleBg)
}

func (screen *BrowserScreen) drawFooter(style Style) {
	if screen.messageTimeout > 0 && time.Since(screen.messageTime) > screen.messageTimeout {
		screen.clearMessage()
	}
	_, height := termbox.Size()
	termboxUtil.DrawStringAtPoint(screen.message, 0, height-1, style.defaultFg, style.defaultBg)
}

func (screen *BrowserScreen) buildLeftPane(style Style) {
	screen.leftPaneBuffer = nil
	if len(screen.currentPath) == 0 {
		screen.currentPath = screen.db.getNextVisiblePath(nil, screen.filter)
	}
	for i := range screen.db.buckets {
		screen.leftPaneBuffer = append(screen.leftPaneBuffer, screen.bucketToLines(&screen.db.buckets[i], style)...)
	}
	// Find the cursor in the leftPane
	for k, v := range screen.leftPaneBuffer {
		if v.Fg == style.cursorFg {
			screen.leftViewPort.scrollRow = k
			break
		}
	}
}

func (screen *BrowserScreen) drawLeftPane(style Style) {
	screen.buildLeftPane(style)
	w, h := termbox.Size()
	if w > 80 {
		w = w / 2
	}
	screen.leftViewPort.bytesPerRow = w
	screen.leftViewPort.numberOfRows = h - 2
	termboxUtil.FillWithChar('=', 0, 1, w, 1, style.defaultFg, style.defaultBg)
	startX, startY := 0, 3
	screen.leftViewPort.firstRow = startY
	treeOffset := 0
	maxCursor := screen.leftViewPort.numberOfRows * 2 / 3

	if screen.leftViewPort.scrollRow > maxCursor {
		treeOffset = screen.leftViewPort.scrollRow - maxCursor
	}
	if len(screen.leftPaneBuffer) > 0 {
		for k, v := range screen.leftPaneBuffer[treeOffset:] {
			termboxUtil.DrawStringAtPoint(v.Text, startX, (startY + k - 1), v.Fg, v.Bg)
		}
	}
}

func (screen *BrowserScreen) buildRightPane(style Style) {
	screen.rightPaneBuffer = nil
	b, p, err := screen.db.getGenericFromPath(screen.currentPath)
	if err == nil {
		if b != nil {
			screen.rightPaneBuffer = append(screen.rightPaneBuffer,
				Line{fmt.Sprintf("Path: %s", strings.Join(stringifyPath(b.GetPath()), " → ")), style.defaultFg, style.defaultBg})
			screen.rightPaneBuffer = append(screen.rightPaneBuffer,
				Line{fmt.Sprintf("Buckets: %d", len(b.buckets)), style.defaultFg, style.defaultBg})
			screen.rightPaneBuffer = append(screen.rightPaneBuffer,
				Line{fmt.Sprintf("Pairs: %d", len(b.pairs)), style.defaultFg, style.defaultBg})
		} else if p != nil {
			screen.rightPaneBuffer = append(screen.rightPaneBuffer,
				Line{fmt.Sprintf("Path: %s", strings.Join(stringifyPath(p.GetPath()), " → ")), style.defaultFg, style.defaultBg})
			screen.rightPaneBuffer = append(screen.rightPaneBuffer,
				Line{fmt.Sprintf("Key: %s", stringify([]byte(p.key))), style.defaultFg, style.defaultBg})

			value := strings.Split(string(formatValue([]byte(p.val))), "\n")
			if len(value) == 1 {
				screen.rightPaneBuffer = append(screen.rightPaneBuffer,
					Line{fmt.Sprintf("Value: %s", value[0]), style.defaultFg, style.defaultBg})
			} else {
				screen.rightPaneBuffer = append(screen.rightPaneBuffer,
					Line{"Value:", style.defaultFg, style.defaultBg})
				for _, v := range value {
					screen.rightPaneBuffer = append(screen.rightPaneBuffer,
						Line{v, style.defaultFg, style.defaultBg})
				}
			}
		}
	} else {
		screen.rightPaneBuffer = append(screen.rightPaneBuffer,
			Line{fmt.Sprintf("Path: %s", strings.Join(stringifyPath(screen.currentPath), " → ")), style.defaultFg, style.defaultBg})
		screen.rightPaneBuffer = append(screen.rightPaneBuffer,
			Line{err.Error(), termbox.ColorRed, termbox.ColorBlack})
	}
}

func (screen *BrowserScreen) drawRightPane(style Style) {
	screen.buildRightPane(style)
	w, h := termbox.Size()
	if w > 80 {
		screen.rightViewPort.bytesPerRow = w / 2
		screen.rightViewPort.numberOfRows = h - 2
		// Screen is wide enough, split it
		termboxUtil.FillWithChar('=', 0, 1, w, 1, style.defaultFg, style.defaultBg)
		termboxUtil.FillWithChar('|', (w / 2), screen.rightViewPort.firstRow-1, (w / 2), h, style.defaultFg, style.defaultBg)
		// Clear the right pane
		termboxUtil.FillWithChar(' ', (w/2)+1, screen.rightViewPort.firstRow+2, w, h, style.defaultFg, style.defaultBg)

		startX := (w / 2) + 2
		startY := 3
		maxScroll := len(screen.rightPaneBuffer) - screen.rightViewPort.numberOfRows
		if maxScroll < 0 {
			maxScroll = 0
		}
		if screen.rightViewPort.scrollRow > maxScroll {
			screen.rightViewPort.scrollRow = maxScroll
		}
		if len(screen.rightPaneBuffer) > 0 {
			for k, v := range screen.rightPaneBuffer[screen.rightViewPort.scrollRow:] {
				termboxUtil.DrawStringAtPoint(v.Text, startX, (startY + k - 1), v.Fg, v.Bg)
			}
		}
	}
}

func formatValue(val []byte) []byte {
	// Attempt JSON parsing and formatting
	out, err := formatValueJSON(val)
	if err == nil {
		return out
	}
	return []byte(stringify([]byte(val)))
}

func formatValueJSON(val []byte) ([]byte, error) {
	var jsonOut interface{}
	err := json.Unmarshal(val, &jsonOut)
	if err != nil {
		return val, err
	}
	out, err := json.MarshalIndent(jsonOut, "", "  ")
	if err != nil {
		return val, err
	}
	return out, nil
}

func (screen *BrowserScreen) bucketToLines(bkt *BoltBucket, style Style) []Line {
	var ret []Line
	bfg, bbg := style.defaultFg, style.defaultBg
	if comparePaths(screen.currentPath, bkt.GetPath()) {
		bfg, bbg = style.cursorFg, style.cursorBg
	}
	bktPrefix := strings.Repeat(" ", len(bkt.GetPath())*2)
	if bkt.expanded {
		ret = append(ret, Line{bktPrefix + "- " + stringify([]byte(bkt.name)), bfg, bbg})
		for i := range bkt.buckets {
			ret = append(ret, screen.bucketToLines(&bkt.buckets[i], style)...)
		}
		for _, bp := range bkt.pairs {
			if screen.filter != "" && !strings.Contains(bp.key, screen.filter) {
				continue
			}
			pfg, pbg := style.defaultFg, style.defaultBg
			if comparePaths(screen.currentPath, bp.GetPath()) {
				pfg, pbg = style.cursorFg, style.cursorBg
			}
			prPrefix := strings.Repeat(" ", len(bp.GetPath())*2)
			var pairString string
			if AppArgs.NoValue {
				pairString = fmt.Sprintf("%s%s", prPrefix, stringify([]byte(bp.key)))
			} else {
				pairString = fmt.Sprintf("%s%s: %s", prPrefix, stringify([]byte(bp.key)), stringify([]byte(bp.val)))
			}
			ret = append(ret, Line{pairString, pfg, pbg})
		}
	} else {
		ret = append(ret, Line{bktPrefix + "+ " + stringify([]byte(bkt.name)), bfg, bbg})
	}
	return ret
}

func (screen *BrowserScreen) startDeleteItem() bool {
	b, p, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil {
		w, h := termbox.Size()
		inpW, inpH := (w / 2), 6
		inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
		mod := termboxUtil.CreateConfirmModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
		if b != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Delete Bucket '%s'?", b.name), inpW-1, termboxUtil.AlignCenter))
		} else if p != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Delete Pair '%s'?", p.key), inpW-1, termboxUtil.AlignCenter))
		}
		mod.Show()
		mod.SetText(termboxUtil.AlignText("This cannot be undone!", inpW-1, termboxUtil.AlignCenter))
		screen.confirmModal = mod
		screen.mode = modeDelete
		return true
	}
	return false
}

func (screen *BrowserScreen) startFilter() bool {
	_, _, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil {
		w, h := termbox.Size()
		inpW, inpH := (w / 2), 6
		inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
		mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
		mod.SetTitle(termboxUtil.AlignText("Filter", inpW, termboxUtil.AlignCenter))
		mod.SetValue(screen.filter)
		mod.Show()
		screen.inputModal = mod
		screen.mode = modeFilter
		return true
	}
	return false
}

func (screen *BrowserScreen) startEditItem() bool {
	_, p, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil {
		w, h := termbox.Size()
		inpW, inpH := (w / 2), 6
		inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
		mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
		if p != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Input new value for '%s'", p.key), inpW, termboxUtil.AlignCenter))
			mod.SetValue(p.val)
		}
		mod.Show()
		screen.inputModal = mod
		screen.mode = modeChangeVal
		return true
	}
	return false
}

func (screen *BrowserScreen) startRenameItem() bool {
	b, p, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil {
		w, h := termbox.Size()
		inpW, inpH := (w / 2), 6
		inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
		mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
		if b != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Rename Bucket '%s' to:", b.name), inpW, termboxUtil.AlignCenter))
			mod.SetValue(b.name)
		} else if p != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Rename Key '%s' to:", p.key), inpW, termboxUtil.AlignCenter))
			mod.SetValue(p.key)
		}
		mod.Show()
		screen.inputModal = mod
		screen.mode = modeChangeKey
		return true
	}
	return false
}

func (screen *BrowserScreen) startInsertItemAtParent(tp BoltType) bool {
	w, h := termbox.Size()
	inpW, inpH := w-1, 7
	if w > 80 {
		inpW, inpH = (w / 2), 7
	}
	inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
	mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
	screen.inputModal = mod
	if len(screen.currentPath) <= 0 {
		// in the root directory
		if tp == typeBucket {
			mod.SetTitle(termboxUtil.AlignText("Create Root Bucket", inpW, termboxUtil.AlignCenter))
			screen.mode = modeInsertBucket | modeModToParent
			mod.Show()
			return true
		}
	} else {
		var insPath string
		_, p, e := screen.db.getGenericFromPath(screen.currentPath[:len(screen.currentPath)-1])
		if e == nil && p != nil {
			insPath = strings.Join(screen.currentPath[:len(screen.currentPath)-2], " → ") + " → "
		} else {
			insPath = strings.Join(screen.currentPath[:len(screen.currentPath)-1], " → ") + " → "
		}
		titlePrfx := ""
		if tp == typeBucket {
			titlePrfx = "New Bucket: "
		} else if tp == typePair {
			titlePrfx = "New Pair: "
		}
		titleText := titlePrfx + insPath
		if len(titleText) > inpW {
			truncW := len(titleText) - inpW
			titleText = titlePrfx + "..." + insPath[truncW+3:]
		}
		if tp == typeBucket {
			mod.SetTitle(termboxUtil.AlignText(titleText, inpW, termboxUtil.AlignCenter))
			screen.mode = modeInsertBucket | modeModToParent
			mod.Show()
			return true
		} else if tp == typePair {
			mod.SetTitle(termboxUtil.AlignText(titleText, inpW, termboxUtil.AlignCenter))
			mod.Show()
			screen.mode = modeInsertPair | modeModToParent
			return true
		}
	}
	return false
}

func (screen *BrowserScreen) startInsertItem(tp BoltType) bool {
	w, h := termbox.Size()
	inpW, inpH := w-1, 7
	if w > 80 {
		inpW, inpH = (w / 2), 7
	}
	inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
	mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
	//mod.SetInputWrap(true)
	screen.inputModal = mod
	var insPath string
	_, p, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil && p != nil {
		insPath = strings.Join(screen.currentPath[:len(screen.currentPath)-1], " → ") + " → "
	} else {
		insPath = strings.Join(screen.currentPath, " → ") + " → "
	}
	titlePrfx := ""
	if tp == typeBucket {
		titlePrfx = "New Bucket: "
	} else if tp == typePair {
		titlePrfx = "New Pair: "
	}
	titleText := titlePrfx + insPath
	if len(titleText) > inpW {
		truncW := len(titleText) - inpW
		titleText = titlePrfx + "..." + insPath[truncW+3:]
	}
	if tp == typeBucket {
		mod.SetTitle(termboxUtil.AlignText(titleText, inpW, termboxUtil.AlignCenter))
		screen.mode = modeInsertBucket
		mod.Show()
		return true
	} else if tp == typePair {
		mod.SetTitle(termboxUtil.AlignText(titleText, inpW, termboxUtil.AlignCenter))
		mod.Show()
		screen.mode = modeInsertPair
		return true
	}
	return false
}

func (screen *BrowserScreen) startExportValue() bool {
	_, p, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil && p != nil {
		w, h := termbox.Size()
		inpW, inpH := (w / 2), 6
		inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
		mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
		mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Export value of '%s' to:", p.key), inpW, termboxUtil.AlignCenter))
		mod.SetValue("")
		mod.Show()
		screen.inputModal = mod
		screen.mode = modeExportValue
		return true
	}
	screen.setMessage("Couldn't do string export on " + screen.currentPath[len(screen.currentPath)-1] + "(did you mean 'X'?)")
	return false
}

func (screen *BrowserScreen) startExportJSON() bool {
	b, p, e := screen.db.getGenericFromPath(screen.currentPath)
	if e == nil {
		w, h := termbox.Size()
		inpW, inpH := (w / 2), 6
		inpX, inpY := ((w / 2) - (inpW / 2)), ((h / 2) - inpH)
		mod := termboxUtil.CreateInputModal("", inpX, inpY, inpW, inpH, termbox.ColorWhite, termbox.ColorBlack)
		if b != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Export JSON of '%s' to:", b.name), inpW, termboxUtil.AlignCenter))
			mod.SetValue("")
		} else if p != nil {
			mod.SetTitle(termboxUtil.AlignText(fmt.Sprintf("Export JSON of '%s' to:", p.key), inpW, termboxUtil.AlignCenter))
			mod.SetValue("")
		}
		mod.Show()
		screen.inputModal = mod
		screen.mode = modeExportJSON
		return true
	}
	return false
}

func (screen *BrowserScreen) setMessage(msg string) {
	screen.message = msg
	screen.messageTime = time.Now()
	screen.messageTimeout = time.Second * 2
}

/* setMessageWithTimeout lets you specify the timeout for the message
 * setting it to -1 means it won't timeout
 */
func (screen *BrowserScreen) setMessageWithTimeout(msg string, timeout time.Duration) {
	screen.message = msg
	screen.messageTime = time.Now()
	screen.messageTimeout = timeout
}

func (screen *BrowserScreen) clearMessage() {
	screen.message = ""
	screen.messageTimeout = -1
}

func (screen *BrowserScreen) refreshDatabase() {
	shadowDB := screen.db
	screen.db = screen.db.refreshDatabase()
	screen.db.syncOpenBuckets(shadowDB)
}

func comparePaths(p1, p2 []string) bool {
	return strings.Join(p1, " → ") == strings.Join(p2, " → ")
}
