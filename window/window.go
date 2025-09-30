package window

import "image/color"

type Screen interface {
	// Create a NewTextWindow. If a graphics screen, x, y, w and h
	// are in pixels.
	NewTextWindow(px, py, pw, ph int) (TextWindow, error)
	// Returns size in the smallest possible unit (e.g. pixels)
	ScreenSize() (int, int)
}

type TextArea interface {
	// Write a charaacter to the display, wrap, newlines, and
	// scrolling should all be supported.
	Write(byte) error

	// Erase the display
	Erase()

	// Returns the dimensions of the screen as text cells
	TextWidth() int
	TextHeight() int
	TextSize() (int, int)

	// Get and set the character position in text cells
	CursorX() int
	CursorY() int
	CursorXY() (int, int)
	SetCursorX(int)
	SetCursorY(int)
	SetCursorXY(x, y int)

	TextColor(ColorChar)
}

// WindowBase contains common methods to all windows
type WindowBase interface {
	// Change the window size and location (in pixels)
	ResizeWindow(x, y, w, h int) error

	// Activate / remove display borders
	ShowBorder(screenw, screenh int) error

	// windowXY is in pixels
	WindowXY() (int, int)
	WindowSize() (int, int)
}

// TextWindow is output for a screen that displays monospaced text
type TextWindow interface {
	WindowBase
	TextArea

	// Refresh the display
	Refresh()

	// Scroll the display up or down
	Scroll(int)

	// Show/hide cursor, this may affect other windows and should be set back to
	// off if it's turned on.
	Cursor(bool)
}

// PixelWindow is used to display pixels
// Note that nothing is guaranteed to actually
type PixelWindow interface {
	WindowBase

	// Change color
	Color(color.RGBA)

	// Set a point
	SetPoint(x int, y int)

	// Drawing all relative to window x, y most with
	// lower overhead than SetPoint
	HLine(x, y, w int)
	VLine(x, y, h int)
	FilledRect(x, y, w, h int)
	Text(s string, x, y int)
}
