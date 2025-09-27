package plotwin

import (
	"fmt"
	"log"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

type PlotCommands struct {
	root   *window.WindowRoot
	screen window.Screen
}

func InitPlotCommands(r *rpn.RPN, root *window.WindowRoot, screen window.Screen) *PlotCommands {
	conceptHelp := map[string]string{
		"plot": "Plot functions using plot. Plot will push an 'x' value to the stack,\n" +
			"run the provided string, and pop the value as y value.\n" +
			"Examples:\n" +
			"    '2 *' plot # plots y = x * 2\n" +
			"    'sq' plot # plots y = x * x\n" +
			"    'sin' plot # plots y = sin(x)\n" +
			"Various properties can be set on the plot window to change the number\n" +
			"of points and the boundaries of the plot.\n" +
			"There are some special variables that plot uses:\n" +
			"    plot.win  : Name of the window to send plots to (there can be more than one)\n" +
			"                at a time.\n" +
			"    plot.init : If no plot window exists, this string is executed and is expected\n" +
			"                create one. Making this a variable allows for user customization.\n" +
			"See Also: window.props, plot.parametric",

		"plot.parametric": "Plot parametric functions using pplot. pplot will push a 't' value to\n" +
			"the stack, run the provided string then pop y, then x to determine the plot point x, y\n" +
			"Examples:\n" +
			"    '$0 cos 1> sin' pplot # draws an arc or full circle, depending on t range\n" +
			"    't= $t sin $t * $t cos $t *' pplot # draw a spiral\n" +
			"    '1 sw' draw a vertical line\n",
	}
	r.RegisterConceptHelp(conceptHelp)

	pc := PlotCommands{root: root, screen: screen}
	r.Register("plot", pc.Plot, rpn.CatPlot, PlotHelp)
	r.Register("pplot", pc.PPlot, rpn.CatPlot, PPlotHelp)
	return &pc
}

const PlotHelp = "Run a value from $plot.min to $plot.max with $plot.steps\n" +
	"executing the provided string and plotting to the window $plot.win\n" +
	"Example (y = x * x): 'c *' plot"

func (pc *PlotCommands) Plot(r *rpn.RPN) error {
	return pc.plotInternal(r, false)
}

const PPlotHelp = "Run a value from $plot.min to $plot.max with $plot.steps\n" +
	"executing the provided string and plotting (x, y) to the window $plot.win\n" +
	"Example (arc): 'c sin sw cos"

func (wc *PlotCommands) PPlot(r *rpn.RPN) error {
	return wc.plotInternal(r, true)
}

func (pc *PlotCommands) plotInternal(r *rpn.RPN, isParametric bool) error {
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
	pw := pc.root.FindWindow(wname)
	if pw == nil {
		if err := pc.initPlot(r); err != nil {
			return err
		}
		pw = pc.root.FindWindow(wname)
	}
	if pw == nil {
		return fmt.Errorf("executing $plot.init did not result in a window: %s", wname)
	}

	if pw.Type() != "plot" {
		return fmt.Errorf("%s has the wrong window type: %s", wname, pw.Type())
	}

	return pw.(*PlotWindow).AddPlot(r, fields, isParametric)
}

func (wc *PlotCommands) initPlot(r *rpn.RPN) error {
	macro, err := r.GetStringVariable("plot.init")
	if err != nil {
		return err
	}
	fields, err := parse.Fields(macro)
	if err != nil {
		return err
	}
	err = r.Exec(fields)
	w, h := wc.screen.ScreenSize()
	if uerr := wc.root.Update(r, w, h, false); uerr != nil {
		log.Printf("initPlot.Update error: %v", uerr)
	}
	if err != nil {
		return fmt.Errorf("while executing $plot.init: %v", err)
	}
	return nil
}
