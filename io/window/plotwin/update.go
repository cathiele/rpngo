package plotwin

import (
	"fmt"
	"mattwach/rpngo/rpn"
)

func (pw *PlotWindow) Update(r *rpn.RPN) error {
	points, err := pw.createPoints(r)
	if err != nil {
		return err
	}
	pw.txtw.Erase()
	defer pw.txtw.Refresh()
	if pw.autox {
		pw.adjustAutoX(points)
	}
	if pw.autoy {
		pw.adjustAutoY(points)
	}
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
	// open up the y a bit (20% or so)
	delta := (pw.maxy - pw.miny) / 10
	pw.maxy += delta
	pw.miny -= delta
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
		pw.txtw.SetCursorXY(x, y)
		if err := pw.txtw.Write('*'); err != nil {
			return err
		}
	}
	return nil
}

func (pw *PlotWindow) transformX(x float64) (int, bool) {
	x = (x - pw.minx) / (pw.maxx - pw.minx)
	if x < 0 {
		return 0, false
	}
	xi := int(float64(pw.txtw.TextWidth())*x + 0.5)
	if xi < 0 || xi >= pw.txtw.TextWidth() {
		return 0, false
	}
	return xi, true
}

func (pw *PlotWindow) transformY(y float64) (int, bool) {
	y = (y - pw.miny) / (pw.maxy - pw.miny)
	if y < 0 {
		return 0, false
	}
	py := pw.txtw.TextHeight() - int(float64(pw.txtw.TextHeight())*y+0.5) - 1
	if py < 0 || py > pw.txtw.TextHeight() {
		return 0, false
	}
	return py, true
}
