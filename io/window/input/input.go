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
	r.Input = iw.Input
	return iw, nil
}

func (iw *InputWindow) Print(msg string) {
	window.Print(iw.txtw, msg)
	if strings.Contains(msg, "\n") {
		iw.txtw.Refresh()
	}
}

func (iw *InputWindow) Input(r *rpn.RPN) (string, error) {
	return iw.gl.get(r)
}

func (iw *InputWindow) Update(r *rpn.RPN) error {
	if iw.firstInput {
		if err := iw.firstRun(r); err != nil {
			return err
		}
	}
	if err := iw.txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return err
	}
	// clear any pending ctrl-c
	select {
	case <-r.Interrupt:
		break
	default:
		break
	}
	r.TextWidth = iw.txtw.TextWidth()
	window.Print(iw.gl.txtd, "> ")
	line, err := iw.gl.get(r)
	if err := iw.txtw.Color(0, 31, 31, 0, 0, 0); err != nil {
		return err
	}
	if err != nil {
		window.PrintErr(iw.txtw, err)
		return nil
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
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
			window.Print(iw.txtw, frame.String(true))
			window.PutByte(iw.txtw, '\n')
		} else if !errors.Is(err, rpn.ErrNotEnoughStackFrames) {
			window.PrintErr(iw.txtw, err)
		}
	}
	iw.txtw.Refresh()
	return nil
}

func (iw *InputWindow) firstRun(r *rpn.RPN) error {
	iw.firstInput = false
	r.RegisterConceptHelp(map[string]string{"exiting": "Enter exit or type Ctrl-D to exit the program"})
	if err := iw.txtw.Color(0, 31, 31, 0, 0, 0); err != nil {
		return err
	}
	window.Print(iw.txtw, "Enter ? to list help topics, topic? for more details\n")
	if err := iw.txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return err
	}
	return nil
}

func (iw *InputWindow) ResizeWindow(x, y, w, h int) error {
	return iw.txtw.ResizeWindow(x, y, w, h)
}

func (iw *InputWindow) ShowBorder(screenw, screenh int) error {
	return iw.txtw.ShowBorder(screenw, screenh)
}

func (iw *InputWindow) WindowXY() (int, int) {
	return iw.txtw.WindowXY()
}

func (iw *InputWindow) WindowSize() (int, int) {
	return iw.txtw.WindowSize()
}

func (iw *InputWindow) Type() string {
	return "input"
}

func (iw *InputWindow) SetProp(name string, val rpn.Frame) error {
	return rpn.ErrNotSupported
}

func (iw *InputWindow) GetProp(name string) (rpn.Frame, error) {
	return rpn.Frame{}, rpn.ErrNotSupported
}

func (iw *InputWindow) ListProps() []string {
	return nil
}

func parseLine(r *rpn.RPN, line string) (bool, error) {
	fields, err := parse.Fields(line)
	if err != nil {
		return false, err
	}
	if err := r.Exec(fields); err != nil {
		return false, err
	}
	return true, nil
}
