// Package iface contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package io

// Input gets input from the keyboard/keypad
type Input interface {
	ReadByte() (byte, error)
}

// TextDisplay is output for a screen that displays monospaced text
type TextDisplay interface {
	// Clear the display
	Clear() error

	// Write charaacters to the display, \n causes a newline, characters
	// that exceed therigth margin will be wrapped
	Write([]byte) error

	// Returns the dimensions of the screen
	Width() uint
	Height() uint

	// Get and set the character position
	X() uint
	Y() uint
	SetX(uint)
	SetY(uint)
}
