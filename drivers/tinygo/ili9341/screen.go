package ili9341

import (
	"mattwach/rpngo/window"

	"tinygo.org/x/drivers/ili9341"
)

const FontCharWidth = 8

type Ili9341Screen struct {
	// Control the LCD.
	Device *ili9341.Device
}

func (s *Ili9341Screen) Init() {
	s.Device = InitDisplay()
}

func (s *Ili9341Screen) NewTextWindow() (window.TextWindow, error) {
	tw := &Ili9341TxtW{}
	tw.Init(s.Device)
	return tw, nil
}

func (s *Ili9341Screen) NewPixelWindow() (window.PixelWindow, error) {
	pw := &Ili9341PixW{}
	pw.Init(s.Device)
	return pw, nil
}

func (s *Ili9341Screen) ScreenSize() (int, int) {
	return 320, 240
}
