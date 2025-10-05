// Ili948xPixW implements the PixelWindow interface in
// window/window.go to the ili948x LCD
package ili948x

import (
	"errors"
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/window"
)

type Ili948xPixW struct {
	// screen to send chars to
	device *Ili948x

	// window dimensions in pixels
	wx int16
	wy int16
	ww int16
	wh int16

	// pixel dimension, in pixels
	px int16
	py int16
	pw int16
	ph int16

	// current color
	col color.RGBA
}

// Init initializes a pixel window. x, y, w, and h are all in pixels
func (tw *Ili948xPixW) Init(d *Ili948x) {
	tw.device = d
	tw.ResizeWindow(0, 0, 5, 5)
}

func (tw *Ili948xPixW) Refresh() {}

func (tw *Ili948xPixW) ResizeWindow(x, y, w, h int) error {
	if (tw.wx == int16(x)) && (tw.wy == int16(y)) && (tw.ww == int16(w)) && (tw.wh == int16(h)) {
		return nil
	}
	if (w < 3) || (h < 3) {
		return errors.New("pixelwindow resize too small")
	}
	tw.wx = int16(x)
	tw.wy = int16(y)
	tw.ww = int16(w)
	tw.wh = int16(h)
	tw.px = tw.wx + 1
	tw.py = tw.wy + 1
	tw.pw = tw.ww - 2
	tw.ph = tw.wh - 2
	tw.device.FillRectangle(int16(x), int16(y), int16(w), int16(h), color.RGBA{})
	return nil
}

func (tw *Ili948xPixW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili948xPixW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}

func (tw *Ili948xPixW) PixelSize() (int, int) {
	return int(tw.pw), int(tw.ph)
}

func (tw *Ili948xPixW) ShowBorder(screenw, screenh int) error {
	// Need to expand by one pixel as the main window does not include the border.
	x0 := tw.wx
	x1 := tw.wx + tw.ww - 1
	y0 := tw.wy
	y1 := tw.wy + tw.wh - 1
	tw.device.DrawHLine(x0, x1, y0, window.BorderColor)
	tw.device.DrawHLine(x0, x1, y1, window.BorderColor)
	tw.device.DrawVLine(x0, y0+1, y1-1, window.BorderColor)
	tw.device.DrawVLine(x1, y0+1, y1-1, window.BorderColor)
	return nil
}

func (tw *Ili948xPixW) Color(c color.RGBA) {
	tw.col = c
}

func (tw *Ili948xPixW) SetPoint(x, y int) {
	tw.device.SetPixel(tw.px+int16(x), tw.py+int16(y), tw.col)
}

func (tw *Ili948xPixW) HLine(x, y, w int) {
	tw.device.DrawHLine(tw.px+int16(x), tw.px+int16(x+w)-1, tw.py+int16(y), tw.col)
}

func (tw *Ili948xPixW) VLine(x, y, h int) {
	tw.device.DrawVLine(tw.px+int16(x), tw.py+int16(y), tw.py+int16(y+h)-1, tw.col)
}

func (tw *Ili948xPixW) FilledRect(x, y, w, h int) {
	tw.device.FillRectangle(tw.px+int16(x), tw.py+int16(y), int16(w), int16(h), tw.col)
}

func (tw *Ili948xPixW) Text(s string, x, y int) {
	// do it lower level to avoid importing a bunch of tinyfont code
	for _, r := range s {
		fonts.NimbusMono12p.GetGlyph(rune(r&0xFF)).Draw(
			tw.device, tw.px+int16(x), tw.py+int16(y), tw.col)
		x += FontCharWidth
	}
}
