package termboxUtil

import "github.com/nsf/termbox-go"

// ScrollFrame is a frame for holding other elements
// It manages it's own x/y, tab index
type ScrollFrame struct {
	id                  string
	x, y, width, height int
	scrollX, scrollY    int
	tabIdx              int
	fg, bg              termbox.Attribute
	bordered            bool
	controls            []termboxControl
}

// CreateScrollFrame creates Scrolling Frame at x, y that is w by h
func CreateScrollFrame(x, y, w, h int, fg, bg termbox.Attribute) *ScrollFrame {
	s := ScrollFrame{x: x, y: y, width: w, height: h, fg: fg, bg: bg}
	return &s
}

// GetID returns this control's ID
func (i *ScrollFrame) GetID() string { return i.id }

// SetID sets this control's ID
func (i *ScrollFrame) SetID(newID string) {
	i.id = newID
}

// GetX returns the x position of the scroll frame
func (i *ScrollFrame) GetX() int { return i.x }

// SetX sets the x position of the scroll frame
func (i *ScrollFrame) SetX(x int) {
	i.x = x
}

// GetY returns the y position of the scroll frame
func (i *ScrollFrame) GetY() int { return i.y }

// SetY sets the y position of the scroll frame
func (i *ScrollFrame) SetY(y int) {
	i.y = y
}

// GetWidth returns the current width of the scroll frame
func (i *ScrollFrame) GetWidth() int { return i.width }

// SetWidth sets the current width of the scroll frame
func (i *ScrollFrame) SetWidth(w int) {
	i.width = w
}

// GetHeight returns the current height of the scroll frame
func (i *ScrollFrame) GetHeight() int { return i.height }

// SetHeight sets the current height of the scroll frame
func (i *ScrollFrame) SetHeight(h int) {
	i.height = h
}

// GetFgColor returns the foreground color
func (i *ScrollFrame) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *ScrollFrame) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *ScrollFrame) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *ScrollFrame) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// IsBordered returns true or false if this scroll frame has a border
func (i *ScrollFrame) IsBordered() bool { return i.bordered }

// SetBordered sets whether we render a border around the scroll frame
func (i *ScrollFrame) SetBordered(b bool) {
	i.bordered = b
}

// GetScrollX returns the x distance scrolled
func (i *ScrollFrame) GetScrollX() int {
	return i.scrollX
}

// GetScrollY returns the y distance scrolled
func (i *ScrollFrame) GetScrollY() int {
	return i.scrollY
}

// ScrollDown scrolls the frame down
func (i *ScrollFrame) ScrollDown() {
	i.scrollY++
}

// ScrollUp scrolls the frame up
func (i *ScrollFrame) ScrollUp() {
	if i.scrollY > 0 {
		i.scrollY--
	}
}

// ScrollLeft scrolls the frame left
func (i *ScrollFrame) ScrollLeft() {
	if i.scrollX > 0 {
		i.scrollX--
	}
}

// ScrollRight scrolls the frame right
func (i *ScrollFrame) ScrollRight() {
	i.scrollX++
}

// AddControl adds a control to the frame
func (i *ScrollFrame) AddControl(t termboxControl) {
	i.controls = append(i.controls, t)
}

// DrawControl figures out the relative position of the control,
// sets it, draws it, then resets it.
func (i *ScrollFrame) DrawControl(t termboxControl) {
	if i.IsVisible(t) {
		ctlX, ctlY := t.GetX(), t.GetY()
		t.SetX((i.GetX() + ctlX))
		t.SetY((i.GetY() + ctlY))
		t.Draw()
		t.SetX(ctlX)
		t.SetY(ctlY)
	}
}

// IsVisible takes a Termbox Control and returns whether
// that control would be visible in the frame
func (i *ScrollFrame) IsVisible(t termboxControl) bool {
	// Check if any part of t should be visible
	cX, cY := t.GetX(), t.GetY()
	if cX+t.GetWidth() >= i.scrollX && cX <= i.scrollX+i.width {
		return cY+t.GetHeight() >= i.scrollY && cY <= i.scrollY+i.height
	}
	return false
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (i *ScrollFrame) HandleEvent(event termbox.Event) bool {
	return false
}

// DrawToStrings generates a slice of strings with what should
// be drawn to the screen
func (i *ScrollFrame) DrawToStrings() []string {
	return []string{}
}

// Draw outputs the Scoll Frame on the screen
func (i *ScrollFrame) Draw() {
	maxWidth := i.width
	maxHeight := i.height
	x, y := i.x, i.y
	startX := i.x
	startY := i.y
	if i.bordered {
		DrawBorder(i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)
		maxWidth--
		maxHeight--
		x++
		y++
		startX++
		startY++
	}
	for idx := range i.controls {
		i.DrawControl(i.controls[idx])
	}
}
