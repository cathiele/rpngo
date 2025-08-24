// Package iface contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package io

import (
	"mattwach/rpngo/io/key"
	"mattwach/rpngo/rpn"
	"strings"
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
	gl := initGetLine(input, txtd)
	for {
		line, err := gl.get()
		if err != nil {
			printErr(txtd, err)
			continue
		}
		if line == "exit" {
			return nil
		}
		action, err := parseLine(rpn, line)
		if err != nil {
			printErr(txtd, err)
			continue
		}
		if action {
			frame, err := rpn.Stack.Peek()
			if err != nil {
				printErr(txtd, err)
			} else {
				print(txtd, frame.String())
				putByte(txtd, '\n')
			}
		}
	}
}

func parseLine(rpn *rpn.RPN, line string) (bool, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false, nil
	}
	fields := strings.Fields(line)
	for _, arg := range fields {
		if err := rpn.Exec(arg); err != nil {
			return false, err
		}
	}
	return true, nil
}
