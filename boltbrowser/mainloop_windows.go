// +build windows
package main

// Windows doesn't support process backgrounding like *nix.
// So we have a separate loop for it.

import "github.com/nsf/termbox-go"

func mainLoop(memBolt *BoltDB, style Style) {
	screens := defaultScreensForData(memBolt)
	displayScreen := screens[BrowserScreenIndex]
	layoutAndDrawScreen(displayScreen, style)
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			newScreenIndex := displayScreen.handleKeyEvent(event)
			if newScreenIndex < len(screens) {
				displayScreen = screens[newScreenIndex]
				layoutAndDrawScreen(displayScreen, style)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize {
			layoutAndDrawScreen(displayScreen, style)
		}
	}
}
