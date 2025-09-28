// Package stackwin shows a stack window
package stackwin

import (
	"fmt"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
)

type StackWindow struct {
	txtb window.TextBuffer
	txtw window.TextWindow
}

func Init(txtw window.TextWindow) (*StackWindow, error) {
	w := &StackWindow{
		txtw: txtw,
	}
	w.txtb.TextColor(window.White)
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
	w, h := sw.txtw.TextSize()
	sw.txtb.MaybeResize(int16(w), int16(h))
	sw.txtb.Erase()
	framesBack := h
	if rpn.Size() < framesBack {
		framesBack = rpn.Size()
	}
	for i := 0; i < framesBack; i++ {
		f, err := rpn.PeekFrame(i)
		if err != nil {
			return err
		}
		sw.txtb.SetCursorXY(0, h-i-1)
		s := fmt.Sprintf("%d: %v", i, f.String(true))
		if len(s) > w {
			s = s[:w]
		}
		window.Print(&sw.txtb, s)
	}
	sw.txtb.UpdateTextWindow(sw.txtw)
	sw.txtw.Refresh()
	return nil
}
