// Package stackwin shows a stack window
package stackwin

import (
	"fmt"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
)

type StackWindow struct {
	txtw window.TextWindow
}

func Init(txtw window.TextWindow) (*StackWindow, error) {
	w := &StackWindow{
		txtw: txtw,
	}
	if err := txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return nil, err
	}
	return w, nil
}

func (sw *StackWindow) ResizeWindow(x, y, w, h int) error {
	return sw.txtw.ResizeWindow(x, y, w, h)
}

func (sw *StackWindow) ShowBorder(screenw, screenh int) error {
	return sw.txtw.ShowBorder(screenw, screenh)
}

func (sw *StackWindow) WindowXY() (int, int) {
	return sw.txtw.WindowXY()
}

func (sw *StackWindow) WindowSize() (int, int) {
	return sw.txtw.WindowSize()
}

func (sw *StackWindow) Type() string {
	return "stack"
}

func (sw *StackWindow) SetProp(name string, val rpn.Frame) error {
	return rpn.ErrNotSupported
}

func (sw *StackWindow) GetProp(name string) (rpn.Frame, error) {
	return rpn.Frame{}, rpn.ErrNotSupported
}

func (sw *StackWindow) ListProps() []string {
	return nil
}

func (sw *StackWindow) Update(rpn *rpn.RPN) error {
	sw.txtw.Erase()
	w, h := sw.txtw.TextSize()
	framesBack := h
	if rpn.Size() < framesBack {
		framesBack = rpn.Size()
	}
	for i := 0; i < framesBack; i++ {
		f, err := rpn.PeekFrame(i)
		if err != nil {
			return err
		}
		sw.txtw.SetCursorXY(0, h-i-1)
		s := fmt.Sprintf("%d: %v", i, f.String(true))
		if len(s) > w {
			s = s[:w]
		}
		window.Print(sw.txtw, s)
	}
	sw.txtw.Refresh()
	return nil
}
