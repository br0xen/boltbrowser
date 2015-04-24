package main

import "github.com/nsf/termbox-go"

type Widget interface {
	layoutUnderPressure(pressure int) (int, int)
	drawAtPoint(cursor Cursor, x int, y int, pressure int, style Style) (int, int)
}

type WidgetSlice []Widget

func (widgets WidgetSlice) sizeAtPressure(pressure int) (int, int) {
	total_widget_width := 0
	max_widget_height := 0
	for _, widget := range widgets {
		widget_width, widget_height := widget.layoutUnderPressure(pressure)
		total_widget_width += widget_width
		if widget_height > max_widget_height {
			max_widget_height = widget_height
		}
	}
	return total_widget_width, max_widget_height
}

func (widgets WidgetSlice) numberVisibleAtPressure(pressure int) int {
	count := 0
	for _, widget := range widgets {
		widget_width, _ := widget.layoutUnderPressure(pressure)
		if widget_width > 0 {
			count++
		}
	}
	return count
}

func (widgets WidgetSlice) layout() (int, int) {
	width, _ := termbox.Size()
	pressure := 0
	spacing := 4
	padding := 2
	for ; pressure < 10; pressure++ {
		spacing = 4
		total_widget_width, _ := widgets.sizeAtPressure(pressure)
		num_spaces := widgets.numberVisibleAtPressure(pressure) - 1
		for ; total_widget_width+num_spaces*spacing > (width-2*padding) && spacing > 2; spacing-- {
		}
		if total_widget_width+num_spaces*spacing <= (width - 2*padding) {
			break
		}
	}
	return pressure, spacing
}

func listOfWidgets() WidgetSlice {
	all_widgets := [...]Widget{
		NavigationWidget(0),
		//CursorWidget(0),
		//OffsetWidget(0),
	}
	return all_widgets[:]
}

func heightOfWidgets() int {
	widgets := listOfWidgets()
	pressure, _ := widgets.layout()
	_, max_widget_height := widgets.sizeAtPressure(pressure)
	return max_widget_height
}

func drawWidgets(cursor Cursor, style Style) (int, int) {
	widgets := listOfWidgets()

	width, height := termbox.Size()
	spacing := 4
	padding := 2
	pressure, spacing := widgets.layout()
	total_widget_width, max_widget_height := widgets.sizeAtPressure(pressure)
	num_spaces := widgets.numberVisibleAtPressure(pressure) - 1
	start_x, start_y := (width-2*padding-(total_widget_width+num_spaces*spacing))/2+padding, height-max_widget_height
	x, y := start_x, start_y
	for _, widget := range widgets {
		widget_width, _ := widget.drawAtPoint(cursor, x, y, pressure, style)
		x += widget_width
		if widget_width > 0 {
			x += spacing
		}
	}

	return total_widget_width + num_spaces*spacing, max_widget_height
}
