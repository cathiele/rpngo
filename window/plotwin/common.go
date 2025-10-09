package plotwin

import (
	"fmt"
	"mattwach/rpngo/rpn"
)

type PointStats struct {
	minx        float64
	maxx        float64
	miny        float64
	maxy        float64
	initialized bool
}

func (ps *PointStats) reset() {
	ps.initialized = false
}

func (ps *PointStats) update(x, y float64, colidx uint8) error {
	if !ps.initialized {
		ps.initialized = true
		ps.minx = x
		ps.maxx = x
		ps.miny = y
		ps.maxy = y
	} else {
		if x < ps.minx {
			ps.minx = x
		}
		if x > ps.maxx {
			ps.maxx = x
		}
		if y < ps.miny {
			ps.miny = y
		}
		if y > ps.maxy {
			ps.maxy = y
		}
	}
	return nil
}

type Plot struct {
	fn           []string
	coloridx     uint8
	isParametric bool
}

type plotWindowCommon struct {
	minx      float64
	maxx      float64
	miny      float64
	maxy      float64
	coloridx  uint8
	autox     bool
	autoy     bool
	minv      float64
	maxv      float64
	steps     uint32
	autosteps uint32
	plots     []Plot
	stats     PointStats
}

func (pw *plotWindowCommon) init() {
	pw.autox = true
	pw.autoy = true
	pw.minv = -1
	pw.maxv = 1
	pw.autosteps = 50
	pw.steps = 250
}

func (pw *plotWindowCommon) nextColor(numColors uint8) {
	pw.coloridx++
	if pw.coloridx >= numColors {
		pw.coloridx = 0
	}
}

func (pw *plotWindowCommon) addPlot(r *rpn.RPN, fn []string, isParametric bool, numColors uint8) error {
	pw.nextColor(numColors)
	if len(fn) == 0 {
		return nil
	}
	for i := range pw.plots {
		if slicesAreEqual(fn, pw.plots[i].fn) {
			// plot already exists.  Just update the color
			pw.plots[i].coloridx = pw.coloridx
			return nil
		}
	}
	fncopy := make([]string, len(fn)) // object allocated on the heap (OK)
	copy(fncopy, fn)
	plot := Plot{fn: fncopy, coloridx: pw.coloridx, isParametric: isParametric}

	// do a dry run of creating the points
	err := pw.addPoints(r, plot, pw.steps, func(x, y float64, c uint8) error { return nil })
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

func (pw *plotWindowCommon) setAxisMinMax(r *rpn.RPN) error {
	// first determine the ranges
	if pw.autox || pw.autoy {
		pw.stats.reset()
		for _, plot := range pw.plots {
			if err := pw.addPoints(r, plot, pw.autosteps, pw.stats.update); err != nil {
				pw.plots = nil
				return fmt.Errorf("plot error for %v, removed all plots: %v", plot.fn, err) // object allocated on the heap (OK)
			}
		}
		if pw.autox {
			pw.adjustAutoX()
		}
		if pw.autoy {
			pw.adjustAutoY()
		}
	}
	return nil
}

func (pw *plotWindowCommon) createPoints(r *rpn.RPN, fn func(x, y float64, coloridx uint8) error) error {
	for _, plot := range pw.plots {
		if err := pw.addPoints(r, plot, pw.steps, fn); err != nil {
			pw.plots = nil
			return fmt.Errorf("plot error for %v, removed all plots: %v", plot.fn, err) // object allocated on the heap (OK)
		}
	}
	return nil
}

func (pw *plotWindowCommon) addPoints(r *rpn.RPN, plot Plot, steps uint32, fn func(x, y float64, coloridx uint8) error) error {
	startlen := r.StackLen()
	step := (pw.maxv - pw.minv) / float64(steps)
	var x complex128
	for v := pw.minv; v <= pw.maxv; v += step {
		if err := r.PushComplex(complex(v, 0)); err != nil {
			return err
		}
		if err := r.Exec(plot.fn); err != nil {
			return err
		}
		y, err := r.PopComplex()
		if err != nil {
			return err
		}
		if plot.isParametric {
			x, err = r.PopComplex()
			if err != nil {
				return err
			}
		} else {
			x = complex(v, 0)
		}
		nowlen := r.StackLen()
		if nowlen != startlen {
			return fmt.Errorf(
				"stack changed size running plot string (old: %d, new %d)",
				startlen,
				nowlen)
		}
		rx := real(x)
		ry := real(y)
		if err := fn(rx, ry, plot.coloridx); err != nil {
			return err
		}
	}
	return nil
}

func (pw *plotWindowCommon) adjustAutoX() {
	pw.minx = pw.stats.minx
	pw.maxx = pw.stats.maxx
	if pw.minx == pw.maxx {
		// create a little spread to avoid math issues
		pw.minx -= 1.0
		pw.maxx += 1.0
	}
}

func (pw *plotWindowCommon) adjustAutoY() {
	pw.miny = pw.stats.miny
	pw.maxy = pw.stats.maxy
	if pw.miny == pw.maxy {
		// create a little spread to avoid math issues
		pw.miny -= 1.0
		pw.maxy += 1.0
	}
	// open up the y a bit
	delta := (pw.maxy - pw.miny) / 5
	pw.maxy += delta
	pw.miny -= delta
}

func (pw *plotWindowCommon) transformX(x float64, w int) (int, bool) {
	x = (x - pw.minx) / (pw.maxx - pw.minx)
	if x < 0 {
		return 0, false
	}
	xi := int(float64(w)*x + 0.5)
	if xi < 0 || xi >= w {
		return 0, false
	}
	return xi, true
}

func (pw *plotWindowCommon) transformY(y float64, h int) (int, bool) {
	y = (y - pw.miny) / (pw.maxy - pw.miny)
	if y < 0 {
		return 0, false
	}
	py := h - int(float64(h)*y+0.5) - 1
	if py < 0 || py > h {
		return 0, false
	}
	return py, true
}

// use a nice-looking scale. 1, 0.5, 0.25, 0.1, 0.05, 0.025, 0.01, etc
func searchScaleDownward(cpu, minSpacing float64) float64 {
	tens := 1.0
	partial := 1
	te := 1.0

	for {
		switch partial {
		case 1:
			partial = 2
		case 2:
			partial = 4
		case 4:
			partial = 1
			tens *= 10
		}

		newte := 1.0 / (tens * float64(partial))
		if (cpu * newte) < minSpacing {
			// too far
			break
		}
		te = newte
	}

	return te
}

func searchScaleUpward(cpu, maxSpacing float64) float64 {
	tens := 1.0
	partialDeci := 10
	te := 1.0

	for {
		switch partialDeci {
		case 10:
			partialDeci = 25
		case 25:
			partialDeci = 50
		case 50:
			partialDeci = 10
			tens *= 10
		}

		newte := tens * float64(partialDeci) / 10.0
		if (cpu * newte) > maxSpacing {
			// too far
			break
		}
		te = newte
	}

	return te
}
