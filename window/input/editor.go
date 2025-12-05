package input

import (
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
)

// Provides a UI for editing a multiline string
type editor struct {
	buff      []byte
	clipboard []byte
	txtb      window.TextBuffer
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
	// select Index active if > 0.  This might be to the cursor or from the cursor
	selIdx int

	// if message is set, it will be shown until the user presses any key
	message string

	replaceMode bool
}

type HighlightState uint8

const (
	HIGHLIGHT_NORMAL HighlightState = iota
	HIGHLIGHT_VARIABLE
	HGHLIGHT_SINGLE_QUOTE
	HGHLIGHT_DOUBLE_QUOTE
	HIGHLIGHT_MACRO
	HIGHLIGHT_COMMENT
)

const editHelp = "Invokes an editor on the head value of the stack. " +
	"Press ESC to push edits back to the stack.  The pushed value " +
	"will always be a string (but can be evaluated with @ if needed)."

func (iw *InputWindow) edit(r *rpn.RPN) error {
	var f rpn.Frame
	var err error
	var ed editor
	if len(r.Frames) != 0 {
		f, err = r.PopFrame()
		if err != nil {
			return err
		}
		ed = editor{
			buff:    []byte(f.String(false)),
			ulIdx:   0,
			message: "Press Alt-H For Help",
		}
	}
	ed.selIdx = -1
	ed.txtb.Init(iw.txtb.Txtw, 0)

	quit := func() error {
		tw, th := iw.txtb.Txtw.TextSize()
		iw.txtb.RefreshArea(0, 0, tw, th)
		return r.PushFrame(f)
	}

	for {
		//ed.debugDump()
		if r.Interrupt() {
			return quit()
		}
		ed.renderDisplay()
		c, err := iw.input.GetChar()
		if err != nil {
			return err
		}
		clearSel := true
		selAnchorIdx := ed.cIdx
		switch c {
		case 27, key.KEY_QUIT: // ESC
			return quit()
		case key.KEY_SAVE:
			f = rpn.StringFrame(string(ed.buff), f.Type())
		case key.KEY_UP:
			ed.keyUpPressed()
		case key.KEY_DOWN:
			ed.keyDownPressed()
		case key.KEY_LEFT:
			ed.keyLeftPressed()
		case key.KEY_RIGHT:
			ed.keyRightPressed()
		case key.KEY_DEL:
			ed.delPressed()
		case key.KEY_BACKSPACE:
			ed.backspacePressed()
		case key.KEY_PAGEDOWN:
			ed.pageDownPressed()
		case key.KEY_PAGEUP:
			ed.pageUpPressed()
		case key.KEY_HOME:
			ed.homePressed()
		case key.KEY_END:
			ed.endPressed()
		case key.KEY_INS:
			ed.replaceMode = !ed.replaceMode
		case key.KEY_COPY:
			ed.copySelection()
		case key.KEY_PASTE:
			ed.paste()
		case key.KEY_SUP:
			clearSel = false
			ed.keyUpPressed()
		case key.KEY_SDOWN:
			clearSel = false
			ed.keyDownPressed()
		case key.KEY_SLEFT:
			clearSel = false
			ed.keyLeftPressed()
		case key.KEY_SRIGHT:
			clearSel = false
			ed.keyRightPressed()
		case key.KEY_SHOME:
			clearSel = false
			ed.homePressed()
		case key.KEY_SEND:
			clearSel = false
			ed.endPressed()
		case '\n':
			ed.insertOrReplaceChar(byte(c))
		default:
			if (c >= ' ') && (c <= 127) {
				ed.insertOrReplaceChar(byte(c))
			}
		}
		if c != 0 {
			if clearSel {
				ed.selIdx = -1
			} else if ed.selIdx < 0 {
				ed.selIdx = selAnchorIdx
			}
		}
	}
}

/*
func (ed *editor) debugDump() {
	x, y := ed.txtb.CursorXY()
	if ed.cIdx == len(ed.buff) {
		log.Printf("x=%v y=%v cidx=%v <end>", x, y, ed.cIdx)
	} else {
		log.Printf("x=%v y=%v cidx=%v c=%c", x, y, ed.cIdx, rune(ed.buff[ed.cIdx]))
	}
}
*/

func (ed *editor) renderDisplay() {
	var hs HighlightState = HIGHLIGHT_NORMAL
	ed.txtb.Cursor(false)
	x := 0
	y := 0
	sbegIdx := 0
	sendIdx := 0
	if ed.selIdx >= 0 {
		sbegIdx = ed.cIdx
		sendIdx = ed.selIdx
		if sbegIdx > sendIdx {
			sbegIdx, sendIdx = sendIdx, sbegIdx
		}
	}
	tw, th := ed.txtb.Txtw.TextSize()
	var col window.ColorChar
	var skip bool
	for i, c := range ed.buff[ed.ulIdx:] {
		if !skip {
			hs, col = checkHighlightState(hs, c)
		}
		skip = !skip && (c == '\\')
		if x >= tw {
			x = 0
			y++
		}
		if y < th {
			idx := ed.ulIdx + i
			if (idx >= sbegIdx) && (idx < sendIdx) {
				col |= 0x0100 // blue background
			}
			if c == '\n' {
				ed.txtb.DrawChar(x, y, window.Cyan|window.ColorChar('.'))
				ed.clearScreenToEndOfLine(x+1, y)
				x = 0
				y++
			} else {
				ed.txtb.DrawChar(x, y, col|window.ColorChar(c))
				x++
			}
		}
		if y >= th {
			break
		}
	}
	ed.clearScreenToBottomRightCorner(x, y)
	if ed.message != "" {
		ed.showMessage()
	}
	// update changed characters
	ed.txtb.Update()
	ed.txtb.Cursor(true)
}

func (ed *editor) showMessage() {
	tw, th := ed.txtb.Txtw.TextSize()
	x := tw - 1
	y := th - 1
	for i := len(ed.message) - 1; i >= 0; i-- {
		c := ed.message[i]
		if c == '\n' {
			x = 0
		} else {
			ed.txtb.DrawChar(x, y, window.Yellow|window.ColorChar(c))
		}
		x--
		if x < 0 {
			x = tw - 1
			y--
			if y < 0 {
				break
			}
		}
	}
	ed.message = ""
}

func checkHighlightState(hs HighlightState, c byte) (HighlightState, window.ColorChar) {
	var col window.ColorChar
	switch hs {
	case HIGHLIGHT_NORMAL:
		switch c {
		case '\'':
			hs = HGHLIGHT_SINGLE_QUOTE
			col = window.Red
		case '"':
			hs = HGHLIGHT_DOUBLE_QUOTE
			col = window.Red
		case '$':
			hs = HIGHLIGHT_VARIABLE
			col = window.Green
		case '@':
			hs = HIGHLIGHT_MACRO
			col = window.Yellow
		case '#':
			hs = HIGHLIGHT_COMMENT
			col = window.Blue
		default:
			col = window.White
		}
	case HGHLIGHT_SINGLE_QUOTE:
		if c == '\'' {
			hs = HIGHLIGHT_NORMAL
		}
		col = window.Red
	case HGHLIGHT_DOUBLE_QUOTE:
		if c == '"' {
			hs = HIGHLIGHT_NORMAL
		}
		col = window.Red
	case HIGHLIGHT_VARIABLE:
		if isWhiteSpace(c) {
			hs = HIGHLIGHT_NORMAL
		}
		col = window.Green
	case HIGHLIGHT_MACRO:
		if isWhiteSpace(c) {
			hs = HIGHLIGHT_NORMAL
		}
		col = window.Yellow
	case HIGHLIGHT_COMMENT:
		if c == '\n' {
			hs = HIGHLIGHT_NORMAL
		}
		col = window.Blue
	}
	return hs, col
}

func isWhiteSpace(c byte) bool {
	return (c == ' ') || (c == '\t') || (c == '\n')
}

func (ed *editor) clearScreenToEndOfLine(x, y int) {
	w := ed.txtb.Txtw.TextWidth()
	for x < w {
		ed.txtb.DrawChar(x, y, ' ')
		x++
	}
}

func (ed *editor) clearScreenToBottomRightCorner(x, y int) {
	w, h := ed.txtb.Txtw.TextSize()
	for y < h {
		for x < w {
			ed.txtb.DrawChar(x, y, ' ')
			x++
		}
		x = 0
		y++
	}
}

func (ed *editor) keyUpPressed() {
	x, y := ed.txtb.CursorXY()
	w := ed.txtb.Txtw.TextWidth()
	// we want to try and end up at the same x on the previous
	// line but this may not be possible if the line is short or
	// we hit the start of the buffer
	wantx := x
	for ed.cIdx > 0 {
		x--
		ed.cIdx--
		if ed.buff[ed.cIdx] == '\n' {
			x = ed.findX()
			y--
			break
		} else if x < 0 {
			x = w - 1
			y--
			break
		}
	}
	if x > wantx {
		ed.cIdx -= (x - wantx)
		x = wantx
	}
	y = ed.checkScroll(y)
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) keyDownPressed() {
	x, y := ed.txtb.CursorXY()
	w := ed.txtb.Txtw.TextWidth()
	// we want to try and end up at the same x on the next
	// line but this may not be possible if the line is short or
	// we hit the end of the buffer
	wantx := x
	for ed.cIdx < len(ed.buff) {
		x++
		ed.cIdx++
		if (x >= w) || (ed.buff[ed.cIdx-1] == '\n') {
			x = 0
			y++
			break
		}
	}
	for ed.cIdx < len(ed.buff) {
		if x == wantx {
			break
		}
		if ed.buff[ed.cIdx] == '\n' {
			break
		}
		x++
		ed.cIdx++
	}
	y = ed.checkScroll(y)
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) pageDownPressed() {
	lines := ed.txtb.Txtw.TextHeight() / 2
	for i := 0; (ed.cIdx < len(ed.buff)) && i < lines; i++ {
		ed.keyDownPressed()
	}
}

func (ed *editor) pageUpPressed() {
	lines := ed.txtb.Txtw.TextHeight() / 2
	for i := 0; (ed.cIdx > 0) && i < lines; i++ {
		ed.keyUpPressed()
	}
}

func (ed *editor) homePressed() {
	for ed.cIdx > 0 {
		if ed.buff[ed.cIdx-1] == '\n' {
			break
		}
		ed.keyLeftPressed()
	}
}

func (ed *editor) endPressed() {
	for ed.cIdx < len(ed.buff) {
		if ed.buff[ed.cIdx] == '\n' {
			break
		}
		ed.keyRightPressed()
	}
}

func (ed *editor) keyLeftPressed() {
	if ed.cIdx <= 0 {
		return
	}
	x, y := ed.txtb.CursorXY()
	ed.cIdx--
	x--
	if ed.buff[ed.cIdx] == '\n' {
		x = ed.findX()
		y--
	} else if x < 0 {
		x = ed.txtb.Txtw.TextWidth() - 1
		y--
	}
	y = ed.checkScroll(y)
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) findX() int {
	x := 0
	w := ed.txtb.Txtw.TextWidth()
	for i := ed.cIdx - 1; i >= 0; i-- {
		if ed.buff[i] == '\n' {
			break
		}
		x++
		if x >= w {
			x = 0
		}
	}
	return x
}

func (ed *editor) keyRightPressed() {
	if ed.cIdx >= len(ed.buff) {
		return
	}
	x, y := ed.txtb.CursorXY()
	if (ed.buff[ed.cIdx] == '\n') || (x >= (ed.txtb.Txtw.TextWidth() - 1)) {
		x = 0
		y++
	} else {
		x++
	}
	ed.cIdx++
	y = ed.checkScroll(y)
	ed.txtb.SetCursorXY(x, y)
}

func (ed *editor) insertOrReplaceChar(c byte) {
	if ed.selIdx >= 0 {
		ed.removeSelection()
	}
	if ed.replaceMode && (ed.cIdx < len(ed.buff)) && (c == '\n') {
		ed.keyDownPressed()
		ed.homePressed()
		return
	} else if !ed.replaceMode || (ed.cIdx >= len(ed.buff)) || (ed.buff[ed.cIdx] == '\n') {
		ed.buff = append(ed.buff, 0)
		copy(ed.buff[ed.cIdx+1:], ed.buff[ed.cIdx:])
		ed.buff[ed.cIdx] = c
	} else {
		ed.buff[ed.cIdx] = c
	}
	ed.keyRightPressed()
}

func (ed *editor) removeSelection() {
	beg := ed.cIdx
	end := ed.selIdx
	if beg > end {
		beg, end = end, beg
		// need to do it this way to preserve the x, y of the cursor
		for ed.cIdx > beg {
			ed.keyLeftPressed()
		}
	}
	copy(ed.buff[beg:], ed.buff[end:])
	ed.buff = ed.buff[:len(ed.buff)-end+beg]
	ed.selIdx = -1
}

func (ed *editor) copySelection() {
	if ed.selIdx < 0 {
		return
	}
	beg := ed.cIdx
	end := ed.selIdx
	if beg > end {
		beg, end = end, beg
	}
	if beg == end {
		return
	}
	if cap(ed.clipboard) < (end - beg) {
		ed.clipboard = make([]byte, end-beg)
	}
	copy(ed.clipboard, ed.buff[beg:end])
	ed.clipboard = ed.clipboard[:end-beg]
}

func (ed *editor) paste() {
	if len(ed.clipboard) == 0 {
		return
	}
	if ed.selIdx >= 0 {
		ed.removeSelection()
	}
	// allocate space if needed
	ed.buff = append(ed.buff, ed.clipboard...)
	// make a gap
	copy(ed.buff[ed.cIdx+len(ed.clipboard):], ed.buff[ed.cIdx:])
	// fill the gap
	copy(ed.buff[ed.cIdx:], ed.clipboard)
	for i := 0; i < len(ed.clipboard); i++ {
		ed.keyRightPressed()
	}
}

func (ed *editor) backspacePressed() {
	if ed.selIdx >= 0 {
		ed.removeSelection()
		return
	}
	if ed.cIdx <= 0 {
		return
	}
	ed.keyLeftPressed()
	ed.delPressed()
}

func (ed *editor) delPressed() {
	if ed.selIdx >= 0 {
		ed.removeSelection()
		return
	}
	if ed.cIdx < 0 {
		return
	}
	copy(ed.buff[ed.cIdx:], ed.buff[ed.cIdx+1:])
	ed.buff = ed.buff[:len(ed.buff)-1]
}

// Checks if y is off the screen and adjusts ed.ulIdx and y to correct
// as-needed
func (ed *editor) checkScroll(y int) int {
	x := 0
	w, h := ed.txtb.Txtw.TextSize()
	for y < 0 {
		// go back one position
		ed.ulIdx--
		// at this point we are either on a '\n' for the end-of-line
		// case or not (for the wrapping case)
		if ed.buff[ed.ulIdx] == '\n' {
			// we need to count the number of characters to the end of
			// the previous line so we can figure out the overhand of this one
			linelen := 0
			for {
				idx := ed.ulIdx - linelen - 1
				if idx < 0 || ed.buff[idx] == '\n' {
					break
				}
				linelen++
			}
			ed.ulIdx -= linelen % w
			y++
		} else {
			// jump to the start of the line
			ed.ulIdx -= w - 1
			y++
		}
	}
	for y >= h {
		// need to scroll down
		for {
			if x >= w {
				y--
				break
			}
			if ed.buff[ed.ulIdx] == '\n' {
				ed.ulIdx++
				y--
				break
			}
			x++
			ed.ulIdx++
		}
	}
	return y
}
