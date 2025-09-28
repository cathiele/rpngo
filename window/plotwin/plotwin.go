// Package plotwin shows a 2d plot
package plotwin

import (
	"fmt"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
)

// colors to autorotate through
var colorWheel = []window.ColorChar{
	window.Red,
	window.Green,
	window.Blue,
	window.Yellow,
	window.Magenta,
	window.Cyan,
}

type Point struct {
	x     float64
	y     float64
	color window.ColorChar
}

type Plot struct {
	fn           []string
	color        window.ColorChar
	isParametric bool
}

type PlotWindow struct {
	txtw     window.TextWindow
	minx     float64
	maxx     float64
	miny     float64
	maxy     float64
	color    window.ColorChar
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
	txtw.TextColor(window.White)
	return w, nil
}

func (pw *PlotWindow) ResizeWindow(x, y, w, h int) error {
	return pw.txtw.ResizeWindow(x, y, w, h)
}

func (pw *PlotWindow) ShowBorder(screenw, screenh int) error {
	return pw.txtw.ShowBorder(screenw, screenh)
}

func (pw *PlotWindow) WindowXY() (int, int) {
	return pw.txtw.WindowXY()
}

func (pw *PlotWindow) WindowSize() (int, int) {
	return pw.txtw.WindowSize()
}

func (pw *PlotWindow) Type() string {
	return "plot"
}

func (pw *PlotWindow) nextColor() {
	pw.coloridx++
	if pw.coloridx >= len(colorWheel) {
		pw.coloridx = 0
	}
	pw.color = colorWheel[pw.coloridx]
}

func (pw *PlotWindow) AddPlot(r *rpn.RPN, fn []string, isParametric bool) error {
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
	plot := Plot{fn: fncopy, color: pw.color, isParametric: isParametric}

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
	var x complex128
	for v := pw.minv; v <= pw.maxv; v += step {
		if err := r.PushComplex(complex(v, 0)); err != nil {
			return nil, err
		}
		if err := r.Exec(plot.fn); err != nil {
			return nil, err
		}
		y, err := r.PopComplex()
		if err != nil {
			return nil, err
		}
		if plot.isParametric {
			x, err = r.PopComplex()
			if err != nil {
				return nil, err
			}
		} else {
			x = complex(v, 0)
		}
		nowlen := r.StackLen()
		if nowlen != startlen {
			return nil, fmt.Errorf(
				"stack changed size running plot string (old: %d, new %d)",
				startlen,
				nowlen)
		}
		points = append(points, Point{x: real(x), y: real(y), color: plot.color})
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
