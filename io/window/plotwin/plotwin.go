// Package plotwin shows a 2d plot
package plotwin

import (
	"errors"
	"fmt"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
)

// colors to autorotate through
var colorWheel = []uint16{
	31 << 10,               // red
	31 << 5,                // green
	31,                     // blue
	(31 << 10) | (31 << 5), // yellow
	(31 << 10) | 31,        // magenta
	(31 << 5) | 31,         // cyan
}

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
	txtw     window.TextWindow
	minx     float64
	maxx     float64
	miny     float64
	maxy     float64
	color    uint16
	coloridx int
	autox    bool
	autoy    bool
	minv     float64
	maxv     float64
	steps    uint32
	plots    []Plot
}

func Init(txtw window.TextWindow) (*PlotWindow, error) {
	w := &PlotWindow{
		txtw:  txtw,
		color: 31 << 5, // green
		autox: true,
		autoy: true,
		minv:  -1,
		maxv:  1,
		steps: 250,
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

func (pw *PlotWindow) nextColor() {
	pw.coloridx++
	if pw.coloridx >= len(colorWheel) {
		pw.coloridx = 0
	}
	pw.color = colorWheel[pw.coloridx]
}

func (pw *PlotWindow) AddPlot(r *rpn.RPN, fn []string) error {
	pw.nextColor()
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
