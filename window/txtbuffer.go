package window

// TextBuffer serves multiple purposes and allows character drivers (of which
// we have several LCD and curses variants) have a reduced interface and
// implementation requirements.  This also increases consistency (by does
// involve reinventing already-solved issues for some targets - particurlarily
// the ncurses one).
//
// Base requirements are:
//  1. Support color text.  Maybe backgrounds too, although that had not yet
//     seen any use
//  2. Do not redraw characters with the same character.  Drawing characters
//     on a color LCD via SPI is a somehwat expensive process.  A low hanginng
//     fruit optimization is to not do it when possible.
//  3. Support scrolling.  It is reasonable for the user to need to look back at
//     history, especially when a command (such as help) has a multiscreen
//     result.
//  5. Keepthe door open for editing support.  Although that will come later.
//
// Implementation.
//
// We have a char, which is a 16-bit value that combines:
//   - An 8-bit ASCII byte (expansion to unicode might be needed later, if this
//     project takes on any international interet).
//   - A foreground and background color (4 bits each and why we might drop
//     or reduce the bg bits later.
//
// We have a chars slice that contains the entire buffer as a single block.
// The format is subsequent bytes tracing to the right and wrapping to the next
// line, which is a typical implementation.
//
// We have headidx which is a pointer to 0,0 of the visible area of the LCD
// Characters before headidx are what you would see if you scrollup and
// characters after headix+screensize are what you will see if you scroll down.
type TextBuffer struct {
	// holds the characters that make up the entire text area
	buffer []ColorChar
	// holds the characters the represent the visible screen
	screen []ColorChar

	// to make scrolling cheaper, we have a head index which starts
	// at 0 and moves forward when we scroll up
	headidx int
	// Scrollback bytes is used to determine how large to make chars and
	// h.  Set it to zero for no scrolling support
	scrollbytes int

	// character position, on the screen and not the buffer
	cx int16
	cy int16

	// height of chars buffer.  Note that bw should equal the TextWindow
	// width but ch could be larger.
	bw int16
	bh int16

	// text color
	col ColorChar

	// text area to update
	Txtw TextWindow
}

func (tb *TextBuffer) Init(txtw TextWindow, scrollbytes int) error {
	tb.Txtw = txtw
	tb.scrollbytes = scrollbytes
	x, y := txtw.WindowXY()
	w, h := txtw.WindowSize()
	return tb.ResizeWindow(x, y, w, h)
}

func (tb *TextBuffer) Update() {
	var bi int = tb.headidx
	var si int = 0
	w, h := tb.Txtw.TextSize()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if tb.screen[si] != tb.buffer[bi] {
				tb.screen[si] = tb.buffer[bi]
				tb.Txtw.DrawChar(int(x), int(y), tb.screen[si])
			}
			bi = (bi + 1) % len(tb.buffer)
			si++
		}
	}
	tb.Txtw.Refresh() // ncurses needs this, LCDs do not
}

func (tb *TextBuffer) ResizeWindow(x, y, w, h int) error {
	if err := tb.Txtw.ResizeWindow(x, y, w, h); err != nil {
		return err
	}
	tw, th := tb.Txtw.TextSize()
	scrollh := tb.scrollbytes / tw
	if (int(tb.bw) == tw) && (int(tb.bh) == (scrollh + th)) {
		// already the right size
		return nil
	}
	tb.bw = int16(tw)
	tb.bh = int16(th + scrollh)
	tb.buffer = make([]ColorChar, tb.bw*tb.bh) // object allocated on the heap (OK)
	tb.screen = make([]ColorChar, tw*th)       // object allocated on the heap (OK)
	// maybe we can reflow the text instead of erasing it after the changes
	// are proven as stable.
	tb.Erase()
	return nil
}

func (tb *TextBuffer) Erase() {
	b := tb.col | ColorChar(' ')
	for i := range tb.buffer {
		tb.buffer[i] = b
	}
	tb.headidx = 0
	tb.Update()
}

func (tb *TextBuffer) Write(b byte, updatenow bool) error {
	tw, th := tb.Txtw.TextSize()
	if (b == '\n') || (int(tb.cx) >= tw) {
		// next line
		tb.cx = 0
		tb.cy++
	}
	scrolled := false
	if int(tb.cy) >= th {
		tb.Scroll(1)
		tb.cy--
		scrolled = true
	}
	if b != '\n' {
		bidx := (tb.headidx + int(tb.cy*tb.bw) + int(tb.cx)) % len(tb.buffer)
		tb.buffer[bidx] = tb.col | ColorChar(b)
		if scrolled {
			// erase rest of the line, which might contain old data
			for i := 1; i < tw; i++ {
				tb.buffer[bidx+i] = tb.col | ColorChar(' ')
			}
		}
		if updatenow {
			if scrolled {
				// need to update the entire screen
				tb.Update()
			} else {
				// just update this character
				sidx := tb.cy*tb.bw + tb.cx
				if tb.screen[sidx] != tb.buffer[bidx] {
					tb.screen[sidx] = tb.buffer[bidx]
					tb.Txtw.DrawChar(int(tb.cx), int(tb.cy), tb.screen[sidx])
				}
			}
		}
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
	// scrolling is as easy as moving the headidx
	tb.headidx += i * int(tb.bw)
	for tb.headidx >= len(tb.buffer) {
		tb.headidx -= int(tb.bw)
	}
	for tb.headidx < 0 {
		tb.headidx += int(tb.bw)
	}
}
