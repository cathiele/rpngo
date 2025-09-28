package plotwin

import (
	"fmt"
	"mattwach/rpngo/window"
)

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
		pw.txtw.SetCursorXY(x, y)
		if err := pw.txtw.Write('+'); err != nil {
			return err
		}
	}

	return nil
}

func (pw *PlotWindow) drawVerticalAxis(x int) error {
	h := pw.txtw.TextHeight()
	for y := 0; y < h; y++ {
		pw.txtw.SetCursorXY(x, y)
		if err := pw.txtw.Write('|'); err != nil {
			return err
		}
	}
	pw.drawVerticalTickMarks(x)
	return nil
}

// About axis drawing
// the goal is to make it look nice while also being useful
// We are going for rounded tick marks (e.g. 0.25, 1.0, ...)
// and reasonable pixels per mark

const minVerticalSpacing = 5.0
const maxVerticalSpacing = 10.0

func (pw *PlotWindow) drawVerticalTickMarks(wx int) {
	wh := pw.txtw.TextHeight() // height in characters
	yr := pw.maxy - pw.miny    // units
	cpu := float64(wh) / yr    // characters / unit
	var te float64 = 1         // ticks every (0.5, 1, etc)
	if cpu > maxVerticalSpacing {
		te = searchScaleDownward(cpu, minVerticalSpacing)
	} else if cpu < minVerticalSpacing {
		te = searchScaleUpward(cpu, maxVerticalSpacing)
	}

	// becuase te was carefully selected to provide a limited number
	// of tick marks, we can use a simple loop to determine the min te
	stepsBack := 0
	for (float64(stepsBack-1) * te) > pw.miny {
		stepsBack--
	}

	for {
		y := float64(stepsBack) * te
		if y >= pw.maxy {
			break
		}
		if stepsBack != 0 {
			pw.drawVerticalTick(wx, y)
		}
		stepsBack++
	}
}

func (pw *PlotWindow) drawVerticalTick(wx int, y float64) {
	ww := pw.txtw.TextWidth()
	wy, _ := pw.transformY(y)
	pw.txtw.SetCursorXY(wx, wy)
	window.PutByte(pw.txtw, '+')
	if (wx + 10) < ww {
		pw.txtw.SetCursorXY(wx+3, wy)
		window.Print(pw.txtw, fmt.Sprintf("%.2f", y))
	}
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

func (pw *PlotWindow) drawHorizontalAxis(y int) error {
	w := pw.txtw.TextWidth()
	pw.txtw.SetCursorXY(0, y)
	for x := 0; x < w; x++ {
		if err := pw.txtw.Write('-'); err != nil {
			return err
		}
	}
	pw.drawHorizontalTickMarks(y)
	return nil
}

const minHorizontalSpacing = 9.0
const maxHorizontalSpacing = 18.0

func (pw *PlotWindow) drawHorizontalTickMarks(wy int) {
	ww := pw.txtw.TextWidth() // width in characters
	xr := pw.maxx - pw.minx   // units
	cpu := float64(ww) / xr   // characters / unit
	var te float64 = 1        // ticks every (0.5, 1, etc)
	if cpu > maxHorizontalSpacing {
		te = searchScaleDownward(cpu, minHorizontalSpacing)
	} else if cpu < minHorizontalSpacing {
		te = searchScaleUpward(cpu, maxHorizontalSpacing)
	}

	stepsBack := 0
	for (float64(stepsBack-1) * te) > pw.minx {
		stepsBack--
	}

	for {
		x := float64(stepsBack) * te
		if x >= pw.maxx {
			break
		}
		if stepsBack != 0 {
			pw.drawHorizontalTick(x, wy)
		}
		stepsBack++
	}
}

const horizontalNumberPad = 5

func (pw *PlotWindow) drawHorizontalTick(x float64, wy int) {
	ww, wh := pw.txtw.TextSize()
	wx, _ := pw.transformX(x)
	if wy >= 0 {
		pw.txtw.SetCursorXY(wx, wy)
		window.PutByte(pw.txtw, '+')
	}
	if (wx > horizontalNumberPad) && (wx < (ww - horizontalNumberPad)) && (wy+1) < wh {
		s := fmt.Sprintf("%.2f", x)
		pw.txtw.SetCursorXY(wx-len(s)/2, wy+1)
		window.Print(pw.txtw, s)
	}
}
