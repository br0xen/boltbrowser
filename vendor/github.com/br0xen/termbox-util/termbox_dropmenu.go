package termboxUtil

import "github.com/nsf/termbox-go"

// DropMenu is a title that, when active drops a menu down
type DropMenu struct {
	id                     string
	title                  string
	x, y, width, height    int
	bg, fg                 termbox.Attribute
	selectedBg, selectedFg termbox.Attribute
	menu                   *Menu
	menuSelected           bool
	showMenu               bool
	bordered               bool
	tabSkip                bool
}

// CreateDropMenu Creates a menu with the specified attributes
func CreateDropMenu(title string, options []string, x, y, width, height int, fg, bg, selectedFg, selectedBg termbox.Attribute) *DropMenu {
	i := DropMenu{
		title: title,
		x:     x, y: y, width: width, height: height,
		fg: fg, bg: bg,
		selectedFg: fg, selectedBg: bg,
	}
	i.menu = CreateMenu("", options, x, y+2, width, height, fg, bg)
	return &i
}

// GetID returns this control's ID
func (i *DropMenu) GetID() string { return i.id }

// SetID sets this control's ID
func (i *DropMenu) SetID(newID string) {
	i.id = newID
}

// GetTitle returns the current title of the menu
func (i *DropMenu) GetTitle() string { return i.title }

// SetTitle sets the current title of the menu to s
func (i *DropMenu) SetTitle(s string) {
	i.title = s
}

// GetMenu returns the menu for this dropmenu
func (i *DropMenu) GetMenu() *Menu {
	return i.menu
}

// GetX returns the current x coordinate of the menu
func (i *DropMenu) GetX() int { return i.x }

// SetX sets the current x coordinate of the menu to x
func (i *DropMenu) SetX(x int) {
	i.x = x
}

// GetY returns the current y coordinate of the menu
func (i *DropMenu) GetY() int { return i.y }

// SetY sets the current y coordinate of the menu to y
func (i *DropMenu) SetY(y int) {
	i.y = y
}

// GetWidth returns the current width of the menu
func (i *DropMenu) GetWidth() int { return i.width }

// SetWidth sets the current menu width to width
func (i *DropMenu) SetWidth(width int) {
	i.width = width
}

// GetHeight returns the current height of the menu
func (i *DropMenu) GetHeight() int { return i.height }

// SetHeight set the height of the menu to height
func (i *DropMenu) SetHeight(height int) {
	i.height = height
}

// GetFgColor returns the foreground color
func (i *DropMenu) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *DropMenu) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *DropMenu) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *DropMenu) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// IsBordered returns the bordered flag
func (i *DropMenu) IsBordered() bool { return i.bordered }

// SetBordered sets the bordered flag
func (i *DropMenu) SetBordered(b bool) {
	i.bordered = b
	i.menu.SetBordered(b)
}

// IsDone returns whether the user has answered the modal
func (i *DropMenu) IsDone() bool { return i.menu.isDone }

// SetDone sets whether the modal has completed it's purpose
func (i *DropMenu) SetDone(b bool) {
	i.menu.isDone = b
}

// IsTabSkipped returns whether this modal has it's tabskip flag set
func (i *DropMenu) IsTabSkipped() bool {
	return i.tabSkip
}

// SetTabSkip sets the tabskip flag for this control
func (i *DropMenu) SetTabSkip(b bool) {
	i.tabSkip = b
}

// ShowMenu tells the menu to draw the options
func (i *DropMenu) ShowMenu() {
	i.showMenu = true
	i.menuSelected = true
}

// HideMenu tells the menu to hide the options
func (i *DropMenu) HideMenu() {
	i.showMenu = false
	i.menuSelected = false
}

// HandleEvent handles the termbox event and returns whether it was consumed
func (i *DropMenu) HandleEvent(event termbox.Event) bool {
	moveUp := (event.Key == termbox.KeyArrowUp || (i.menu.vimMode && event.Ch == 'k'))
	moveDown := (event.Key == termbox.KeyArrowDown || (i.menu.vimMode && event.Ch == 'j'))
	if i.menuSelected {
		selIdx := i.menu.GetSelectedIndex()
		if (moveUp && selIdx == 0) || (moveDown && selIdx == (len(i.menu.options)-1)) {
			i.menuSelected = false
		} else {
			if i.menu.HandleEvent(event) {
				if i.menu.IsDone() {
					i.HideMenu()
				}
				return true
			}
		}
	} else {
		i.ShowMenu()
		return true
	}
	return false
}

// Draw draws the menu
func (i *DropMenu) Draw() {
	// The title
	ttlFg, ttlBg := i.fg, i.bg
	if !i.menuSelected {
		ttlFg, ttlBg = i.selectedFg, i.selectedBg
	}
	ttlTxt := i.title
	if i.showMenu {
		ttlTxt = ttlTxt + "-Showing Menu"
	}
	DrawStringAtPoint(AlignText(i.title, i.width, AlignLeft), i.x, i.y, ttlFg, ttlBg)
	if i.showMenu {
		i.menu.Draw()
	}
}
