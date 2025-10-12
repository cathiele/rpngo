// Package input contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package input

import (
	"errors"
	"mattwach/rpngo/key"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
	"strings"
)

var ErrExit = errors.New("exit")

// Input gets input from the keyboard/keypad
type Input interface {
	GetChar() (key.Key, error)
}

type InputWindow struct {
	input          Input
	txtb           window.TextBuffer
	gl             *getLine
	firstInput     bool
	scrollbackMode bool
}

func (iw *InputWindow) Init(input Input, txtw window.TextWindow, r *rpn.RPN, scrollbytes int) {
	iw.input = input
	// Make this >0 when we are ready to try scrolling.
	iw.txtb.Init(txtw, scrollbytes)
	iw.gl = initGetLine(input, &iw.txtb)
	iw.firstInput = true
	r.Print = iw.Print
	r.Input = iw.Input
}

func (iw *InputWindow) Print(msg string) {
	iw.txtb.Print(msg, true)
}

func (iw *InputWindow) Input(r *rpn.RPN) (string, error) {
	return iw.gl.get(r)
}

func (iw *InputWindow) Update(r *rpn.RPN) error {
	iw.txtb.CheckSize()
	if iw.firstInput {
		if err := iw.firstRun(r); err != nil {
			return err
		}
	}
	iw.txtb.TextColor(window.White)
	r.TextWidth = iw.txtb.Txtw.TextWidth()
	iw.txtb.Print("> ", true)
	line, err := iw.gl.get(r)
	iw.txtb.TextColor(window.Cyan)
	if err != nil {
		iw.txtb.PrintErr(err, true)
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
		iw.txtb.PrintErr(err, true)
		return nil
	}
	if action {
		frame, err := r.PeekFrame(0)
		if err == nil {
			iw.txtb.Print(frame.String(true), false)
			iw.txtb.Write('\n', false)
		} else if !errors.Is(err, rpn.ErrNotEnoughStackFrames) {
			iw.txtb.PrintErr(err, true)
		}
	}
	iw.txtb.Update()
	return nil
}

func (iw *InputWindow) firstRun(r *rpn.RPN) error {
	iw.firstInput = false
	r.RegisterConceptHelp(map[string]string{"exiting": "Enter exit or type Ctrl-D to exit the program"})
	iw.txtb.TextColor(window.Cyan)
	iw.txtb.Print("Enter ? to list help topics, topic? for more details\n", true)
	iw.txtb.TextColor(window.White)
	return nil
}

func (iw *InputWindow) ResizeWindow(x, y, w, h int) error {
	err := iw.txtb.Txtw.ResizeWindow(x, y, w, h)
	if err != nil {
		return err
	}
	iw.txtb.CheckSize()
	return nil
}

func (iw *InputWindow) ShowBorder(screenw, screenh int) error {
	return iw.txtb.Txtw.ShowBorder(screenw, screenh)
}

func (iw *InputWindow) WindowXY() (int, int) {
	return iw.txtb.Txtw.WindowXY()
}

func (iw *InputWindow) WindowSize() (int, int) {
	return iw.txtb.Txtw.WindowSize()
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

var fields = make([]string, 128)

func parseLine(r *rpn.RPN, line string) (bool, error) {
	fields = fields[:0]
	var err error
	fields, err = parse.Fields(line, fields)
	if err != nil {
		return false, err
	}
	if err := r.Exec(fields); err != nil {
		return false, err
	}
	return true, nil
}
