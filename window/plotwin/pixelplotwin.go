package plotwin

import (
	"image/color"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
)

// colors to autorotate through
var colorWheelPixel = []color.RGBA{
	{R: 255, G: 0, B: 0, A: 255},
	{R: 0, G: 255, B: 0, A: 255},
	{R: 0, G: 0, B: 255, A: 255},
	{R: 255, G: 255, B: 0, A: 255},
	{R: 255, G: 0, B: 255, A: 255},
	{R: 0, G: 255, B: 255, A: 255},
}

type PixelPlotWindow struct {
	pixw       window.PixelWindow
	common     plotWindowCommon
	lastcolidx uint8
}

func AddPixelPlotFn(w window.WindowWithProps, r *rpn.RPN, fn []string, isParametric bool) error {
	return w.(*PixelPlotWindow).common.addPlot(r, fn, isParametric, uint8(len(colorWheelPixel)))
}

func (pw *PixelPlotWindow) Init(pixw window.PixelWindow) {
	pw.pixw = pixw
	pw.common.init()
}

func (pw *PixelPlotWindow) ResizeWindow(x, y, w, h int) error {
	return pw.pixw.ResizeWindow(x, y, w, h)
}

func (pw *PixelPlotWindow) ShowBorder(screenw, screenh int) error {
	return pw.pixw.ShowBorder(screenw, screenh)
}

func (pw *PixelPlotWindow) WindowXY() (int, int) {
	return pw.pixw.WindowXY()
}

func (pw *PixelPlotWindow) WindowSize() (int, int) {
	return pw.pixw.WindowSize()
}

func (pw *PixelPlotWindow) Type() string {
	return "plot"
}

func (pw *PixelPlotWindow) Update(r *rpn.RPN) error {
	w, h := pw.pixw.WindowSize()
	pw.pixw.Color(color.RGBA{})
	pw.pixw.FilledRect(0, 0, w, h)
	pw.drawAxis()
	pw.lastcolidx = 255
	return pw.common.createPoints(r, pw.plotPoint)
}

func (pw *PixelPlotWindow) SetProp(name string, val rpn.Frame) error {
	return pw.common.setProp(name, val)
}

func (pw *PixelPlotWindow) GetProp(name string) (rpn.Frame, error) {
	return pw.common.getProp(name)
}

func (pw *PixelPlotWindow) ListProps() []string {
	return pw.common.ListProps()
}

func (pw *PixelPlotWindow) plotPoint(x, y float64, colidx uint8) error {
	w, h := pw.pixw.WindowSize()
	if colidx != pw.lastcolidx {
		pw.lastcolidx = colidx
		pw.pixw.Color(colorWheelPixel[colidx])
	}
	wx, xok := pw.common.transformX(x, w)
	if !xok {
		return nil
	}
	wy, yok := pw.common.transformY(y, h)
	if !yok {
		return nil
	}
	pw.pixw.SetPoint(wx, wy)
	return nil
}
