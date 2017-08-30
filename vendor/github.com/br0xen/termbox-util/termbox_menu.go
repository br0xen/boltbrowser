package termboxUtil

import "github.com/nsf/termbox-go"

// Menu is a menu with a list of options
type Menu struct {
	id      string
	title   string
	options []MenuOption
	// If height is -1, then it is adaptive to the menu
	x, y, width, height    int
	showHelp               bool
	cursor                 int
	bg, fg                 termbox.Attribute
	selectedBg, selectedFg termbox.Attribute
	disabledBg, disabledFg termbox.Attribute
	isDone                 bool
	bordered               bool
	vimMode                bool
	tabSkip                bool
}

// CreateMenu Creates a menu with the specified attributes
func CreateMenu(title string, options []string, x, y, width, height int, fg, bg termbox.Attribute) *Menu {
	i := Menu{
		title: title,
		x:     x, y: y, width: width, height: height,
		fg: fg, bg: bg, selectedFg: bg, selectedBg: fg,
		disabledFg: bg, disabledBg: bg,
	}
	for _, line := range options {
		i.options = append(i.options, MenuOption{text: line})
	}
	if len(i.options) > 0 {
		i.SetSelectedOption(&i.options[0])
	}
	return &i
}

// GetID returns this control's ID
func (i *Menu) GetID() string { return i.id }

// SetID sets this control's ID
func (i *Menu) SetID(newID string) {
	i.id = newID
}

// GetTitle returns the current title of the menu
func (i *Menu) GetTitle() string { return i.title }

// SetTitle sets the current title of the menu to s
func (i *Menu) SetTitle(s string) {
	i.title = s
}

// GetOptions returns the current options of the menu
func (i *Menu) GetOptions() []MenuOption {
	return i.options
}

// SetOptions set the menu's options to opts
func (i *Menu) SetOptions(opts []MenuOption) {
	i.options = opts
}

// SetOptionsFromStrings sets the options of this menu from a slice of strings
func (i *Menu) SetOptionsFromStrings(opts []string) {
	var newOpts []MenuOption
	for _, v := range opts {
		newOpts = append(newOpts, *CreateOptionFromText(v))
	}
	i.SetOptions(newOpts)
	i.SetSelectedOption(i.GetOptionFromIndex(0))
}

// GetX returns the current x coordinate of the menu
func (i *Menu) GetX() int { return i.x }

// SetX sets the current x coordinate of the menu to x
func (i *Menu) SetX(x int) {
	i.x = x
}

// GetY returns the current y coordinate of the menu
func (i *Menu) GetY() int { return i.y }

// SetY sets the current y coordinate of the menu to y
func (i *Menu) SetY(y int) {
	i.y = y
}

// GetWidth returns the current width of the menu
func (i *Menu) GetWidth() int { return i.width }

// SetWidth sets the current menu width to width
func (i *Menu) SetWidth(width int) {
	i.width = width
}

// GetHeight returns the current height of the menu
func (i *Menu) GetHeight() int { return i.height }

// SetHeight set the height of the menu to height
func (i *Menu) SetHeight(height int) {
	i.height = height
}

// GetSelectedOption returns the current selected option
func (i *Menu) GetSelectedOption() *MenuOption {
	idx := i.GetSelectedIndex()
	if idx != -1 {
		return &i.options[idx]
	}
	return nil
}

// GetOptionFromIndex Returns the
func (i *Menu) GetOptionFromIndex(idx int) *MenuOption {
	if idx >= 0 && idx < len(i.options) {
		return &i.options[idx]
	}
	return nil
}

// GetOptionFromText Returns the first option with the text v
func (i *Menu) GetOptionFromText(v string) *MenuOption {
	for idx := range i.options {
		testOption := &i.options[idx]
		if testOption.GetText() == v {
			return testOption
		}
	}
	return nil
}

// GetSelectedIndex returns the index of the selected option
// Returns -1 if nothing is selected
func (i *Menu) GetSelectedIndex() int {
	for idx := range i.options {
		if i.options[idx].IsSelected() {
			return idx
		}
	}
	return -1
}

// SetSelectedIndex sets the selection to setIdx
func (i *Menu) SetSelectedIndex(idx int) {
	if len(i.options) > 0 {
		if idx < 0 {
			idx = 0
		} else if idx >= len(i.options) {
			idx = len(i.options) - 1
		}
		i.SetSelectedOption(&i.options[idx])
	}
}

// SetSelectedOption sets the current selected option to v (if it's valid)
func (i *Menu) SetSelectedOption(v *MenuOption) {
	for idx := range i.options {
		if &i.options[idx] == v {
			i.options[idx].Select()
		} else {
			i.options[idx].Unselect()
		}
	}
}

// SelectPrevOption Decrements the selected option (if it can)
func (i *Menu) SelectPrevOption() {
	idx := i.GetSelectedIndex()
	for idx >= 0 {
		idx--
		testOption := i.GetOptionFromIndex(idx)
		if testOption != nil && !testOption.IsDisabled() {
			i.SetSelectedOption(testOption)
			return
		}
	}
}

// SelectNextOption Increments the selected option (if it can)
func (i *Menu) SelectNextOption() {
	idx := i.GetSelectedIndex()
	for idx < len(i.options) {
		idx++
		testOption := i.GetOptionFromIndex(idx)
		if testOption != nil && !testOption.IsDisabled() {
			i.SetSelectedOption(testOption)
			return
		}
	}
}

// SelectPageUpOption Goes up 'menu height' options
func (i *Menu) SelectPageUpOption() {
	idx := i.GetSelectedIndex()
	idx -= i.height
	if idx < 0 {
		idx = 0
	}
	i.SetSelectedIndex(idx)
	return
}

// SelectPageDownOption Goes down 'menu height' options
func (i *Menu) SelectPageDownOption() {
	idx := i.GetSelectedIndex()
	idx += i.height
	if idx >= len(i.options) {
		idx = len(i.options) - 1
	}
	i.SetSelectedIndex(idx)
	return
}

// SelectFirstOption Goes to the top
func (i *Menu) SelectFirstOption() {
	i.SetSelectedIndex(0)
	return
}

// SelectLastOption Goes to the bottom
func (i *Menu) SelectLastOption() {
	i.SetSelectedIndex(len(i.options) - 1)
	return
}

// SetOptionDisabled Disables the specified option
func (i *Menu) SetOptionDisabled(idx int) {
	if len(i.options) > idx {
		i.GetOptionFromIndex(idx).Disable()
	}
}

// SetOptionEnabled Enables the specified option
func (i *Menu) SetOptionEnabled(idx int) {
	if len(i.options) > idx {
		i.GetOptionFromIndex(idx).Enable()
	}
}

// HelpIsShown returns true or false if the help is displayed
func (i *Menu) HelpIsShown() bool { return i.showHelp }

// ShowHelp sets whether or not to display the help text
func (i *Menu) ShowHelp(b bool) {
	i.showHelp = b
}

// GetFgColor returns the foreground color
func (i *Menu) GetFgColor() termbox.Attribute { return i.fg }

// SetFgColor sets the foreground color
func (i *Menu) SetFgColor(fg termbox.Attribute) {
	i.fg = fg
}

// GetBgColor returns the background color
func (i *Menu) GetBgColor() termbox.Attribute { return i.bg }

// SetBgColor sets the current background color
func (i *Menu) SetBgColor(bg termbox.Attribute) {
	i.bg = bg
}

// IsDone returns whether the user has answered the modal
func (i *Menu) IsDone() bool { return i.isDone }

// SetDone sets whether the modal has completed it's purpose
func (i *Menu) SetDone(b bool) {
	i.isDone = b
}

// IsBordered returns true or false if this menu has a border
func (i *Menu) IsBordered() bool { return i.bordered }

// SetBordered sets whether we render a border around the menu
func (i *Menu) SetBordered(b bool) {
	i.bordered = b
}

// EnableVimMode Enables h,j,k,l navigation
func (i *Menu) EnableVimMode() {
	i.vimMode = true
}

// DisableVimMode Disables h,j,k,l navigation
func (i *Menu) DisableVimMode() {
	i.vimMode = false
}

// HandleEvent handles the termbox event and returns whether it was consumed
func (i *Menu) HandleEvent(event termbox.Event) bool {
	if event.Key == termbox.KeyEnter || event.Key == termbox.KeySpace {
		i.isDone = true
		return true
	}
	currentIdx := i.GetSelectedIndex()
	switch event.Key {
	case termbox.KeyArrowUp:
		i.SelectPrevOption()
	case termbox.KeyArrowDown:
		i.SelectNextOption()
	case termbox.KeyArrowLeft:
		i.SelectPageUpOption()
	case termbox.KeyArrowRight:
		i.SelectPageDownOption()
	}
	if i.vimMode {
		switch event.Ch {
		case 'j':
			i.SelectNextOption()
		case 'k':
			i.SelectPrevOption()
		}
		if event.Key == termbox.KeyCtrlF {
			i.SelectPageDownOption()
		} else if event.Key == termbox.KeyCtrlB {
			i.SelectPageUpOption()
		}
	}
	if i.GetSelectedIndex() != currentIdx {
		return true
	}
	return false
}

// Draw draws the modal
func (i *Menu) Draw() {
	// First blank out the area we'll be putting the menu
	FillWithChar(' ', i.x, i.y, i.x+i.width, i.y+i.height, i.fg, i.bg)
	// Now draw the border
	optionStartX := i.x
	optionStartY := i.y
	optionWidth := i.width
	optionHeight := i.height
	if optionHeight == -1 {
		optionHeight = len(i.options)
	}
	if i.bordered {
		if i.height == -1 {
			DrawBorder(i.x, i.y, i.x+i.width, i.y+optionHeight+1, i.fg, i.bg)
		} else {
			DrawBorder(i.x, i.y, i.x+i.width, i.y+optionHeight, i.fg, i.bg)
		}
		optionStartX = i.x + 1
		optionStartY = i.y + 1
		optionWidth = i.width - 1
		optionHeight -= 2
	}

	// The title
	if i.title != "" {
		DrawStringAtPoint(AlignText(i.title, optionWidth, AlignCenter), optionStartX, optionStartY, i.fg, i.bg)
		optionStartY++
		if i.bordered {
			FillWithChar('-', optionStartX, optionStartY, optionWidth, optionStartY, i.fg, i.bg)
			optionStartY++
			optionHeight--
		}
		optionHeight--
	}

	if len(i.options) > 0 {
		// If the currently selected option is disabled, move to the next
		if i.GetSelectedOption().IsDisabled() {
			i.SelectNextOption()
		}

		// Print the options
		bldHeight := (optionHeight / 2)
		startIdx := i.GetSelectedIndex()
		endIdx := i.GetSelectedIndex()
		for bldHeight > 0 && startIdx >= 1 {
			startIdx--
			bldHeight--
		}
		bldHeight += (optionHeight / 2)
		for bldHeight > 0 && endIdx < len(i.options) {
			endIdx++
			bldHeight--
		}

		for idx := startIdx; idx < endIdx; idx++ { //i.options {
			if i.GetSelectedIndex()-idx >= optionHeight-1 {
				// Skip this one
				continue
			}
			currOpt := &i.options[idx]
			outTxt := currOpt.GetText()
			if len(outTxt) >= i.width {
				outTxt = outTxt[:i.width]
			}
			if currOpt.IsDisabled() {
				DrawStringAtPoint(outTxt, optionStartX, optionStartY, i.disabledFg, i.disabledBg)
			} else if i.GetSelectedOption() == currOpt {
				DrawStringAtPoint(outTxt, optionStartX, optionStartY, i.selectedFg, i.selectedBg)
			} else {
				DrawStringAtPoint(outTxt, optionStartX, optionStartY, i.fg, i.bg)
			}
			optionStartY++
			if optionStartY > i.y+optionHeight-1 {
				break
			}
		}
	}
}

/* MenuOption Struct & methods */

// MenuOption An option in the menu
type MenuOption struct {
	id       string
	text     string
	selected bool
	disabled bool
	helpText string
	subMenu  []MenuOption
}

// CreateOptionFromText just returns a MenuOption object
// That only has it's text value set.
func CreateOptionFromText(s string) *MenuOption {
	return &MenuOption{text: s}
}

// SetText Sets the text for this option
func (i *MenuOption) SetText(s string) {
	i.text = s
}

// GetText Returns the text for this option
func (i *MenuOption) GetText() string { return i.text }

// Disable Sets this option to disabled
func (i *MenuOption) Disable() {
	i.disabled = true
}

// Enable Sets this option to enabled
func (i *MenuOption) Enable() {
	i.disabled = false
}

// IsDisabled returns whether this option is enabled
func (i *MenuOption) IsDisabled() bool {
	return i.disabled
}

// IsSelected Returns whether this option is selected
func (i *MenuOption) IsSelected() bool {
	return i.selected
}

// Select Sets this option to selected
func (i *MenuOption) Select() {
	i.selected = true
}

// Unselect Sets this option to not selected
func (i *MenuOption) Unselect() {
	i.selected = false
}

// SetHelpText Sets this option's help text to s
func (i *MenuOption) SetHelpText(s string) {
	i.helpText = s
}

// GetHelpText Returns the help text for this option
func (i *MenuOption) GetHelpText() string { return i.helpText }

// AddToSubMenu adds a slice of MenuOptions to this option
func (i *MenuOption) AddToSubMenu(sub *MenuOption) {
	i.subMenu = append(i.subMenu, *sub)
}
