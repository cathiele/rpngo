// Ili9341PixW implements the PixelWindow interface in
// window/window.go to the ili9341 LCD
package ili9341

import (
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"

	"tinygo.org/x/drivers/ili9341"
)

type Ili9341PixW struct {
	// screen to send chars to
	device *ili9341.Device

	// window dimensions in pixels
	wx int16
	wy int16
	ww int16
	wh int16

	// current color
	col color.RGBA
}

// Init initializes a pixel window. x, y, w, and h are all in pixels
func (tw *Ili9341PixW) Init(d *ili9341.Device, w, h int) {
	tw.device = d
	tw.ResizeWindow(0, 0, w, h)
}

func (tw *Ili9341PixW) ResizeWindow(x, y, w, h int) error {
	if (tw.wx == int16(x)) && (tw.wy == int16(y)) && (tw.ww == int16(w)) && (tw.wh == int16(h)) {
		return nil
	}
	tw.wx = int16(x)
	tw.wy = int16(y)
	tw.device.FillRectangle(int16(x), int16(y), int16(w), int16(h), color.RGBA{})
	return nil
}

func (tw *Ili9341PixW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili9341PixW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}

func (tw *Ili9341PixW) ShowBorder(screenw, screenh int) error {
	c := color.RGBA{R: 100, G: 0, B: 100}
	tw.device.DrawFastHLine(tw.wx, tw.wx+tw.ww-1, tw.wy, c)
	tw.device.DrawFastHLine(tw.wx, tw.wx+tw.ww-1, tw.wy+tw.wh-1, c)
	tw.device.DrawFastVLine(tw.wx, tw.wy, tw.wy+tw.wh-1, c)
	tw.device.DrawFastVLine(tw.wx+tw.ww-1, tw.wy, tw.wy+tw.wh-1, c)
	return nil
}

func (tw *Ili9341PixW) Color(c color.RGBA) {
	tw.col = c
}

func (tw *Ili9341PixW) SetPoint(x, y int) {
	tw.device.SetPixel(int16(x), int16(y), tw.col)
}

func (tw *Ili9341PixW) HLine(x, y, w int) {
	tw.device.DrawFastHLine(int16(x), int16(x+w), int16(y), tw.col)
}

func (tw *Ili9341PixW) VLine(x, y, h int) {
	tw.device.DrawFastVLine(int16(x), int16(y), int16(y+h), tw.col)
}

func (tw *Ili9341PixW) FilledRect(x, y, w, h int) {
	tw.device.FillRectangle(int16(x), int16(y), int16(w), int16(h), tw.col)
}

func (tw *Ili9341PixW) Text(s string, x, y int) {
	// do it lower level to avoid importing a bunch of tinyfont code
	for _, r := range s {
		fonts.NimbusMono12p.GetGlyph(rune(r&0xFF)).Draw(
			tw.device, int16(x), int16(y), tw.col)
		x += FontCharWidth
	}
}
