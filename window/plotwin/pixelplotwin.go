package plotwin

import (
	"image/color"
	"mattwach/rpngo/window"
)

// colors to autorotate through
var colorWheelPixel = []color.RGBA{
	color.RGBA{R: 255, G: 0, B: 0, A: 255},
	color.RGBA{R: 0, G: 255, B: 0, A: 255},
	color.RGBA{R: 0, G: 0, B: 255, A: 255},
	color.RGBA{R: 255, G: 255, B: 0, A: 255},
	color.RGBA{R: 255, G: 0, B: 255, A: 255},
	color.RGBA{R: 0, G: 255, B: 255, A: 255},
}

type PixelPlotWindow struct {
	pixw   window.PixelWindow
	common PlotWindowCommon
}
