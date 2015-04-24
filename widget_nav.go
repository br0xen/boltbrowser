package main

type NavigationWidget int

func (widget NavigationWidget) layoutUnderPressure(pressure int) (int, int) {
	layouts := map[int]string{
		0: "Navigate: ↓j ↑k",
	}
	runeCount := 0
	for _, _ = range layouts[0] {
		runeCount++
	}
	return runeCount, 2
}

func (widget NavigationWidget) drawAtPoint(cursor Cursor, x int, y int, pressure int, style Style) (int, int) {
	fg := style.default_fg
	bg := style.default_bg
	x_pos := x
	x_pos += drawStringAtPoint("Navigate: ↓j ↑k", x_pos, y, fg, bg)
	return x_pos - x, 2
}
