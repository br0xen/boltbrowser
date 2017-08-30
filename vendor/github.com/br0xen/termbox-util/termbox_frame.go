package termboxUtil

import "github.com/nsf/termbox-go"

// Frame is a frame for holding other elements
// It manages it's own x/y, tab index
type Frame struct {
	id                  string
	x, y, width, height int
	tabIdx              int
	fg, bg              termbox.Attribute
	bordered            bool
	controls            []termboxControl
	tabSkip             bool
}

// CreateFrame creates a Frame at x, y that is w by h
func CreateFrame(x, y, w, h int, fg, bg termbox.Attribute) *Frame {
	s := Frame{x: x, y: y, width: w, height: h, fg: fg, bg: bg, bordered: true}
	return &s
}

// GetID returns this control's ID
func (i *Frame) GetID() string { return i.id }

// SetID sets this control's ID
func (i *Frame) SetID(newID string) {
	i.id = newID
}

// GetX returns the x position of the frame
func (i *Frame) GetX() int { return i.x }

// SetX sets the x position of the frame
func (i *Frame) SetX(x int) {
	i.x = x
}

// GetY returns the y position of the frame
func (i *Frame) GetY() int { return i.y }

// SetY sets the y position of the frame
func (i *Frame) SetY(y int) {
	i.y = y
}

// GetWidth returns the current width of the frame
func (i *Frame) GetWidth() int { return i.width }

// SetWidth sets the current width of the frame
func (i *Frame) SetWidth(w int) {
	i.width = w
}

// GetHeight returns the current height of the frame
func (i *Frame) GetHeight() int { return i.height }

// SetHeight sets the current height of the frame
func (i *Frame) SetHeight(h int) {
	i.height = h
}

// GetFgColor returns the foreground color
func (i *Frame) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *Frame) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *Frame) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *Frame) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// IsBordered returns true or false if this frame has a border
func (i *Frame) IsBordered() bool { return i.bordered }

// SetBordered sets whether we render a border around the frame
func (i *Frame) SetBordered(b bool) {
	i.bordered = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (i *Frame) IsTabSkipped() bool {
	return i.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (i *Frame) SetTabSkip(b bool) {
	i.tabSkip = b
}

// AddControl adds a control to the frame
func (i *Frame) AddControl(t termboxControl) {
	i.controls = append(i.controls, t)
}

// GetActiveControl returns the control at tabIdx
func (i *Frame) GetActiveControl() termboxControl {
	if len(i.controls) >= i.tabIdx {
		return i.controls[i.tabIdx]
	}
	return nil
}

// GetControls returns a slice of all controls
func (i *Frame) GetControls() []termboxControl {
	return i.controls
}

// GetControl returns the control at index i
func (i *Frame) GetControl(idx int) termboxControl {
	if len(i.controls) >= idx {
		return i.controls[idx]
	}
	return nil
}

// GetControlCount returns the number of controls contained
func (i *Frame) GetControlCount() int {
	return len(i.controls)
}

// GetLastControl returns the last control contained
func (i *Frame) GetLastControl() termboxControl {
	return i.controls[len(i.controls)-1]
}

// DrawControl figures out the relative position of the control,
// sets it, draws it, then resets it.
func (i *Frame) DrawControl(t termboxControl) {
	ctlX, ctlY := t.GetX(), t.GetY()
	t.SetX((i.GetX() + ctlX))
	t.SetY((i.GetY() + ctlY))
	t.Draw()
	t.SetX(ctlX)
	t.SetY(ctlY)
}

// GetBottomY returns the y of the lowest control in the frame
func (i *Frame) GetBottomY() int {
	var ret int
	for idx := range i.controls {
		if i.controls[idx].GetY()+i.controls[idx].GetHeight() > ret {
			ret = i.controls[idx].GetY() + i.controls[idx].GetHeight()
		}
	}
	return ret
}

// HandleEvent accepts the termbox event and returns whether it was consumed
func (i *Frame) HandleEvent(event termbox.Event) bool {
	if event.Key == termbox.KeyTab {
		i.FindNextTabStop()
		return true
	}
	return i.controls[i.tabIdx].HandleEvent(event)
}

// FindNextTabStop finds the next control that can be tabbed to
// A return of true means it found a different one than we started on.
func (i *Frame) FindNextTabStop() bool {
	startTab := i.tabIdx
	i.tabIdx = (i.tabIdx + 1) % len(i.controls)
	for i.controls[i.tabIdx].IsTabSkipped() {
		i.tabIdx = (i.tabIdx + 1) % len(i.controls)
		if i.tabIdx == startTab {
			break
		}
	}
	return i.tabIdx != startTab
}

// Draw outputs the Scoll Frame on the screen
func (i *Frame) Draw() {
	maxWidth := i.width
	maxHeight := i.height
	x, y := i.x, i.y
	startX := i.x
	startY := i.y
	if i.bordered {
		FillWithChar(' ', i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)
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
