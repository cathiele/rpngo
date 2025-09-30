// Package plotwin shows a 2d plot
package plotwin

import (
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
)

// colors to autorotate through
var colorWheelTxt = []window.ColorChar{
	window.Red,
	window.Green,
	window.Blue,
	window.Yellow,
	window.Magenta,
	window.Cyan,
}

type TxtPlotWindow struct {
	txtw   window.TextWindow
	common plotWindowCommon
}

func AddTxtPlotFn(w window.WindowWithProps, r *rpn.RPN, fn []string, isParametric bool) error {
	return w.(*TxtPlotWindow).common.addPlot(r, fn, isParametric, uint8(len(colorWheelTxt)))
}

func (pw *TxtPlotWindow) Init(txtw window.TextWindow) {
	pw.txtw = txtw
	pw.common.init()
}

func (pw *TxtPlotWindow) ResizeWindow(x, y, w, h int) error {
	return pw.txtw.ResizeWindow(x, y, w, h)
}

func (pw *TxtPlotWindow) ShowBorder(screenw, screenh int) error {
	return pw.txtw.ShowBorder(screenw, screenh)
}

func (pw *TxtPlotWindow) WindowXY() (int, int) {
	return pw.txtw.WindowXY()
}

func (pw *TxtPlotWindow) WindowSize() (int, int) {
	return pw.txtw.WindowSize()
}

func (pw *TxtPlotWindow) Type() string {
	return "plot"
}

func (pw *TxtPlotWindow) Update(r *rpn.RPN) error {
	points, err := pw.common.createPoints(r)
	if err != nil {
		return err
	}
	pw.txtw.Erase()
	defer pw.txtw.Refresh()
	if err := pw.drawAxis(); err != nil {
		return err
	}
	if err := pw.plotPoints(points); err != nil {
		return err
	}
	return nil
}

func (pw *TxtPlotWindow) SetProp(name string, val rpn.Frame) error {
	return pw.common.setProp(name, val)
}

func (pw *TxtPlotWindow) GetProp(name string) (rpn.Frame, error) {
	return pw.common.getProp(name)
}

func (pw *TxtPlotWindow) ListProps() []string {
	return pw.common.ListProps()
}

func (pw *TxtPlotWindow) plotPoints(points []Point) error {
	var lastcolidx uint8 = 255
	w := pw.txtw.TextWidth()
	h := pw.txtw.TextHeight()
	for _, p := range points {
		if p.coloridx != lastcolidx {
			lastcolidx = p.coloridx
			pw.txtw.TextColor(colorWheelTxt[lastcolidx])
		}
		x, xok := pw.common.transformX(p.x, w)
		if !xok {
			continue
		}
		y, yok := pw.common.transformY(p.y, h)
		if !yok {
			continue
		}
		pw.txtw.SetCursorXY(x, y)
		if err := pw.txtw.Write('*'); err != nil {
			return err
		}
	}
	return nil
}
