package commands

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/io/window/plotwin"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const PlotHelp = "Run a value from $plot.min to $plot.max with $plot.steps\n" +
	"executing the provided string and plotting to the window $plot.win\n" +
	"Example (y = x * x): 'c *' plot"

func (wc *WindowCommands) Plot(r *rpn.RPN) error {
	macro, err := r.PopString()
	if err != nil {
		return err
	}
	fields, err := parse.Fields(macro)
	if err != nil {
		return err
	}
	wname, err := r.GetStringVariable("plot.win")
	if err != nil {
		return err
	}
	pw := wc.root.FindWindow(wname)
	if pw == nil {
		if err := wc.initPlot(r); wc != nil {
			return err
		}
		pw = wc.root.FindWindow(wname)
	}
	if pw == nil {
		return fmt.Errorf("executing $plot.init did not result in a window: %s", wname)
	}

	if pw.Type() != "plot" {
		return fmt.Errorf("%s has the wrong window type: %s", wname, pw.Type())
	}

	xmin, err := r.GetComplexVariable("plot.min")
	if err != nil {
		return err
	}

	xmax, err := r.GetComplexVariable("plot.max")
	if err != nil {
		return err
	}

	steps, err := r.GetComplexVariable("plot.steps")
	if err != nil {
		return err
	}
	rsteps := real(steps)
	if rsteps < 2 {
		return errors.New("plot.steps must be 2 or greater")
	}
	if rsteps > 50000 {
		return errors.New("plot.steps must be 50000 or less")
	}

	return makePlot(r, pw.(*plotwin.PlotWindow), fields, real(xmin), real(xmax), int(rsteps))
}

func (wc *WindowCommands) initPlot(r *rpn.RPN) error {
	macro, err := r.GetStringVariable("plot.init")
	if err != nil {
		return err
	}
	fields, err := parse.Fields(macro)
	if err != nil {
		return err
	}
	err = r.Exec(fields)
	w, h := wc.screen.Size()
	if uerr := wc.root.Update(r, w, h, false); uerr != nil {
		log.Printf("initPlot.Update error: %v", uerr)
	}
	if err != nil {
		return fmt.Errorf("while executing $plot.init: %v", err)
	}
	return nil
}

func makePlot(r *rpn.RPN, pw *plotwin.PlotWindow, fields []string, xmin, xmax float64, steps int) error {
	if xmax <= xmin {
		return errors.New("plot.max must be > plot.min")
	}
	startlen := r.StackLen()
	step := (xmax - xmin) / float64(steps)
	for x := xmin; x <= xmax; x += step {
		if err := r.PushComplex(complex(x, 0)); err != nil {
			return err
		}
		if err := r.Exec(fields); err != nil {
			return err
		}
		y, err := r.PopComplex()
		if err != nil {
			return err
		}
		nowlen := r.StackLen()
		if nowlen != startlen {
			return fmt.Errorf("Stack changed size running plot string (old: %d, new %d)", startlen, nowlen)
		}
		pw.SetPoint(x, real(y))
	}
	return nil
}
