// Package iface contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package io

import (
	"mattwach/rpngo/io/key"
	"mattwach/rpngo/rpn"
)

// Input gets input from the keyboard/keypad
type Input interface {
	GetChar() (key.Key, error)
}

// TextDisplay is output for a screen that displays monospaced text
type TextDisplay interface {
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

	// Scroll the display up or down
	Scroll(int)
}

func Loop(rpn *rpn.RPN, input Input, txtd TextDisplay) error {
	for {
		line, err := getLine(input, txtd)
		print(txtd, string(line))
		putByte(txtd, '\n')
		if err != nil {
			return err
		}
		if line == "exit" {
			return nil
		}
	}
}
