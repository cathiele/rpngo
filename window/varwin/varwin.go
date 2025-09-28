// Package varwin shows a variable window
package varwin

import (
	"fmt"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
	"strings"
)

type VariableWindow struct {
	txtb      window.TextBuffer
	txtw      window.TextWindow
	multiline bool
}

func Init(txtw window.TextWindow) (*VariableWindow, error) {
	w := &VariableWindow{txtw: txtw}
	w.txtb.TextColor(window.White)
	return w, nil
}

func (vw *VariableWindow) ResizeWindow(x, y, w, h int) error {
	return vw.txtw.ResizeWindow(x, y, w, h)
}

func (vw *VariableWindow) ShowBorder(screenw, screenh int) error {
	return vw.txtw.ShowBorder(screenw, screenh)
}

func (vw *VariableWindow) WindowXY() (int, int) {
	return vw.txtw.WindowXY()
}

func (vw *VariableWindow) WindowSize() (int, int) {
	return vw.txtw.WindowSize()
}

func (vw *VariableWindow) Type() string {
	return "var"
}

func (vw *VariableWindow) SetProp(name string, val rpn.Frame) error {
	if name == "multiline" {
		if val.Type != rpn.BOOL_FRAME {
			return rpn.ErrExpectedABoolean
		}
		vw.multiline = val.Int != 0
	} else {
		return rpn.ErrUnknownProperty
	}
	return nil
}

func (vw *VariableWindow) GetProp(name string) (rpn.Frame, error) {
	if name == "multiline" {
		return rpn.BoolFrame(vw.multiline), nil
	}
	return rpn.Frame{}, rpn.ErrUnknownProperty
}

func (vw *VariableWindow) ListProps() []string {
	return []string{"multiline"}
}

func (vw *VariableWindow) Update(rpn *rpn.RPN) error {
	w, h := vw.txtw.TextSize()
	vw.txtb.MaybeResize(int16(w), int16(h))
	vw.txtb.Erase()
	nv := rpn.AllVariableNamesAndValues()
	n := len(nv)
	if n > h {
		n = h - 1
	}
	vw.txtb.SetCursorXY(0, 0)
	for i := 0; i < n; i++ {
		name := nv[i].Name + ": "
		val := framesToString(nv[i].Values)
		if !vw.multiline {
			val = makeSingleLine(val, w-len(name))
		}
		vw.txtb.TextColor(window.White)
		window.Print(&vw.txtb, name)
		vw.txtb.TextColor(window.Cyan)
		window.Print(&vw.txtb, val)
		window.PutByte(&vw.txtb, '\n')
	}
	if len(nv) > h {
		window.Print(&vw.txtb, fmt.Sprintf("+ %d more\n", len(nv)-h+1))
	}
	vw.txtb.UpdateTextWindow(vw.txtw)
	vw.txtw.Refresh()
	return nil
}

func framesToString(frames []rpn.Frame) string {
	var parts []string
	for _, f := range frames {
		parts = append(parts, f.String(true))
	}
	return strings.Join(parts, " -> ")
}

func makeSingleLine(line string, width int) string {
	if width < 0 {
		return ""
	}
	if strings.Contains(line, "\n") {
		line = removeCRsAndComments(line)
	}
	if len(line) < width {
		return line
	}
	if width < 4 {
		return line[:width]
	}
	return line[:width-4] + "..."
}

func removeCRsAndComments(line string) string {
	if !strings.Contains(line, "#") {
		// no comments
		return strings.ReplaceAll(line, "\n", " ")
	}
	var parts []string
	for _, part := range strings.Split(line, "\n") {
		commentIdx := strings.Index(part, "#")
		if commentIdx >= 0 {
			part = part[:commentIdx]
		}
		if len(part) == 0 {
			continue
		}
		parts = append(parts, strings.Fields(part)...)
	}
	return strings.Join(parts, " ")
}
