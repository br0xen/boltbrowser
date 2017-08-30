package termboxUtil

import "github.com/nsf/termbox-go"

// ProgressBar Just contains the data needed to display a progress bar
type ProgressBar struct {
	id             string
	total          int
	progress       int
	allowOverflow  bool
	allowUnderflow bool
	fullChar       rune
	emptyChar      rune
	bordered       bool
	alignment      TextAlignment
	colorized      bool

	x, y          int
	width, height int
	bg, fg        termbox.Attribute
}

// CreateProgressBar Create a progress bar object
func CreateProgressBar(tot, x, y int, fg, bg termbox.Attribute) *ProgressBar {
	i := ProgressBar{total: tot,
		fullChar: '#', emptyChar: ' ',
		x: x, y: y, height: 1, width: 10,
		bordered: true, fg: fg, bg: bg,
		alignment: AlignLeft,
	}
	return &i
}

// GetID returns this control's ID
func (i *ProgressBar) GetID() string { return i.id }

// SetID sets this control's ID
func (i *ProgressBar) SetID(newID string) {
	i.id = newID
}

// GetProgress returns the curret progress value
func (i *ProgressBar) GetProgress() int {
	return i.progress
}

// SetProgress sets the current progress of the bar
func (i *ProgressBar) SetProgress(p int) {
	if (p <= i.total || i.allowOverflow) || (p >= 0 || i.allowUnderflow) {
		i.progress = p
	}
}

// IncrProgress increments the current progress of the bar
func (i *ProgressBar) IncrProgress() {
	if i.progress < i.total || i.allowOverflow {
		i.progress++
	}
}

// DecrProgress decrements the current progress of the bar
func (i *ProgressBar) DecrProgress() {
	if i.progress > 0 || i.allowUnderflow {
		i.progress--
	}
}

// GetPercent returns the percent full of the bar
func (i *ProgressBar) GetPercent() int {
	return int(float64(i.progress) / float64(i.total) * 100)
}

// EnableOverflow Tells the progress bar that it can go over the total
func (i *ProgressBar) EnableOverflow() {
	i.allowOverflow = true
}

// DisableOverflow Tells the progress bar that it can NOT go over the total
func (i *ProgressBar) DisableOverflow() {
	i.allowOverflow = false
}

// EnableUnderflow Tells the progress bar that it can go below zero
func (i *ProgressBar) EnableUnderflow() {
	i.allowUnderflow = true
}

// DisableUnderflow Tells the progress bar that it can NOT go below zero
func (i *ProgressBar) DisableUnderflow() {
	i.allowUnderflow = false
}

// GetFullChar returns the rune used for 'full'
func (i *ProgressBar) GetFullChar() rune {
	return i.fullChar
}

// SetFullChar sets the rune used for 'full'
func (i *ProgressBar) SetFullChar(f rune) {
	i.fullChar = f
}

// GetEmptyChar gets the rune used for 'empty'
func (i *ProgressBar) GetEmptyChar() rune {
	return i.emptyChar
}

// SetEmptyChar sets the rune used for 'empty'
func (i *ProgressBar) SetEmptyChar(f rune) {
	i.emptyChar = f
}

// GetX Return the x position of the Progress Bar
func (i *ProgressBar) GetX() int { return i.x }

// SetX set the x position of the ProgressBar to x
func (i *ProgressBar) SetX(x int) {
	i.x = x
}

// GetY Return the y position of the ProgressBar
func (i *ProgressBar) GetY() int { return i.y }

// SetY Set the y position of the ProgressBar to y
func (i *ProgressBar) SetY(y int) {
	i.y = y
}

// GetHeight returns the height of the progress bar
// Defaults to 1 (3 if bordered)
func (i *ProgressBar) GetHeight() int {
	return i.height
}

// SetHeight Sets the height of the progress bar
func (i *ProgressBar) SetHeight(h int) {
	i.height = h
}

// GetWidth returns the width of the progress bar
func (i *ProgressBar) GetWidth() int {
	return i.width
}

// SetWidth Sets the width of the progress bar
func (i *ProgressBar) SetWidth(w int) {
	i.width = w
}

// GetFgColor returns the foreground color
func (i *ProgressBar) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *ProgressBar) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *ProgressBar) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *ProgressBar) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// Align Tells which direction the progress bar empties
func (i *ProgressBar) Align(a TextAlignment) {
	i.alignment = a
}

// SetColorized sets whether the progress bar should be colored
// depending on how full it is:
//  10% - Red
//	50% - Yellow
//	80% - Green
func (i *ProgressBar) SetColorized(c bool) {
	i.colorized = c
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (i *ProgressBar) HandleEvent(event termbox.Event) bool {
	return false
}

// Draw outputs the input field on the screen
func (i *ProgressBar) Draw() {
	// For now, just draw a [####  ] bar
	// TODO: make this more advanced
	useFg := i.fg
	if i.colorized {
		if i.GetPercent() < 10 {
			useFg = termbox.ColorRed
		} else if i.GetPercent() < 50 {
			useFg = termbox.ColorYellow
		} else {
			useFg = termbox.ColorGreen
		}
	}
	drawX, drawY := i.x, i.y
	fillWidth, fillHeight := i.width-2, i.height
	DrawStringAtPoint("[", drawX, drawY, i.fg, i.bg)
	numFull := int(float64(fillWidth) * float64(i.progress) / float64(i.total))
	FillWithChar(i.fullChar, drawX+1, drawY, drawX+1+numFull, drawY+(fillHeight-1), useFg, i.bg)
	DrawStringAtPoint("]", drawX+i.width-1, drawY, i.fg, i.bg)

	/*
		drawX, drawY := i.x, i.y
		drawWidth, drawHeight := i.width, i.height
		if i.bordered {
			if i.height == 1 && i.width > 2 {
				// Just using [ & ] for the border
				DrawStringAtPoint("[", drawX, drawY, i.fg, i.bg)
				DrawStringAtPoint("]", drawX+i.width-1, drawY, i.fg, i.bg)
				drawX++
				drawWidth -= 2
			} else if i.height >= 3 {
				DrawBorder(drawX, drawY, drawX+i.width, drawY+i.height, i.fg, i.bg)
				drawX++
				drawY++
				drawWidth -= 2
				drawHeight -= 2
			}
		}

		// Figure out how many chars are full
		numFull := drawWidth * (i.progress / i.total)
		switch i.alignment {
		case AlignRight: // TODO: Fill from right to left
		case AlignCenter: // TODO: Fill from middle out
		default: // Fill from left to right
			FillWithChar(i.fullChar, drawX, drawY, drawX+numFull, drawY+(drawHeight-1), i.fg, i.bg)
			if numFull < drawWidth {
				FillWithChar(i.emptyChar, drawX+numFull, drawY, drawX+drawWidth-1, drawY+(drawHeight-1), i.fg, i.bg)
			}
		}
	*/
}
