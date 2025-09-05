// Package varwin shows a variable window
package varwin

import (
	"fmt"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
	"strings"
)

type VariableWindow struct {
	txtw      window.TextWindow
	multiline bool
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

func (vw *VariableWindow) WindowXY() (int, int) {
	return vw.txtw.WindowXY()
}

func (vw *VariableWindow) Size() (int, int) {
	return vw.txtw.Size()
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
		return fmt.Errorf("unknown property: %s", name)
	}
	return nil
}

func (vw *VariableWindow) GetProp(name string) (rpn.Frame, error) {
	if name == "multiline" {
		return rpn.BoolFrame(vw.multiline), nil
	}
	return rpn.Frame{}, fmt.Errorf("unknown property: %s", name)
}

func (vw *VariableWindow) ListProps() []string {
	return []string{"multiline"}
}

func (vw *VariableWindow) Update(rpn *rpn.RPN) error {
	vw.txtw.Erase()
	w, h := vw.txtw.Size()
	nv := rpn.AllVariableNamesAndValues()
	n := len(nv)
	allShown := true
	if n > h {
		n = h - 1
		allShown = false
	}
	vw.txtw.SetXY(0, 0)
	for i := 0; i < n; i++ {
		name := nv[i].Name + ": "
		val := framesToString(nv[i].Values)
		if !vw.multiline {
			val = makeSingleLine(val, w-len(name))
		}
		vw.txtw.Color(31, 31, 31, 0, 0, 0)
		window.Print(vw.txtw, name)
		vw.txtw.Color(0, 31, 31, 0, 0, 0)
		window.Print(vw.txtw, val)
		window.PutByte(vw.txtw, '\n')
	}
	if !allShown {
		window.Print(vw.txtw, fmt.Sprintf("+ %d more\n", len(nv)-h))
	}
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
		parts = append(parts, part)
	}
	return strings.Join(parts, "\n")
}
