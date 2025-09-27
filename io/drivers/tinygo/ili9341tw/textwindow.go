// Implements textwindow for the tinygo env.
//
// Currently this targets the ili9341.  If/when more
// devices are supported, some refactoring may need to occcur.
package ili9341tw

import (
	"image/color"
	"mattwach/rpngo/io/drivers/tinygo/fonts"
	"mattwach/rpngo/io/drivers/tinygo/pixel565"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/pixel"
)

// The character and color info packed into 1 16-bit value
//
// Format:
// RGGB RGGB CCCC CCCC
// FFFF BBBB HHHH HHHH
type lcdchar uint16

// rgb are 5 bit values
func newLCDCharFGColor(r, g, b uint16) lcdchar {
	return lcdchar(((r >> 4) << 15) | ((g >> 3) << 13) | ((b >> 4) << 12))
}

// rgb are 5 bit values
func newLCDCharBGColor(r, g, b uint16) lcdchar {
	return lcdchar(((r >> 4) << 11) | ((g >> 3) << 9) | ((b >> 4) << 8))
}

func (l lcdchar) char() byte {
	return byte(l & 0xFF)
}

func (l lcdchar) FGColor() color.RGBA {
	var r uint8 = 0
	if (l & 0x8000) != 0 {
		r = 0xFF
	}
	var g uint8 = 0
	switch l & 0x6000 {
	case 0x6000:
		g = 0xFF
	case 0x4000:
		g = 0xAA
	case 0x2000:
		g = 0x55
	}
	var b uint8 = 0
	if (l & 0x1000) != 0 {
		b = 0xFF
	}
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}
}

func (l lcdchar) BGColor() pixel.RGB565BE {
	var r uint8 = 0
	if (l & 0x800) != 0 {
		r = 0xFF
	}
	var g uint8 = 0
	switch l & 0x600 {
	case 0x600:
		g = 0xFF
	case 0x400:
		g = 0xAA
	case 0x200:
		g = 0x55
	}
	var b uint8 = 0
	if (l & 0x100) != 0 {
		b = 0xFF
	}
	return pixel.NewRGB565BE(r, g, b)
}

type Ili9341TW struct {
	// holds the characters that make up the text grid
	chars []lcdchar

	// screen to send chars to
	device *ili9341.Device

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
	image pixel565.Pixel565
	// saves a little performance in drawing
	lastr lcdchar

	// text color
	fgcol lcdchar
	bgcol lcdchar

	// cusror flash state
	cursorEn bool
	// fg and bg color of original
	cursorCol     lcdchar
	cursorShowing bool
	cursorShowX   int16
	cursorShowY   int16
}

// Init initializes a text window. x, y, w, and h are all in pixels
func (tw *Ili9341TW) Init(d *ili9341.Device, x, y, w, h int) {
	tw.cw = 11
	tw.ch = 16
	tw.cyoffset = 11
	tw.image.Init(tw.cw, tw.ch)
	tw.device = d
	tw.cursorEn = true
	tw.cursorCol = 0x0000
	tw.cursorShowing = false
	tw.ResizeWindow(x, y, w, h)
}

func (tw *Ili9341TW) ResizeWindow(x, y, w, h int) error {
	tw.wx = int16(x)
	tw.wy = int16(y)
	tw.cx = 0
	tw.cy = 0
	tw.textw = int16(w) / tw.cw
	tw.texth = int16(h) / tw.ch
	tw.ww = int16(w)
	tw.wh = int16(h)
	tw.chars = make([]lcdchar, int(tw.textw)*int(tw.texth))
	tw.Erase()
	return nil
}

func (tw *Ili9341TW) Refresh() {
	// maybe no need to do this?
}

func (tw *Ili9341TW) updateCharAt(tx, ty int16, r lcdchar) {
	oldr := tw.chars[ty*tw.textw+tx]
	if r == oldr {
		return
	}
	tw.chars[ty*tw.textw+tx] = r
	if r != tw.lastr {
		tw.lastr = r
		tw.image.Image.FillSolidColor(r.BGColor())
		fonts.NotoMonoRegular8p.GetGlyph(rune(r&0xFF)).Draw(&tw.image, 0, tw.cyoffset, r.FGColor())
	}
	tw.device.DrawBitmap(tw.wx+tx*tw.cw, tw.wy+ty*tw.ch, tw.image.Image)
}

func (tw *Ili9341TW) Erase() {
	var j int16
	b := tw.fgcol | tw.bgcol | lcdchar(' ')
	for j = 0; j < tw.texth; j++ {
		var i int16
		for i = 0; i < tw.textw; i++ {
			tw.updateCharAt(i, j, b)
		}
	}
}

func (tw *Ili9341TW) ShowBorder(screenw, screenh int) error {
	// implement later
	return nil
}

func (tw *Ili9341TW) Write(b byte) error {
	tw.ShowCursorIfEnabled(false)
	if b != '\n' {
		tw.updateCharAt(tw.cx, tw.cy, tw.fgcol|tw.bgcol|lcdchar(b))
		tw.cx++
	}
	if (b == '\n') || (tw.cx >= tw.textw) {
		// next line
		tw.cx = 0
		tw.cy++
	}
	if tw.cy >= tw.texth {
		tw.Scroll(int(tw.texth - tw.cy - 1))
	}
	return nil
}

func (tw *Ili9341TW) TextWidth() int {
	return int(tw.textw)
}

func (tw *Ili9341TW) TextHeight() int {
	return int(tw.texth)
}

func (tw *Ili9341TW) TextSize() (int, int) {
	return int(tw.texth), int(tw.texth)
}

func (tw *Ili9341TW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili9341TW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}

func (tw *Ili9341TW) CursorX() int {
	return int(tw.cx)
}

func (tw *Ili9341TW) CursorY() int {
	return int(tw.cy)
}

func (tw *Ili9341TW) CursorXY() (int, int) {
	return int(tw.cx), int(tw.cy)
}

func (tw *Ili9341TW) SetCursorX(x int) {
	tw.ShowCursorIfEnabled(false)
	tw.cx = int16(x)
}

func (tw *Ili9341TW) SetCursorY(y int) {
	tw.ShowCursorIfEnabled(false)
	tw.cy = int16(y)
}

func (tw *Ili9341TW) SetCursorXY(x, y int) {
	tw.ShowCursorIfEnabled(false)
	tw.cx = int16(x)
	tw.cy = int16(y)
}

func (tw *Ili9341TW) Color(fr, fg, fb, br, bg, bb int) error {
	tw.fgcol = newLCDCharFGColor(uint16(fr), uint16(fg), uint16(fb))
	tw.bgcol = newLCDCharBGColor(uint16(br), uint16(bg), uint16(bb))
	return nil
}

func (tw *Ili9341TW) Scroll(i int) {
	if i < 0 {
		tw.scrollUp(-i)
	} else if i > 0 {
		tw.scrollDown(i)
	}
}

func (tw *Ili9341TW) scrollUp(i int) {
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
			tw.updateCharAt(x, y, tw.chars[offset])
			offset++
		}
	}
	b := tw.fgcol | tw.bgcol | lcdchar(' ')
	for y < tw.texth {
		var x int16
		for x = 0; x < tw.textw; x++ {
			tw.updateCharAt(x, y, b)
		}
		y++
	}
}

func (tw *Ili9341TW) scrollDown(i int) {
	// not yet implemented
}

func (tw *Ili9341TW) Cursor(en bool) {
	tw.ShowCursorIfEnabled(en)
	tw.cursorEn = en
}

func (tw *Ili9341TW) ShowCursorIfEnabled(show bool) {
	if !tw.cursorEn {
		return
	}
	if show == tw.cursorShowing {
		return
	}
	tw.cursorShowing = !tw.cursorShowing
	if show {
		ch := tw.chars[tw.cy*tw.textw+tw.cx]
		tw.cursorCol = ch & 0xFF00
		tw.updateCharAt(tw.cx, tw.cy, 0x0F00|(ch&0x00FF))
		tw.cursorShowX = tw.cx
		tw.cursorShowY = tw.cy
	} else {
		ch := tw.chars[tw.cursorShowY*tw.textw+tw.cursorShowX]
		tw.updateCharAt(tw.cursorShowX, tw.cursorShowY, tw.cursorCol|(ch&0x00FF))
	}
}
