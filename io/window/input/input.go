// Package input contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package input

import (
	"errors"
	"mattwach/rpngo/io/key"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"strings"
)

var ErrExit = errors.New("exit")

// Input gets input from the keyboard/keypad
type Input interface {
	GetChar() (key.Key, error)
}

type InputWindow struct {
	input      Input
	txtw       window.TextWindow
	gl         *getLine
	firstInput bool
}

func Init(input Input, txtw window.TextWindow, r *rpn.RPN) (*InputWindow, error) {
	iw := &InputWindow{
		input:      input,
		txtw:       txtw,
		gl:         initGetLine(input, txtw),
		firstInput: true,
	}
	r.Print = iw.Print
	return iw, nil
}

func (iw *InputWindow) Print(msg string) {
	window.Print(iw.txtw, msg)
	if strings.Contains(msg, "\n") {
		iw.txtw.Refresh()
	}
}

func (iw *InputWindow) Update(r *rpn.RPN) error {
	if iw.firstInput {
		if err := iw.firstRun(r); err != nil {
			return err
		}
	}
	line, err := iw.gl.get()
	if err != nil {
		window.PrintErr(iw.txtw, err)
		return nil
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}
	if line[len(line)-1] == '?' {
		if err := r.PrintHelp(line[:len(line)-1], iw.txtw.Width()); err != nil {
			window.PrintErr(iw.txtw, err)
		}
		return nil
	}
	if line == "exit" {
		return ErrExit
	}
	action, err := parseLine(r, line)
	if err != nil {
		window.PrintErr(iw.txtw, err)
		return nil
	}
	if action {
		frame, err := r.PeekFrame(0)
		if err == nil {
			window.Print(iw.txtw, frame.String())
			window.PutByte(iw.txtw, '\n')
		} else if !errors.Is(err, rpn.ErrStackEmpty) {
			window.PrintErr(iw.txtw, err)
		}
	}
	iw.txtw.Refresh()
	return nil
}

func (iw *InputWindow) firstRun(r *rpn.RPN) error {
	iw.firstInput = false
	r.RegisterConceptHelp("exiting", "Enter exit or type Ctrl-D to exit the program")
	if err := iw.txtw.Color(0, 31, 31, 0, 0, 0); err != nil {
		return err
	}
	window.Print(iw.txtw, "Enter ? to list help topics, topic? for more details\n")
	if err := iw.txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return err
	}
	return nil
}

func (iw *InputWindow) Resize(x, y, w, h int) error {
	return iw.txtw.Resize(x, y, w, h)
}

func (iw *InputWindow) ShowBorder(screenw, screenh int) error {
	return iw.txtw.ShowBorder(screenw, screenh)
}

func (iw *InputWindow) WindowXY() (int, int) {
	return iw.txtw.WindowXY()
}

func (iw *InputWindow) Size() (int, int) {
	return iw.txtw.Size()
}

func (iw *InputWindow) Type() string {
	return "input"
}

func parseLine(r *rpn.RPN, line string) (bool, error) {
	fields, err := parse.Fields(line)
	if err != nil {
		return false, err
	}
	for _, arg := range fields {
		if err := r.Exec(arg); err != nil {
			return false, err
		}
	}
	return true, nil
}
