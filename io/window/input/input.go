// Package input contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package input

import (
	"errors"
	"mattwach/rpngo/io/key"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
	"strings"
)

var ErrExit = errors.New("exit")

// Input gets input from the keyboard/keypad
type Input interface {
	GetChar() (key.Key, error)
}

type InputWindow struct {
	input Input
	txtw  window.TextWindow
	gl    *getLine
}

func Init(input Input, txtw window.TextWindow) (*InputWindow, error) {
	iw := &InputWindow{
		input: input,
		txtw:  txtw,
		gl:    initGetLine(input, txtw),
	}
	if err := txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return nil, err
	}
	return iw, nil
}

func (iw *InputWindow) Update(rpn *rpn.RPN) error {
	line, err := iw.gl.get()
	if err != nil {
		window.PrintErr(iw.txtw, err)
		return nil
	}
	if line == "exit" {
		return ErrExit
	}
	action, err := parseLine(rpn, line)
	if err != nil {
		window.PrintErr(iw.txtw, err)
		return nil
	}
	if action {
		frame, err := rpn.Stack.Peek(0)
		if err != nil {
			window.PrintErr(iw.txtw, err)
		} else {
			window.Print(iw.txtw, frame.String())
			window.PutByte(iw.txtw, '\n')
		}
	}
	return nil
}

func (iw *InputWindow) Resize(x, y, w, h int) {
	iw.txtw.Resize(x, y, w, h)
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
