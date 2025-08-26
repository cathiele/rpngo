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

func (sw *StackWindow) Resize(x, y, w, h int) {
	sw.txtw.Resize(x, y, w, h)
}

func (sw *StackWindow) Update(rpn *rpn.RPN) error {
	sw.txtw.Erase()
	w, h := sw.txtw.Size()
	framesBack := h
	if rpn.Stack.Size() < framesBack {
		framesBack = rpn.Stack.Size()
	}
	if framesBack <= 0 {
		// nothing to do
		return nil
	}
	for i := 0; i < framesBack; i++ {
		f, err := rpn.Stack.Peek(i)
		if err != nil {
			return err
		}
		sw.txtw.SetXY(0, h-i-1)
		s := fmt.Sprintf("%d: %v", i+1, f.String())
		if len(s) > w {
			s = s[:w]
		}
		window.Print(sw.txtw, s)
	}
	sw.txtw.Refresh()
	return nil
}
