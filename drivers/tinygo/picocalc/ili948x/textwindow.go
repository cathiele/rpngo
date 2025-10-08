// Implements textwindow for the tinygo env.
//
// Currently this targets the ili948x.  If/when more
// devices are supported, some refactoring may need to occcur.
package ili948x

import (
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/window"
)

const textPad = 3

type Ili948xTxtW struct {
	// holds the characters that make up the text grid
	chars []window.ColorChar

	// screen to send chars to
	device *Ili948x

	// window dimensions in pixels
	wx int16
	wy int16
	ww int16
	wh int16

	// character position in text cells
	cx int16
	cy int16

	// character dimension in pixels
	cw int16
	ch int16

	// character y offset
	cyoffset int16

	// width and height as text cells
	textw int16
	texth int16

	// A cell to draw the characters in
	bitmap Bitmap
	// saves a little performance in drawing
	lastr window.ColorChar

	// text color
	col window.ColorChar

	// cusror flash state
	cursorEn bool
	// fg and bg color of original
	cursorCol     window.ColorChar
	cursorShowing bool
	cursorShowX   int16
	cursorShowY   int16
}

// Init initializes a text window. x, y, w, and h are all in pixels
func (tw *Ili948xTxtW) Init(d *Ili948x) {
	tw.cw = FontCharWidth
	tw.ch = 13
	tw.cyoffset = 10
	tw.bitmap.Init(tw.cw, tw.ch)
	tw.device = d
	tw.cursorEn = true
	tw.cursorCol = 0x0000
	tw.cursorShowing = false
	tw.ResizeWindow(0, 0, 1, 1)
}

func (tw *Ili948xTxtW) ResizeWindow(x, y, w, h int) error {
	if (tw.wx == int16(x)) && (tw.wy == int16(y)) && (tw.ww == int16(w)) && (tw.wh == int16(h)) {
		return nil
	}
	tw.wx = int16(x)
	tw.wy = int16(y)
	tw.cx = 0
	tw.cy = 0
	tw.cursorShowing = false
	tw.cursorShowX = 0
	tw.cursorShowY = 0

	if (tw.ww != int16(w)) || (tw.wh != int16(h)) {
		tw.textw = (int16(w) - textPad*2) / tw.cw
		if tw.textw <= 0 {
			tw.textw = 1
		}
		tw.texth = (int16(h) - textPad*2) / tw.ch
		if tw.texth <= 0 {
			tw.texth = 1
		}
		tw.ww = int16(w)
		tw.wh = int16(h)
		tw.chars = make([]window.ColorChar, int(tw.textw)*int(tw.texth))
	}
	tw.device.FillRectangle(int16(x), int16(y), int16(w), int16(h), 0) // fill with black
	var j int16
	b := window.ColorChar(' ')
	for j = 0; j < tw.texth; j++ {
		tw.chars[j] = b
	}
	return nil
}

func (tw *Ili948xTxtW) Refresh() {
	// maybe no need to do this?
}

func fgColor(c window.ColorChar) color.RGBA {
	r, g, b := c.FGColor8()
	return color.RGBA{R: r, G: g, B: b}
}

func bgColor(c window.ColorChar) RGB565 {
	r, g, b := c.BGColor8()
	return NewRGB565(r>>3, g>>2, b>>3)
}

func (tw *Ili948xTxtW) updateCharAt(tx, ty int16, r window.ColorChar) {
	idx := ty*tw.textw + tx
	//tinygo.Check("updateCharAt", int(idx), len(tw.chars))
	oldr := tw.chars[idx]
	if r == oldr {
		return
	}
	tw.chars[idx] = r
	if r != tw.lastr {
		tw.lastr = r
		tw.bitmap.FillWith(bgColor(r))
		fonts.NimbusMono12p.GetGlyph(rune(r&0xFF)).Draw(&tw.bitmap, 0, tw.cyoffset, fgColor(r))
	}
	tw.device.DrawBitmap(tw.wx+tx*tw.cw+textPad, tw.wy+ty*tw.ch+textPad, &tw.bitmap)
}

func (tw *Ili948xTxtW) Erase() {
	var j int16
	b := tw.col | window.ColorChar(' ')
	for j = 0; j < tw.texth; j++ {
		var i int16
		for i = 0; i < tw.textw; i++ {
			tw.updateCharAt(i, j, b)
		}
	}
}

func (tw *Ili948xTxtW) ShowBorder(screenw, screenh int) error {
	tw.device.DrawHLine(tw.wx, tw.wx+tw.ww-1, tw.wy, MAGENTA)
	tw.device.DrawHLine(tw.wx, tw.wx+tw.ww-1, tw.wy+tw.wh-1, MAGENTA)
	tw.device.DrawVLine(tw.wx, tw.wy, tw.wy+tw.wh-1, MAGENTA)
	tw.device.DrawVLine(tw.wx+tw.ww-1, tw.wy, tw.wy+tw.wh-1, MAGENTA)
	return nil
}

func (tw *Ili948xTxtW) Write(b byte) error {
	tw.ShowCursorIfEnabled(false)
	if (b == '\n') || (tw.cx >= tw.textw) {
		// next line
		tw.cx = 0
		tw.cy++
	}
	if tw.cy >= tw.texth {
		tw.Scroll(int(tw.texth - tw.cy - 1))
	}
	if b != '\n' {
		tw.updateCharAt(tw.cx, tw.cy, tw.col|window.ColorChar(b))
		tw.cx++
	}
	return nil
}

func (tw *Ili948xTxtW) TextWidth() int {
	return int(tw.textw)
}

func (tw *Ili948xTxtW) TextHeight() int {
	return int(tw.texth)
}

func (tw *Ili948xTxtW) TextSize() (int, int) {
	return int(tw.textw), int(tw.texth)
}

func (tw *Ili948xTxtW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili948xTxtW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}

func (tw *Ili948xTxtW) CursorX() int {
	return int(tw.cx)
}

func (tw *Ili948xTxtW) CursorY() int {
	return int(tw.cy)
}

func (tw *Ili948xTxtW) CursorXY() (int, int) {
	return int(tw.cx), int(tw.cy)
}

func (tw *Ili948xTxtW) SetCursorX(x int) {
	tw.ShowCursorIfEnabled(false)
	tw.cx = int16(x)
}

func (tw *Ili948xTxtW) SetCursorY(y int) {
	tw.ShowCursorIfEnabled(false)
	tw.cy = int16(y)
}

func (tw *Ili948xTxtW) SetCursorXY(x, y int) {
	tw.ShowCursorIfEnabled(false)
	tw.cx = int16(x)
	tw.cy = int16(y)
}

func (tw *Ili948xTxtW) TextColor(col window.ColorChar) {
	tw.col = col
}

func (tw *Ili948xTxtW) Scroll(i int) {
	if i < 0 {
		tw.scrollUp(-i)
	} else if i > 0 {
		tw.scrollDown(i)
	}
}

func (tw *Ili948xTxtW) scrollUp(i int) {
	if i >= int(tw.texth) {
		tw.Erase()
		tw.cy = 0
		return
	}
	tw.cy -= int16(i)
	maxy := tw.texth - int16(i)
	var y int16
	var offset int = i * int(tw.textw)
	for y = 0; y < maxy; y++ {
		var x int16
		for x = 0; x < tw.textw; x++ {
			//tinygo.Check("scrollUp", int(offset), len(tw.chars))
			tw.updateCharAt(x, y, tw.chars[offset])
			offset++
		}
	}
	b := tw.col | window.ColorChar(' ')
	for y < tw.texth {
		var x int16
		for x = 0; x < tw.textw; x++ {
			tw.updateCharAt(x, y, b)
		}
		y++
	}
}

func (tw *Ili948xTxtW) scrollDown(i int) {
	// not yet implemented
}

func (tw *Ili948xTxtW) Cursor(en bool) {
	tw.ShowCursorIfEnabled(en)
	tw.cursorEn = en
}

func (tw *Ili948xTxtW) ShowCursorIfEnabled(show bool) {
	if !tw.cursorEn {
		return
	}
	if show == tw.cursorShowing {
		return
	}
	tw.cursorShowing = !tw.cursorShowing
	if show {
		if tw.cx >= tw.textw {
			// next line
			tw.cx = 0
			tw.cy++
		}
		if tw.cy >= tw.texth {
			tw.Scroll(int(tw.texth - tw.cy - 1))
		}
		//tinygo.Check("show", int(tw.cy*tw.textw+tw.cx), len(tw.chars))
		ch := tw.chars[tw.cy*tw.textw+tw.cx]
		tw.cursorCol = ch & 0xFF00
		tw.updateCharAt(tw.cx, tw.cy, 0x0F00|(ch&0x00FF))
		tw.cursorShowX = tw.cx
		tw.cursorShowY = tw.cy
	} else {
		//tinygo.Check("hide", int(tw.cursorShowY*tw.textw+tw.cursorShowX), len(tw.chars))
		ch := tw.chars[tw.cursorShowY*tw.textw+tw.cursorShowX]
		tw.updateCharAt(tw.cursorShowX, tw.cursorShowY, tw.cursorCol|(ch&0x00FF))
	}
}
