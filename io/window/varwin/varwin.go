// Package varwin shows a variable window
package varwin

import (
	"fmt"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
	"sort"
)

type VariableWindow struct {
	txtw window.TextWindow
}

func Init(txtw window.TextWindow) (*VariableWindow, error) {
	w := &VariableWindow{txtw: txtw}
	if err := txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return nil, err
	}
	return w, nil
}

func (vw *VariableWindow) Resize(x, y, w, h int) error {
	return vw.txtw.Resize(x, y, w, h)
}

func (vw *VariableWindow) ShowBorder(screenw, screenh int) error {
	return vw.txtw.ShowBorder(screenw, screenh)
}

func (vw *VariableWindow) Update(rpn *rpn.RPN) error {
	vw.txtw.Erase()
	h := vw.txtw.Height()
	var names []string
	for n := range rpn.Variables {
		names = append(names, n)
	}
	sort.Strings(names)
	n := len(names)
	allShown := true
	if n > h {
		n = h - 1
		allShown = false
	}
	vw.txtw.SetXY(0, 0)
	for i := 0; i < n; i++ {
		window.Print(vw.txtw, names[i])
		window.Print(vw.txtw, ": ")
		f := rpn.Variables[names[i]]
		window.Print(vw.txtw, f.String())
		window.PutByte(vw.txtw, '\n')
	}
	if !allShown {
		window.Print(vw.txtw, fmt.Sprintf("+ %d more\n", len(names)-h))
	}
	vw.txtw.Refresh()
	return nil
}
