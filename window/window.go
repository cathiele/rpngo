package window

type Screen interface {
	// Create a NewTextWindow. If a graphics screen, x, y, w and h
	// are in pixels.
	NewTextWindow(px, py, pw, ph int) (TextWindow, error)
	// Returns size in the smallest possible unit (e.g. pixels)
	ScreenSize() (int, int)
}

// TextWindow is output for a screen that displays monospaced text
type TextWindow interface {
	// Refresh the display
	Refresh()

	// Erase the display
	Erase()

	// Activate / remove display borders
	ShowBorder(screenw, screenh int) error

	// Write a charaacter to the display, wrap, newlines, and
	// scrolling should all be supported.
	Write(byte) error

	// Returns the dimensions of the screen as text cells
	TextWidth() int
	TextHeight() int
	TextSize() (int, int)

	// windowXY is in pixels
	WindowXY() (int, int)
	WindowSize() (int, int)
	ResizeWindow(px, py, pw, ph int) error

	// Get and set the character position in text cells
	CursorX() int
	CursorY() int
	CursorXY() (int, int)
	SetCursorX(int)
	SetCursorY(int)
	SetCursorXY(x, y int)

	// Change the foreground and background colors (approximately) to the given
	// r, g, b values. each value ranges from 0 to 32 (foreground, then background)
	// If the display does not support color, these commands do nothing.
	Color(fr, fg, fb, br, bg, bb int) error

	// Scroll the display up or down
	Scroll(int)

	// Show/hide cursor, this may affect other windows and should be set back to
	// off if it's turned on.
	Cursor(bool)
}
