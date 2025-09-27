package ili9341tw

import (
	"mattwach/rpngo/io/window"

	"tinygo.org/x/drivers/ili9341"
)

type Ili9341Screen struct {
	// Control the LCD.
	Device *ili9341.Device
}

func (s *Ili9341Screen) Init() {
	s.Device = InitDisplay()
}

// NewTextWindow creates a new text window on the given screen.
// x, y, w, and h are all in pixels.
func (s *Ili9341Screen) NewTextWindow(x, y, w, h int) (window.TextWindow, error) {
	tw := &Ili9341TW{}
	tw.Init(s.Device, x, y, w, h)
	return tw, nil
}

func (s *Ili9341Screen) ScreenSize() (int, int) {
	return 320, 240
}
