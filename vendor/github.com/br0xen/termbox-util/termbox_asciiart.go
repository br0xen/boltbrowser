package termboxUtil

import (
	"strings"

	"github.com/nsf/termbox-go"
)

// ASCIIArt is a []string with more functions
type ASCIIArt struct {
	id       string
	contents []string
	x, y     int
	bg, fg   termbox.Attribute
	bordered bool
	tabSkip  bool
}

// CreateASCIIArt Create an ASCII art object from a string slice
func CreateASCIIArt(c []string, x, y int, fg, bg termbox.Attribute) *ASCIIArt {
	i := ASCIIArt{contents: c, x: x, y: y, fg: fg, bg: bg, bordered: false, tabSkip: true}
	return &i
}

// GetID returns this control's ID
func (i *ASCIIArt) GetID() string { return i.id }

// SetID sets this control's ID
func (i *ASCIIArt) SetID(newID string) {
	i.id = newID
}

// GetX Return the x position of the modal
func (i *ASCIIArt) GetX() int { return i.x }

// SetX set the x position of the modal to x
func (i *ASCIIArt) SetX(x int) {
	i.x = x
}

// GetY Return the y position of the modal
func (i *ASCIIArt) GetY() int { return i.y }

// SetY Set the y position of the modal to y
func (i *ASCIIArt) SetY(y int) {
	i.y = y
}

// GetHeight Returns the number of strings in the contents slice
func (i *ASCIIArt) GetHeight() int {
	return len(i.contents)
}

// SetHeight truncates lines from the bottom of the ascii art
func (i *ASCIIArt) SetHeight(h int) {
	if len(i.contents) > h {
		i.contents = i.contents[:h]
	} else {
		for j := len(i.contents); j < h; j++ {
			i.contents = append(i.contents, "")
		}
	}
}

// GetWidth Returns the number of strings in the contents slice
func (i *ASCIIArt) GetWidth() int {
	// Find the longest line
	var ret int
	for j := range i.contents {
		if len(i.contents[j]) > ret {
			ret = len(i.contents[j])
		}
	}
	return ret
}

// SetWidth Sets all lines in the contents to width w
func (i *ASCIIArt) SetWidth(w int) {
	// Find the longest line
	for j := range i.contents {
		mkUp := w - len(i.contents[j])
		if mkUp > 0 {
			i.contents[j] = i.contents[j] + strings.Repeat(" ", mkUp)
		} else {
			i.contents[j] = i.contents[j][:w]
		}
	}
}

// SetContents Sets the contents of i to c
func (i *ASCIIArt) SetContents(c []string) {
	i.contents = c
}

// GetContents returns the ascii art
func (i *ASCIIArt) GetContents() []string {
	return i.contents
}

// SetContentLine Sets a specific line of the contents to s
func (i *ASCIIArt) SetContentLine(s string, idx int) {
	if idx >= 0 && idx < len(i.contents) {
		i.contents[idx] = s
	}
}

// GetFgColor returns the foreground color
func (i *ASCIIArt) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *ASCIIArt) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *ASCIIArt) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *ASCIIArt) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// Align Align the Ascii art over width width with alignment a
func (i *ASCIIArt) Align(a TextAlignment, width int) {
	// First get the width of the longest string in the slice
	var newContents []string
	incomingLength := 0
	for _, line := range i.contents {
		if len(line) > incomingLength {
			incomingLength = len(line)
		}
	}
	for _, line := range i.contents {
		newContents = append(newContents, AlignText(AlignText(line, incomingLength, AlignLeft), width, a))
	}
	i.contents = newContents
}

// IsBordered returns whether this modal is bordered or not
func (i *ASCIIArt) IsBordered() bool {
	return i.bordered
}

// SetBordered sets whether we render a border around the frame
func (i *ASCIIArt) SetBordered(b bool) {
	i.bordered = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (i *ASCIIArt) IsTabSkipped() bool {
	return i.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (i *ASCIIArt) SetTabSkip(b bool) {
	i.tabSkip = b
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (i *ASCIIArt) HandleEvent(event termbox.Event) bool {
	return false
}

// Draw outputs the input field on the screen
func (i *ASCIIArt) Draw() {
	drawX, drawY := i.x, i.y
	for _, line := range i.contents {
		DrawStringAtPoint(line, drawX, drawY, i.fg, i.bg)
		drawY++
	}
}
