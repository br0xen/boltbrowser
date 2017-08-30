package termboxUtil

import (
	"github.com/nsf/termbox-go"
)

// InputModal A modal for text input
type InputModal struct {
	id                  string
	title               string
	text                string
	input               *InputField
	x, y, width, height int
	showHelp            bool
	cursor              int
	bg, fg              termbox.Attribute
	isDone              bool
	isAccepted          bool
	isVisible           bool
	bordered            bool
	tabSkip             bool
	inputSelected       bool
}

// CreateInputModal Create an input modal with the given attributes
func CreateInputModal(title string, x, y, width, height int, fg, bg termbox.Attribute) *InputModal {
	i := InputModal{title: title, x: x, y: y, width: width, height: height, fg: fg, bg: bg, bordered: true}
	i.input = CreateInputField(i.x+2, i.y+3, i.width-2, 2, i.fg, i.bg)
	i.showHelp = true
	i.input.bordered = true
	i.isVisible = true
	i.inputSelected = true
	return &i
}

// GetID returns this control's ID
func (i *InputModal) GetID() string { return i.id }

// SetID sets this control's ID
func (i *InputModal) SetID(newID string) {
	i.id = newID
}

// GetTitle Return the title of the modal
func (i *InputModal) GetTitle() string { return i.title }

// SetTitle Sets the title of the modal to s
func (i *InputModal) SetTitle(s string) {
	i.title = s
}

// GetText Return the text of the modal
func (i *InputModal) GetText() string { return i.text }

// SetText Set the text of the modal to s
func (i *InputModal) SetText(s string) {
	i.text = s
}

// GetX Return the x position of the modal
func (i *InputModal) GetX() int { return i.x }

// SetX set the x position of the modal to x
func (i *InputModal) SetX(x int) {
	i.x = x
}

// GetY Return the y position of the modal
func (i *InputModal) GetY() int { return i.y }

// SetY Set the y position of the modal to y
func (i *InputModal) SetY(y int) {
	i.y = y
}

// GetWidth Return the width of the modal
func (i *InputModal) GetWidth() int { return i.width }

// SetWidth Set the width of the modal to width
func (i *InputModal) SetWidth(width int) {
	i.width = width
}

// GetHeight Return the height of the modal
func (i *InputModal) GetHeight() int { return i.height }

// SetHeight Set the height of the modal to height
func (i *InputModal) SetHeight(height int) {
	i.height = height
}

// SetMultiline returns whether this is a multiline modal
func (i *InputModal) SetMultiline(m bool) {
	i.input.multiline = m
}

// IsMultiline returns whether this is a multiline modal
func (i *InputModal) IsMultiline() bool {
	return i.input.multiline
}

// IsBordered returns whether this control is bordered or not
func (i *InputModal) IsBordered() bool {
	return i.bordered
}

// SetBordered sets whether we render a border around the frame
func (i *InputModal) SetBordered(b bool) {
	i.bordered = b
}

// IsTabSkipped returns whether this control has it's tabskip flag set
func (i *InputModal) IsTabSkipped() bool {
	return i.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (i *InputModal) SetTabSkip(b bool) {
	i.tabSkip = b
}

// HelpIsShown Returns whether the modal is showing it's help text or not
func (i *InputModal) HelpIsShown() bool { return i.showHelp }

// ShowHelp Set the "Show Help" flag
func (i *InputModal) ShowHelp(b bool) {
	i.showHelp = b
}

// GetFgColor returns the foreground color
func (i *InputModal) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *InputModal) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *InputModal) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *InputModal) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// Show Sets the visibility flag to true
func (i *InputModal) Show() {
	i.isVisible = true
}

// Hide Sets the visibility flag to false
func (i *InputModal) Hide() {
	i.isVisible = false
}

// IsVisible returns the isVisible flag
func (i *InputModal) IsVisible() bool {
	return i.isVisible
}

// SetDone Sets the flag that tells whether this modal has completed it's purpose
func (i *InputModal) SetDone(b bool) {
	i.isDone = b
}

// IsDone Returns the "isDone" flag
func (i *InputModal) IsDone() bool {
	return i.isDone
}

// IsAccepted Returns whether the modal has been accepted
func (i *InputModal) IsAccepted() bool {
	return i.isAccepted
}

// GetValue Return the current value of the input
func (i *InputModal) GetValue() string { return i.input.GetValue() }

// SetValue Sets the value of the input to s
func (i *InputModal) SetValue(s string) {
	i.input.SetValue(s)
}

// SetInputWrap sets whether the input field will wrap long text or not
func (i *InputModal) SetInputWrap(b bool) {
	i.input.SetWrap(b)
}

// Clear Resets all non-positional parameters of the modal
func (i *InputModal) Clear() {
	i.title = ""
	i.text = ""
	i.input.SetValue("")
	i.isDone = false
	i.isVisible = false
}

// HandleEvent Handle the termbox event, return true if it was consumed
func (i *InputModal) HandleEvent(event termbox.Event) bool {
	if event.Key == termbox.KeyEnter {
		if !i.input.IsMultiline() || !i.inputSelected {
			// Done editing
			i.isDone = true
			i.isAccepted = true
		} else {
			i.input.HandleEvent(event)
		}
		return true
	} else if event.Key == termbox.KeyTab {
		if i.input.IsMultiline() {
			i.inputSelected = !i.inputSelected
		}
	} else if event.Key == termbox.KeyEsc {
		// Done editing
		i.isDone = true
		i.isAccepted = false
		return true
	}
	return i.input.HandleEvent(event)
}

// Draw Draw the modal
func (i *InputModal) Draw() {
	if i.isVisible {
		// First blank out the area we'll be putting the modal
		FillWithChar(' ', i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)
		nextY := i.y + 1
		// The title
		if i.title != "" {
			if len(i.title) > i.width {
				diff := i.width - len(i.title)
				DrawStringAtPoint(i.title[:len(i.title)+diff-1], i.x+1, nextY, i.fg, i.bg)
			} else {
				DrawStringAtPoint(i.title, i.x+1, nextY, i.fg, i.bg)
			}
			nextY++
			FillWithChar('-', i.x+1, nextY, i.x+i.width-1, nextY, i.fg, i.bg)
			nextY++
		}
		if i.text != "" {
			DrawStringAtPoint(i.text, i.x+1, nextY, i.fg, i.bg)
			nextY++
		}
		i.input.SetY(nextY)
		i.input.Draw()
		nextY += 3
		if i.showHelp {
			helpString := " (ENTER) to Accept. (ESC) to Cancel. "
			helpX := (i.x + i.width - len(helpString)) - 1
			DrawStringAtPoint(helpString, helpX, nextY, i.fg, i.bg)
		}
		if i.bordered {
			// Now draw the border
			DrawBorder(i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)
		}
	}
}
