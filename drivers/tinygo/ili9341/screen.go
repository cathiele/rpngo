package ili9341

import (
	"image/color"
	"mattwach/rpngo/window"

	"tinygo.org/x/drivers/ili9341"
)

type Ili9341Screen struct {
	// Control the LCD.
	Device *ili9341.Device
}

func (s *Ili9341Screen) Init() {
	s.Device = InitDisplay()
}

func (s *Ili9341Screen) NewTextWindow() (window.TextWindow, error) {
	tw := &Ili9341TxtW{} // object allocated on the heap (OK)
	tw.Init(s.Device)
	return tw, nil
}

func (s *Ili9341Screen) NewPixelWindow() (window.PixelWindow, error) {
	pw := &Ili9341PixW{} // object allocated on the heap (OK)
	pw.Init(s.Device)
	return pw, nil
}

func (s *Ili9341Screen) ScreenSize() (int, int) {
	return 320, 240
}

// fastvline and fasthline are bugged for some reason, we we do it the slow way for now
func slowHline(d *ili9341.Device, x0 int16, x1 int16, y int16, c color.RGBA) error {
	for x0 <= x1 {
		d.SetPixel(x0, y, c)
		x0++
	}
	return nil
}

func slowVline(d *ili9341.Device, x int16, y0 int16, y1 int16, c color.RGBA) error {
	for y0 <= y1 {
		d.SetPixel(x, y0, c)
		y0++
	}
	return nil
}
