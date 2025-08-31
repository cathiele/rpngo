// Package plotwin shows a 2d plot
package plotwin

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
)

type Point struct {
	x     float64
	y     float64
	color uint16
}

type Plot struct {
	fn    []string
	color uint16
}

type PlotWindow struct {
	txtw  window.TextWindow
	minx  float64
	maxx  float64
	miny  float64
	maxy  float64
	color uint16
	autox bool
	autoy bool
	minv  float64
	maxv  float64
	steps uint32
	plots []Plot
}

func Init(txtw window.TextWindow) (*PlotWindow, error) {
	w := &PlotWindow{
		txtw:  txtw,
		color: 31 << 5, // green
		autox: true,
		autoy: true,
		minv:  -1,
		maxv:  1,
		steps: 400,
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
	return "plot"
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

func (pw *PlotWindow) AddPlot(r *rpn.RPN, fn []string) error {
	if len(fn) == 0 {
		return nil
	}
	for i := range pw.plots {
		if slicesAreEqual(fn, pw.plots[i].fn) {
			// plot already exists.  Just update the color
			pw.plots[i].color = pw.color
			return nil
		}
	}
	fncopy := make([]string, len(fn))
	copy(fncopy, fn)
	plot := Plot{fn: fncopy, color: pw.color}

	// do a dry run of creating the points
	_, err := pw.addPoints(r, nil, plot)
	if err != nil {
		return err
	}

	pw.plots = append(pw.plots, plot)
	return nil
}

func slicesAreEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (pw *PlotWindow) Update(r *rpn.RPN) error {
	points, err := pw.createPoints(r)
	if err != nil {
		return err
	}
	pw.txtw.Erase()
	defer pw.txtw.Refresh()
	pw.maybeAutoAdjustAxes(points)
	if err := pw.drawAxis(); err != nil {
		return err
	}
	if err := pw.plotPoints(points); err != nil {
		return err
	}
	return nil
}

func (pw *PlotWindow) createPoints(r *rpn.RPN) ([]Point, error) {
	var points []Point
	for _, plot := range pw.plots {
		var err error
		points, err = pw.addPoints(r, points, plot)
		if err != nil {
			pw.plots = nil
			return nil, fmt.Errorf("plot error for %v, removed all plots: %v", plot.fn, err)
		}
	}
	return points, nil
}

func (pw *PlotWindow) addPoints(r *rpn.RPN, points []Point, plot Plot) ([]Point, error) {
	startlen := r.StackLen()
	step := (pw.maxv - pw.minv) / float64(pw.steps)
	for x := pw.minv; x <= pw.maxv; x += step {
		if err := r.PushComplex(complex(x, 0)); err != nil {
			return nil, err
		}
		if err := r.Exec(plot.fn); err != nil {
			return nil, err
		}
		y, err := r.PopComplex()
		if err != nil {
			return nil, err
		}
		nowlen := r.StackLen()
		if nowlen != startlen {
			return nil, fmt.Errorf(
				"stack changed size running plot string (old: %d, new %d)",
				startlen,
				nowlen)
		}
		points = append(points, Point{x: x, y: real(y), color: plot.color})
	}
	return points, nil
}

func (pw *PlotWindow) maybeAutoAdjustAxes(points []Point) {
	if pw.autox {
		pw.adjustAutoX(points)
	}
	if pw.autoy {
		pw.adjustAutoY(points)
	}
	if pw.autox && pw.autoy {
		pw.makeAxesSquare()
	}
}

func (pw *PlotWindow) adjustAutoX(points []Point) {
	if len(points) == 0 {
		pw.minx = 0
		pw.maxx = 0
	} else {
		pw.minx = points[0].x
		pw.maxx = points[0].x
	}
	for _, p := range points {
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

func (pw *PlotWindow) adjustAutoY(points []Point) {
	if len(points) == 0 {
		pw.miny = 0
		pw.maxy = 0
	} else {
		pw.miny = points[0].y
		pw.maxy = points[0].y
	}
	for _, p := range points {
		if p.y < pw.miny {
			pw.miny = p.y
		} else if p.y > pw.maxy {
			pw.maxy = p.y
		}
	}
	if pw.miny == pw.maxy {
		// create a little spread to avoid math issues
		pw.miny -= 1.0
		pw.maxy += 1.0
	}
}

func (pw *PlotWindow) makeAxesSquare() {
	w, h := pw.txtw.Size()
	wratio := float64(h) / float64(w)
	log.Printf(
		"window ratio: w=%v h=%v minx=%v maxx=%v miny=%v maxy=%v ratio=%v",
		w,
		h,
		pw.minx,
		pw.maxx,
		pw.miny,
		pw.maxy,
		wratio)

	pratio := (pw.maxy - pw.miny) / (pw.maxx - pw.minx)
	log.Printf(
		"starting plot ratio: w=%v h=%v ratio=%v",
		(pw.maxx - pw.minx),
		(pw.maxy - pw.minv),
		pratio)

	if wratio > pratio {
		// need to expand y
		yspread := wratio * (pw.maxx - pw.minx)
		ydelta := (yspread - (pw.maxy - pw.miny)) / 2
		pw.miny -= ydelta
		pw.maxy += ydelta
		pratio = (pw.maxy - pw.miny) / (pw.maxx - pw.minx)
		log.Printf(
			"expandy plot: yspread=%v, ydelta=%v, miny=%v, maxy=%v pratio=%v",
			yspread,
			ydelta,
			pw.miny,
			pw.maxy,
			pratio)
	} else {
		// need to expand x
		xspread := (pw.maxy - pw.miny) / wratio
		xdelta := (xspread - (pw.maxx - pw.minx)) / 2
		pw.minx -= xdelta
		pw.maxx += xdelta
		pratio = (pw.maxy - pw.miny) / (pw.maxx - pw.minx)
		log.Printf(
			"expandx plot: xspread=%v, xdelta=%v, minx=%v, maxx=%v miny=%v maxy=%v pratio=%v",
			xspread,
			xdelta,
			pw.minx,
			pw.maxx,
			pw.miny,
			pw.maxy,
			pratio)
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

func (pw *PlotWindow) plotPoints(points []Point) error {
	var color uint16
	for _, p := range points {
		if p.color != color {
			color = p.color
			pw.txtw.Color(int(color>>10), int((color>>5)&31), int(color&31), 0, 0, 0)
		}
		x, xok := pw.transformX(p.x)
		if !xok {
			continue
		}
		y, yok := pw.transformY(p.y)
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
	return pw.txtw.Height() - int(float64(pw.txtw.Height())*y) - 1, true
}
