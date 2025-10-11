package window

import "image/color"

var BorderColor = color.RGBA{R: 100, B: 100}

type Screen interface {
	// Create a new TextWindow. If a graphics screen, x, y, w and h
	// are in pixels.
	NewTextWindow() (TextWindow, error)
	// Create a new PixelWindow. may panic on non-graphics displays
	// (don't configure non-graphic dsiplays with pixel windows)
	NewPixelWindow() (PixelWindow, error)
	// Returns size in the smallest possible unit (e.g. pixels)
	ScreenSize() (int, int)
}

type TextArea interface {
	// Updates the character at the given x, y.  It's recommended to
	// have a buffering layer to avoid redrawing identical characters and
	// thus taking a performance hit on slower displays (which includes
	// anything SPI)
	DrawChar(x, y int, char ColorChar)

	// Erase the display
	Erase()

	// Returns the dimensions of the screen as text cells
	TextWidth() int
	TextHeight() int
	TextSize() (int, int)

	// Scroll the display up or down
	Scroll(int)
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

	// Show/hide cursor, this may affect other windows and should be set back to
	// off if it's turned on.
	Cursor(bool)
}

// PixelWindow is used to display pixels
// Note that nothing is guaranteed to actually
type PixelWindow interface {
	WindowBase

	// Refresh the display
	Refresh()

	// Size without the border
	PixelSize() (int, int)

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
