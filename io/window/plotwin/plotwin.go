// Package plotwin shows a 2d plot
package plotwin

import (
	"errors"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
)

type Point struct {
	x     float64
	y     float64
	color uint16
}

type PlotWindow struct {
	txtw   window.TextWindow
	minx   float64
	maxx   float64
	miny   float64
	maxy   float64
	color  uint16
	autox  bool
	autoy  bool
	points []Point
}

func Init(txtw window.TextWindow) (*PlotWindow, error) {
	w := &PlotWindow{
		txtw:  txtw,
		minx:  -1.0,
		maxx:  1.0,
		color: 31 << 5, // green
		autoy: true,
	}
	if err := txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return nil, err
	}
	return w, nil
}

func (pw *PlotWindow) Resize(x, y, w, h int) error {
	return pw.txtw.Resize(x, y, w, h)
}

func (pw *PlotWindow) ShowBorder(screenw, screenh int) error {
	return pw.txtw.ShowBorder(screenw, screenh)
}

func (pw *PlotWindow) WindowXY() (int, int) {
	return pw.txtw.WindowXY()
}

func (pw *PlotWindow) Size() (int, int) {
	return pw.txtw.Size()
}

func (pw *PlotWindow) Type() string {
	return "var"
}

func (pw *PlotWindow) SetProp(name string, val rpn.Frame) error {
	return errors.New("props not supported")
}

func (pw *PlotWindow) GetProp(name string) (rpn.Frame, error) {
	return rpn.Frame{}, errors.New("props not supported")
}

func (pw *PlotWindow) ListProps() []string {
	return nil
}

func (pw *PlotWindow) SetPoint(x, y float64) {
	pw.points = append(pw.points, Point{x, y, pw.color})
}

func (pw *PlotWindow) Update(rpn *rpn.RPN) error {
	pw.txtw.Erase()
	defer pw.txtw.Refresh()
	if pw.autox {
		if len(pw.points) == 0 {
			return nil
		}
		pw.adjustAutoX()
	}
	if pw.autox {
		if len(pw.points) == 0 {
			return nil
		}
		pw.adjustAutoY()
	}
	if err := pw.drawAxis(); err != nil {
		return err
	}
	if err := pw.plotPoints(); err != nil {
		return err
	}
	return nil
}

func (pw *PlotWindow) adjustAutoX() {
	pw.minx = pw.points[0].x
	pw.maxx = pw.points[0].x
	for _, p := range pw.points {
		if p.x < pw.minx {
			pw.minx = p.x
		} else if p.x > pw.maxx {
			pw.maxx = p.x
		}
	}
	if pw.minx == pw.maxx {
		// create a little spread to avoid math issues
		pw.minx -= 1.0
		pw.maxx += 1.0
	}
}

func (pw *PlotWindow) adjustAutoY() {
	pw.miny = pw.points[0].y
	pw.maxy = pw.points[0].y
	for _, p := range pw.points {
		if p.y < pw.miny {
			pw.miny = p.y
		} else if p.y > pw.maxy {
			pw.maxy = p.y
		}
	}
	if pw.minx == pw.maxx {
		// create a little spread to avoid math issues
		pw.miny -= 1.0
		pw.maxy += 1.0
	}
}

func (pw *PlotWindow) drawAxis() error {
	if err := pw.txtw.Color(31, 31, 31, 0, 0, 0); err != nil {
		return err
	}

	x, xok := pw.transformX(0)
	if xok {
		if err := pw.drawVerticalAxis(x); err != nil {
			return err
		}
	}

	y, yok := pw.transformY(0)
	if yok {
		if err := pw.drawHorizontalAxis(y); err != nil {
			return err
		}
	}

	if xok && yok {
		pw.txtw.SetXY(x, y)
		if err := pw.txtw.Write('+'); err != nil {
			return err
		}
	}

	return nil
}

func (pw *PlotWindow) drawVerticalAxis(x int) error {
	h := pw.txtw.Height()
	for y := 0; y < h; y++ {
		pw.txtw.SetXY(x, y)
		if err := pw.txtw.Write('|'); err != nil {
			return err
		}
	}
	return nil
}

func (pw *PlotWindow) drawHorizontalAxis(y int) error {
	w := pw.txtw.Width()
	pw.txtw.SetXY(0, y)
	for x := 0; x < w; x++ {
		if err := pw.txtw.Write('-'); err != nil {
			return err
		}
	}
	return nil
}

func (pw *PlotWindow) plotPoints() error {
	var color uint16
	for _, p := range pw.points {
		if p.color != color {
			color = p.color
			pw.txtw.Color(int(color>>10), int((color>>5)&31), int(color&31), 0, 0, 0)
		}
		x, xok := pw.transformX(p.x)
		if !xok {
			continue
		}
		y, yok := pw.transformY(p.x)
		if !yok {
			continue
		}
		pw.txtw.SetXY(x, y)
		if err := pw.txtw.Write('*'); err != nil {
			return err
		}
	}
	return nil
}

func (pw *PlotWindow) transformX(x float64) (int, bool) {
	x -= pw.minx
	if x < 0 {
		// off the left of the screen
		return 0, false
	}
	// convert x to a ratio between 0 and 1
	x = x / (pw.maxx - pw.minx)
	if x >= 1 {
		// off the right of the screen
		return 0, false
	}
	return int(float64(pw.txtw.Width()) * x), true
}

func (pw *PlotWindow) transformY(y float64) (int, bool) {
	y -= pw.miny
	if y < 0 {
		// off the top of the screen
		return 0, false
	}
	// convert y to a ratio between 0 and 1
	y = y / (pw.maxy - pw.miny)
	if y >= 1 {
		// off the bottom of the screen
		return 0, false
	}
	return int(float64(pw.txtw.Height()) * y), true
}
