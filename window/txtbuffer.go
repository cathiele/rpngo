package window

// A text buffer is intended to avoid the flash effect of erasing, then
// redrawing text.  Due to some smarts of not redrawing the same character,
// it's also faster.
//
// Usage: instead of writing text directly to a text window, write it
// to a text buffer, then have the buffer write to the
type TextBuffer struct {
	// holds the characters that make up the text grid
	chars []ColorChar
	// to make scrolling cheaper, we have a head index which starts
	// at 0 and moves forward when we scroll up
	headidx int

	// character position
	cx int16
	cy int16

	// size of chars buffer
	w int16
	h int16

	// text color
	col ColorChar

	// text area to update
	Txtw TextWindow
}

func (tb *TextBuffer) Init(txtw TextWindow) {
	tb.Txtw = txtw

}

func (tb *TextBuffer) Update() {
	var x int16
	var y int16
	var i int = tb.headidx
	for y = 0; y < tb.h; y++ {
		for x = 0; x < tb.w; x++ {
			if tb.chars[i].IsDirty() {
				tb.chars[i].ClearDirty()
				tb.Txtw.DrawChar(int(x), int(y), tb.col)
			}
			i++
			if i >= len(tb.chars) {
				i = 0
			}
		}
	}
}

func (tb *TextBuffer) MaybeResize() {
	tw, th := tb.Txtw.TextSize()
	if (tw == int(tb.w)) && (th == int(tb.h)) {
		return
	}
	tb.w = int16(tw)
	tb.h = int16(th)
	tb.chars = make([]ColorChar, tb.w*tb.h) // object allocated on the heap (OK)
	tb.Erase()
}

func (tb *TextBuffer) Erase() {
	b := tb.col | ColorChar(' ')
	for i := range tb.chars {
		tb.chars[i] = b
	}
}

func (tb *TextBuffer) Write(b byte) error {
	if (b == '\n') || (tb.cx >= tb.w) {
		// next line
		tb.cx = 0
		tb.cy++
	}
	if tb.cy >= tb.h {
		tb.Scroll(-1)
	}
	if b != '\n' {
		idx := (tb.headidx + int(tb.cy*tb.w) + int(tb.cx)) % len(tb.chars)
		tb.chars[idx] = tb.col | ColorChar(b) | 0x80
		tb.cx++
	}
	return nil
}

func (tb *TextBuffer) CursorX() int {
	return int(tb.cx)
}

func (tb *TextBuffer) CursorY() int {
	return int(tb.cy)
}

func (tb *TextBuffer) CursorXY() (int, int) {
	return int(tb.cx), int(tb.cy)
}

func (tb *TextBuffer) SetCursorX(x int) {
	tb.cx = int16(x)
}

func (tb *TextBuffer) SetCursorY(y int) {
	tb.cy = int16(y)
}

func (tb *TextBuffer) SetCursorXY(x, y int) {
	tb.cx = int16(x)
	tb.cy = int16(y)
}

func (tb *TextBuffer) TextColor(col ColorChar) {
	tb.col = col
}

func (tb *TextBuffer) Scroll(i int) {
	tb.Txtw.Scroll(i)
	oldhead := tb.headidx
	tb.headidx -= i * int(tb.w)
	for tb.headidx > len(tb.chars) {
		tb.headidx -= len(tb.chars)
	}
	for tb.headidx < 0 {
		tb.headidx += len(tb.chars)
	}
	b := tb.col | ColorChar(' ')
	if i < 0 {
		for oldhead != tb.headidx {
			tb.chars[oldhead] = b
			oldhead++
			if oldhead >= len(tb.chars) {
				oldhead = 0
			}
		}
	} else {
		for oldhead != tb.headidx {
			tb.chars[oldhead] = b
			oldhead--
			if oldhead < 0 {
				oldhead = len(tb.chars) - 1
			}
		}
	}
}
