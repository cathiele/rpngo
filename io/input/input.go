// Package input contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package input

import (
	"errors"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
	"strings"
)

// ASCII for most keys, except the special ones below
type Key int

const (
	KEY_LEFT Key = iota + 256
	KEY_RIGHT
	KEY_UP
	KEY_DOWN
	KEY_BACKSPACE
	KEY_DEL
	KEY_INS
	KEY_END
	KEY_HOME
	KEY_EOF
)

// Input gets input from the keyboard/keypad
type Input interface {
	GetChar() (Key, error)
}

func Loop(rpn *rpn.RPN, input Input, root *window.WindowGroup, screen window.Screen) error {
	txtd := root.FindTextWindow("i")
	if txtd == nil {
		return errors.New("could not find window 'i'")
	}
	if err := txtd.Color(31, 31, 31, 0, 0, 0); err != nil {
		return err
	}

	gl := initGetLine(input, txtd)
	for {
		line, err := gl.get()
		if err != nil {
			window.PrintErr(txtd, err)
			continue
		}
		if line == "exit" {
			return nil
		}
		action, err := parseLine(rpn, line)
		if err != nil {
			window.PrintErr(txtd, err)
			continue
		}
		if action {
			frame, err := rpn.Stack.Peek()
			if err != nil {
				window.PrintErr(txtd, err)
			} else {
				window.Print(txtd, frame.String())
				window.PutByte(txtd, '\n')
			}
		}
		screen.Refresh()
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
