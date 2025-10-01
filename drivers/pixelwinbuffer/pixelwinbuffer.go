// Package pixelwinbuffer provides an alternative to
// "erase all pixels, redraw".  Which either creates on-screen
// flashing or takes up memory (significant amount for the rp2040).
// It works like this:
//
// All drawing calls fill in a 1bpp buffer that the pixel is updated
// when it's time to refresh, a previous version of this
// buffer is compared to the latest one.
// Every pixel that is set for the old one and not set
// for the new one is cleared.
//
// This has the effect of a slight double image as the old
// and new will be onscreen at the same time for a very small
// amount of time until the old pixels can be cleared.
package pixelwinbuffer

import (
	"errors"
	"image/color"
	"mattwach/rpngo/window"
)

type PixelBuffer struct {
	previous []uint8
	current  []uint8
	target   window.PixelWindow
	pw       int
	ph       int
}

func (pb *PixelBuffer) Init(target window.PixelWindow) error {
	pb.target = target
	return pb.ResizeWindow(0, 0, 1, 1)
}

func (pb *PixelBuffer) ResizeWindow(x, y, w, h int) error {
	pb.target.ResizeWindow(x, y, w, h)
	size := ((w + 7) * h) / 8
	if size <= 0 {
		return errors.New("pixelbufffer, size <= 0")
	}
	if len(pb.previous) < size {
		pb.previous = make([]uint8, size)
	}
	if len(pb.current) < size {
		pb.current = make([]uint8, size)
	}
	pb.pw, pb.ph = pb.target.PixelSize() // for less overhead
	return nil
}

func (pb *PixelBuffer) ShowBorder(sw, sh int) error {
	w, h := pb.target.WindowSize()
	pb.target.Color(window.BorderColor)
	pb.HLine(-1, -1, w)
	pb.HLine(-1, h-1, w)
	pb.VLine(-1, -1, h)
	pb.VLine(w-1, -1, h)
}

func (pb *PixelBuffer) WindowXY() (int, int) {
	return pb.target.WindowXY()
}

func (pb *PixelBuffer) WindowSize() (int, int) {
	return pb.target.WindowSize()
}

func (pb *PixelBuffer) PixelSize() (int, int) {
	return pb.target.PixelSize()
}

func (pb *PixelBuffer) Color(c color.RGBA) {
	pb.target.Color(c)
}

func (pb *PixelBuffer) setBit(x, y int) {
	poffset := y*pb.pw + x
	pb.current[poffset>>3] |= (0x80 >> (poffset & 7))
}

func (pb *PixelBuffer) setBitHline(x, y, w int) {
	maxy := y + w - 1
	for y <= maxy {
		if ((y & 0x07) == 0) && ((maxy - y) >= 8) {
			poffset := y*pb.pw + x
			pb.current[poffset>>3] = 0xFF
		} else {
			pb.setBit(x, y)
			y++
		}
	}
}

func (pb *PixelBuffer) SetPoint(x, y int) {
	pb.setBit(x, y)
	pb.target.SetPoint(x, y)
}

func (pb *PixelBuffer) HLine(x, y, w int) {
	pb.setBitHline(x, y, w)
	pb.target.HLine(x, y, w)
}
