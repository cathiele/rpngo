// Implements textwindow for the tinygo env.
//
// Currently this targets the ili9341.  If/when more
// devices are supported, some refactoring may need to occcur.
package ili9341tw

import (
	"image/color"
	"mattwach/rpngo/io/drivers/tinygo/pixel565"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/pixel"
	"tinygo.org/x/tinyfont/freemono"
)

type Ili9341TW struct {
	// holds the characters that make up the text grid
	chars []byte

	// screen to send chars to
	device *ili9341.Device

	// upper left corner pixel offset
	x int16
	y int16

	// current character position
	cx int16
	cy int16

	// character width and height in pixels
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
	lastr byte

	// text color
	fgcol color.RGBA
	bgcol pixel.RGB565BE
}

// Init initializes a text window. x, y, w, and h are all in pixels
func (tw *Ili9341TW) Init(d *ili9341.Device, x, y, w, h int) {
	tw.cw = 11
	tw.ch = 16
	tw.cyoffset = 11
	tw.image.Init(tw.cw, tw.ch)
	tw.device = d
	tw.Resize(x, y, w, h)
}

func (tw *Ili9341TW) Resize(x, y, w, h int) error {
	tw.x = int16(x)
	tw.y = int16(y)
	tw.cx = 0
	tw.cy = 0
	tw.textw = int16(w) / tw.cw
	tw.texth = int16(h) / tw.ch
	tw.chars = make([]byte, int(tw.textw)*int(tw.texth))
	tw.Erase()
	return nil
}

func (tw *Ili9341TW) Refresh() {
	// maybe no need to do this?
}

func (tw *Ili9341TW) updateCharAt(tx, ty int16, r byte) {
	oldr := tw.chars[ty*tw.textw+tx]
	if r == oldr {
		return
	}
	tw.chars[ty*tw.textw+tx] = r
	if r != tw.lastr {
		tw.lastr = r
		tw.image.Image.FillSolidColor(tw.bgcol)
		freemono.Regular9pt7b.GetGlyph(rune(r)).Draw(&tw.image, 0, tw.cyoffset, tw.fgcol)
	}
	tw.device.DrawBitmap(tw.x+tx*tw.cw, tw.y+ty*tw.ch, tw.image.Image)
}

func (tw *Ili9341TW) Erase() {
	var j int16
	for j = 0; j < tw.texth; j++ {
		var i int16
		for i = 0; i < tw.textw; i++ {
			tw.updateCharAt(i, j, ' ')
		}
	}
}

func (tw *Ili9341TW) ShowBorder(screenw, screenh int) error {
	// implement later
	return nil
}

func (tw *Ili9341TW) Write(b byte) error {
	if (b == '\n') || (tw.cx >= tw.textw) {
		// next line
		tw.cx = 0
		tw.cy++
	}
	if tw.cy >= tw.texth {
		tw.Scroll(int(tw.texth - tw.cy - 1))
	}
	if b != '\n' {
		tw.updateCharAt(tw.cx, tw.cy, b)
		tw.cx++
	}
	return nil
}

func (tw *Ili9341TW) Width() int {
	return int(tw.textw)
}

func (tw *Ili9341TW) Height() int {
	return int(tw.texth)
}

func (tw *Ili9341TW) Size() (int, int) {
	return int(tw.texth), int(tw.texth)
}

func (tw *Ili9341TW) WindowXY() (int, int) {
	return int(tw.x), int(tw.y)
}

func (tw *Ili9341TW) X() int {
	return int(tw.cx)
}

func (tw *Ili9341TW) Y() int {
	return int(tw.cy)
}

func (tw *Ili9341TW) XY() (int, int) {
	return int(tw.cx), int(tw.cy)
}

func (tw *Ili9341TW) SetX(x int) {
	tw.cx = int16(x)
}

func (tw *Ili9341TW) SetY(y int) {
	tw.cy = int16(y)
}

func (tw *Ili9341TW) SetXY(x, y int) {
	tw.cx = int16(x)
	tw.cy = int16(y)
}

func (tw *Ili9341TW) Color(fr, fg, fb, br, bg, bb int) error {
	tw.fgcol = color.RGBA{
		R: uint8(fr * 8),
		G: uint8(fg * 8),
		B: uint8(fb * 8),
	}
	tw.bgcol = pixel.NewRGB565BE(
		uint8(br*8),
		uint8(bg*8),
		uint8(bb*8),
	)
	// do not reuse image
	tw.lastr = 0xff
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
	for y < tw.texth {
		var x int16
		for x = 0; x < tw.textw; x++ {
			tw.updateCharAt(x, y, ' ')
		}
		y++
	}
}

func (tw *Ili9341TW) scrollDown(i int) {
	// not yet implemented
}

func (tw *Ili9341TW) Cursor(bool) {
	// not implemented yet
}
