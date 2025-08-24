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

	// Write charaacters to the display, newlines and scrolling
	// have to be handled by the client
	Write([]byte) error

	// Returns the dimensions of the screen
	Width() uint
	Height() uint

	// Get and set the character position
	X() uint
	Y() uint
	SetX(uint)
	SetY(uint)

	// Scroll the display up or down
	Scroll(int)
}

func Loop(rpn *rpn.RPN, input Input, txtd TextDisplay) error {
	for {
		line, err := getLine(input, txtd)
		print(txtd, string(line))
		newLine(txtd)
		if err != nil {
			return err
		}
		if line == "exit" {
			return nil
		}
	}
}
