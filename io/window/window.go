package window

// TextWindow is output for a screen that displays monospaced text
type TextWindow interface {
	// Clear the display
	Clear() error

	// Refresh the display
	Refresh()

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
	SetXY(int, int)

	// Change the foreground and background colors (approximately) to the given
	// r, g, b values. each value ranges from 0 to 32 (foreground, then background)
	// If the display does not support color, these commands do nothing.
	Color(int, int, int, int, int, int) error

	// Scroll the display up or down
	Scroll(int)
}
