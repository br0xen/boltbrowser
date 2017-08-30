package termboxUtil

import "github.com/nsf/termbox-go"

// Label is a field for inputting text
type Label struct {
	id                  string
	value               string
	x, y, width, height int
	cursor              int
	fg, bg              termbox.Attribute
	bordered            bool
	wrap                bool
	multiline           bool
}

// CreateLabel creates an input field at x, y that is w by h
func CreateLabel(lbl string, x, y, w, h int, fg, bg termbox.Attribute) *Label {
	i := Label{value: lbl, x: x, y: y, width: w, height: h, fg: fg, bg: bg}
	return &i
}

// GetID returns this control's ID
func (i *Label) GetID() string { return i.id }

// SetID sets this control's ID
func (i *Label) SetID(newID string) {
	i.id = newID
}

// GetValue gets the current text that is in the Label
func (i *Label) GetValue() string { return i.value }

// SetValue sets the current text in the Label to s
func (i *Label) SetValue(s string) {
	i.value = s
}

// GetX returns the x position of the input field
func (i *Label) GetX() int { return i.x }

// SetX sets the x position of the input field
func (i *Label) SetX(x int) {
	i.x = x
}

// GetY returns the y position of the input field
func (i *Label) GetY() int { return i.y }

// SetY sets the y position of the input field
func (i *Label) SetY(y int) {
	i.y = y
}

// GetWidth returns the current width of the input field
func (i *Label) GetWidth() int {
	if i.width == -1 {
		if i.bordered {
			return len(i.value) + 2
		}
		return len(i.value)
	}
	return i.width
}

// SetWidth sets the current width of the input field
func (i *Label) SetWidth(w int) {
	i.width = w
}

// GetHeight returns the current height of the input field
func (i *Label) GetHeight() int { return i.height }

// SetHeight sets the current height of the input field
func (i *Label) SetHeight(h int) {
	i.height = h
}

// GetFgColor returns the foreground color
func (i *Label) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *Label) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *Label) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *Label) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// IsBordered returns true or false if this input field has a border
func (i *Label) IsBordered() bool { return i.bordered }

// SetBordered sets whether we render a border around the input field
func (i *Label) SetBordered(b bool) {
	i.bordered = b
}

// DoesWrap returns true or false if this input field wraps text
func (i *Label) DoesWrap() bool { return i.wrap }

// SetWrap sets whether we wrap the text at width.
func (i *Label) SetWrap(b bool) {
	i.wrap = b
}

// IsMultiline returns true or false if this field can have multiple lines
func (i *Label) IsMultiline() bool { return i.multiline }

// SetMultiline sets whether the field can have multiple lines
func (i *Label) SetMultiline(b bool) {
	i.multiline = b
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (i *Label) HandleEvent(event termbox.Event) bool { return false }

// Draw outputs the input field on the screen
func (i *Label) Draw() {
	maxWidth := i.width
	maxHeight := i.height
	x, y := i.x, i.y
	startX := i.x
	startY := i.y
	if i.bordered {
		DrawBorder(i.x, i.y, i.x+i.GetWidth(), i.y+i.height, i.fg, i.bg)
		maxWidth--
		maxHeight--
		x++
		y++
		startX++
		startY++
	}

	DrawStringAtPoint(i.value, x, y, i.fg, i.bg)
}
