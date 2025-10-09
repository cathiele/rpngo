// Package varwin shows a variable window
package varwin

import (
	"fmt"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
	"strings"
)

type VariableWindow struct {
	txtb           window.TextBuffer
	txtw           window.TextWindow
	showdot        bool
	multiline      bool
	namesAndValues []rpn.NameAndValues
}

func Init(txtw window.TextWindow) (*VariableWindow, error) {
	w := &VariableWindow{txtw: txtw} // object allocated on the heap: escapes at line 21
	w.txtb.TextColor(window.White)
	w.namesAndValues = make([]rpn.NameAndValues, 0, 16)
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
	switch name {
	case "showdot":
		if val.Type != rpn.BOOL_FRAME {
			return rpn.ErrExpectedABoolean
		}
		vw.showdot = val.Bool()
		return nil
	case "multiline":
		if val.Type != rpn.BOOL_FRAME {
			return rpn.ErrExpectedABoolean
		}
		vw.multiline = val.Bool()
		return nil
	default:
		return rpn.ErrUnknownProperty
	}
}

func (vw *VariableWindow) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "showdot":
		return rpn.BoolFrame(vw.showdot), nil
	case "multiline":
		return rpn.BoolFrame(vw.multiline), nil
	default:
		return rpn.Frame{}, rpn.ErrUnknownProperty
	}
}

func (vw *VariableWindow) ListProps() []string {
	return []string{"showdot", "multiline"} // object allocated on the heap: escapes at line 75
}

func (vw *VariableWindow) Update(r *rpn.RPN) error {
	w, h := vw.txtw.TextSize()
	vw.txtb.MaybeResize(int16(w), int16(h))
	vw.txtb.Erase()
	vw.namesAndValues = r.AppendAllVariableNamesAndValues(vw.namesAndValues[:0])
	vw.txtb.SetCursorXY(0, 0)
	hidden := 0
	row := 0
	for i := 0; i < len(vw.namesAndValues); i++ {
		if !vw.showdot && (len(vw.namesAndValues[i].Name) > 0) && (vw.namesAndValues[i].Name[0] == '.') {
			hidden++
			continue
		}
		if row < (h - 1) {
			name := vw.namesAndValues[i].Name + ": "
			val := framesToString(vw.namesAndValues[i].Values)
			if !vw.multiline {
				val = makeSingleLine(val, w-len(name))
			} else {
				row += countCRs(val)
			}
			vw.txtb.TextColor(window.White)
			window.Print(&vw.txtb, name)
			vw.txtb.TextColor(window.Cyan)
			window.Print(&vw.txtb, val)
			window.PutByte(&vw.txtb, '\n')
			row++
		}
	}
	vw.txtb.TextColor(window.Blue)
	vw.txtb.SetCursorXY(0, h-1)
	window.Print(&vw.txtb, fmt.Sprintf("num: %v hidden:%v", len(vw.namesAndValues), hidden))
	vw.txtb.UpdateTextWindow(vw.txtw)
	vw.txtw.Refresh()
	return nil
}

func countCRs(val string) int {
	n := 0
	for _, c := range val {
		if c == '\n' {
			n++
		}
	}
	return n
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
