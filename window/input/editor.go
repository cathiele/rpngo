package input

import (
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
)

// Provides a UI for editing a multiline string
type editor struct {
	buff []byte
	txtb window.TextBuffer
	// Problem statement:
	//
	// - We have a buffer of bytes with possible \n
	// - Lines that are too long will wrap around
	// - We have a character position in our buffer which corresponds
	//   to some cx, cy in the text buffer

	// buffer index of the upper left character
	ulIdx int

	// current character index
	cIdx int
}

const EditHelp = "Invokes an editor on the head value of the stack. " +
	"Press ESC to push edits back to the stack.  The pushed value " +
	"will always be a string (but can be evaluated with @ if needed)."

func (iw *InputWindow) Edit(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	ed := editor{buff: []byte(f.String(false)), ulIdx: 0}
	ed.txtb.Init(iw.txtb.Txtw, 0)
	for {
		ed.renderDisplay()
		c, err := iw.input.GetChar()
		if err != nil {
			return err
		}
		switch c {
		case 27: // ESC
			tw, th := iw.txtb.Txtw.TextSize()
			iw.txtb.RefreshArea(0, 0, tw, th)
			return r.PushFrame(rpn.StringFrame(string(ed.buff)))
		case key.KEY_UP:
			ed.keyUpPressed()
		case key.KEY_DOWN:
			ed.keyDownPressed()
		case key.KEY_LEFT:
			ed.keyLeftPressed()
		case key.KEY_RIGHT:
			ed.keyRightPressed()
		}
	}
}

func (ed *editor) renderDisplay() {
	ed.txtb.Cursor(false)
	x := 0
	y := 0
	tw, th := ed.txtb.Txtw.TextSize()
	for _, c := range ed.buff[ed.ulIdx:] {
		if x >= tw {
			x = 0
			y++
		}
		if y < th {
			if c == '\n' {
				ed.txtb.DrawChar(x, y, window.Cyan|window.ColorChar('.'))
				x = 0
				y++
			} else {
				ed.txtb.DrawChar(x, y, window.White|window.ColorChar(c))
				x++
			}
		}
		if y >= th {
			break
		}
	}
	// update changed characters
	ed.txtb.Update()
	ed.txtb.Cursor(true)
}

func (ed *editor) keyUpPressed() {
	x, y := ed.txtb.CursorXY()
	y--
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) keyDownPressed() {
	x, y := ed.txtb.CursorXY()
	y++
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) keyLeftPressed() {
	if ed.cIdx <= 0 {
		return
	}
	ed.cIdx--
	x, y := ed.txtb.CursorXY()
	if ed.buff[ed.cIdx] == '\n' {
		x = ed.findX()
		y--
	} else {
		x--
	}
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) findX() int {
	x := 0
	for i := ed.cIdx - 1; i >= 0; i-- {
		if ed.buff[i] == '\n' {
			break
		}
		x++
	}
	return x
}

func (ed *editor) keyRightPressed() {
	if ed.cIdx >= len(ed.buff) {
		return
	}
	x, y := ed.txtb.CursorXY()
	if ed.buff[ed.cIdx] == '\n' {
		x = 0
		y++
	} else {
		x++
	}
	ed.cIdx++
	ed.txtb.SetCursorXY(x, y)
}
