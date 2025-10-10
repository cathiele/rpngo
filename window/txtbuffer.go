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

	// character position
	cx int16
	cy int16

	// width and height
	w int16
	h int16

	// text color
	col ColorChar
}

func (tb *TextBuffer) UpdateTextWindow(tw TextWindow) {
	tw.SetCursorXY(0, 0)
	var col ColorChar
	tw.TextColor(Black)
	for i := range tb.chars {
		c := tb.chars[i]
		newcol := c & 0xFF00
		if newcol != col {
			col = newcol
			tw.TextColor(col)
		}
		tw.Write(c.Char())
	}
}

func (tb *TextBuffer) MaybeResize(w, h int16) {
	if (tb.w == w) && (tb.h == h) {
		return
	}
	if w <= 0 {
		w = 1
	}
	if h <= 0 {
		h = 1
	}
	tb.w = w
	tb.h = h
	tb.chars = make([]ColorChar, w*h) // object allocated on the heap (OK)
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
		// no scrolling as it may not be needed
		return nil
	}
	if b != '\n' {
		tb.chars[tb.cy*tb.w+tb.cx] = tb.col | ColorChar(b)
		tb.cx++
	}
	return nil
}

func (tb *TextBuffer) TextWidth() int {
	return int(tb.w)
}

func (tb *TextBuffer) TextHeight() int {
	return int(tb.h)
}

func (tb *TextBuffer) TextSize() (int, int) {
	return int(tb.w), int(tb.h)
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
	// maybe not needed so not implemented
}
