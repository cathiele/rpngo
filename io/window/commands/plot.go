package commands

import (
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
	return wc.plotInternal(r, false)
}

const PPlotHelp = "Run a value from $plot.min to $plot.max with $plot.steps\n" +
	"executing the provided string and plotting (x, y) to the window $plot.win\n" +
	"Example (arc): 'c sin sw cos"

func (wc *WindowCommands) PPlot(r *rpn.RPN) error {
	return wc.plotInternal(r, true)
}

func (wc *WindowCommands) plotInternal(r *rpn.RPN, isParametric bool) error {
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
		if err := wc.initPlot(r); err != nil {
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

	return pw.(*plotwin.PlotWindow).AddPlot(r, fields, isParametric)
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
