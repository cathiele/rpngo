package window

type Screen interface {
	NewTextWindow(x, y, w, h int) TextWindow
}

// TextWindow is output for a screen that displays monospaced text
type TextWindow interface {
	// Refresh the display
	Refresh()

	// Resize the window
	Resize(x, y, w, h int)

	// Clear the display
	Clear() error

	// Write a charaacter to the display, wrap, newlines, and
	// scrolling should all be supported.
	Write(byte) error

	// Returns the dimensions of the screen
	Width() int
	Height() int
	Size() (int, int)

	// Get and set the character position
	X() int
	Y() int
	XY() (int, int)
	SetX(int)
	SetY(int)
	SetXY(x, y int)

	// Change the foreground and background colors (approximately) to the given
	// r, g, b values. each value ranges from 0 to 32 (foreground, then background)
	// If the display does not support color, these commands do nothing.
	Color(fr, fg, fb, br, bg, bb int) error

	// Scroll the display up or down
	Scroll(int)
}
