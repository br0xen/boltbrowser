package termboxUtil

import (
	"github.com/nsf/termbox-go"
)

// ConfirmModal is a modal with yes/no (or similar) buttons
type ConfirmModal struct {
	id                  string
	title               string
	text                string
	x, y, width, height int
	showHelp            bool
	cursor              int
	bg, fg              termbox.Attribute
	isDone              bool
	accepted            bool
	value               string
	isVisible           bool
	bordered            bool
	tabSkip             bool
}

// CreateConfirmModal Creates a confirmation modal with the specified attributes
func CreateConfirmModal(title string, x, y, width, height int, fg, bg termbox.Attribute) *ConfirmModal {
	i := ConfirmModal{title: title, x: x, y: y, width: width, height: height, fg: fg, bg: bg}
	if i.title == "" && i.text == "" {
		i.title = "Confirm?"
	}
	i.showHelp = true
	return &i
}

// GetID returns this control's ID
func (i *ConfirmModal) GetID() string { return i.id }

// SetID sets this control's ID
func (i *ConfirmModal) SetID(newID string) {
	i.id = newID
}

// GetTitle returns the current title of the modal
func (i *ConfirmModal) GetTitle() string { return i.title }

// SetTitle sets the current title of the modal to s
func (i *ConfirmModal) SetTitle(s string) {
	i.title = s
}

// GetText returns the current text of the modal
func (i *ConfirmModal) GetText() string { return i.text }

// SetText sets the text of the modal to s
func (i *ConfirmModal) SetText(s string) {
	i.text = s
}

// GetX returns the current x coordinate of the modal
func (i *ConfirmModal) GetX() int { return i.x }

// SetX sets the current x coordinate of the modal to x
func (i *ConfirmModal) SetX(x int) {
	i.x = x
}

// GetY returns the current y coordinate of the modal
func (i *ConfirmModal) GetY() int { return i.y }

// SetY sets the current y coordinate of the modal to y
func (i *ConfirmModal) SetY(y int) {
	i.y = y
}

// GetWidth returns the current width of the modal
func (i *ConfirmModal) GetWidth() int { return i.width }

// SetWidth sets the current modal width to width
func (i *ConfirmModal) SetWidth(width int) {
	i.width = width
}

// GetHeight returns the current height of the modal
func (i *ConfirmModal) GetHeight() int { return i.height }

// SetHeight set the height of the modal to height
func (i *ConfirmModal) SetHeight(height int) {
	i.height = height
}

// HelpIsShown returns true or false if the help is displayed
func (i *ConfirmModal) HelpIsShown() bool { return i.showHelp }

// ShowHelp sets whether or not to display the help text
func (i *ConfirmModal) ShowHelp(b bool) {
	i.showHelp = b
}

// GetFgColor returns the foreground color
func (i *ConfirmModal) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *ConfirmModal) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *ConfirmModal) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *ConfirmModal) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// IsDone returns whether the user has answered the modal
func (i *ConfirmModal) IsDone() bool { return i.isDone }

// SetDone sets whether the modal has completed it's purpose
func (i *ConfirmModal) SetDone(b bool) {
	i.isDone = b
}

// Show sets the visibility flag of the modal to true
func (i *ConfirmModal) Show() {
	i.isVisible = true
}

// Hide sets the visibility flag of the modal to false
func (i *ConfirmModal) Hide() {
	i.isVisible = false
}

// IsAccepted returns whether the user accepted the modal
func (i *ConfirmModal) IsAccepted() bool { return i.accepted }

// Clear clears all of the non-positional parameters of the modal
func (i *ConfirmModal) Clear() {
	i.title = ""
	i.text = ""
	i.accepted = false
	i.isDone = false
}

// IsBordered returns whether this modal is bordered or not
func (i *ConfirmModal) IsBordered() bool {
	return i.bordered
}

// SetBordered sets whether we render a border around the frame
func (i *ConfirmModal) SetBordered(b bool) {
	i.bordered = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (i *ConfirmModal) IsTabSkipped() bool {
	return i.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (i *ConfirmModal) SetTabSkip(b bool) {
	i.tabSkip = b
}

// HandleEvent handles the termbox event and returns whether it was consumed
func (i *ConfirmModal) HandleEvent(event termbox.Event) bool {
	if event.Ch == 'Y' || event.Ch == 'y' {
		i.accepted = true
		i.isDone = true
		return true
	} else if event.Ch == 'N' || event.Ch == 'n' {
		i.accepted = false
		i.isDone = true
		return true
	}
	return false
}

// Draw draws the modal
func (i *ConfirmModal) Draw() {
	// First blank out the area we'll be putting the modal
	FillWithChar(' ', i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)
	// Now draw the border
	DrawBorder(i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)

	nextY := i.y + 1
	// The title
	if i.title != "" {
		DrawStringAtPoint(i.title, i.x+1, nextY, i.fg, i.bg)
		nextY++
		FillWithChar('-', i.x+1, nextY, i.x+i.width-1, nextY, i.fg, i.bg)
		nextY++
	}
	if i.text != "" {
		DrawStringAtPoint(i.text, i.x+1, nextY, i.fg, i.bg)
	}
	nextY += 2
	if i.showHelp {
		helpString := " (Y/y) Confirm. (N/n) Reject. "
		helpX := (i.x + i.width) - len(helpString) - 1
		DrawStringAtPoint(helpString, helpX, nextY, i.fg, i.bg)
	}
}
