// Package commands is creates window management commands
package commands

import (
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
	"mattwach/rpngo/window/plotwin"
	"mattwach/rpngo/window/stackwin"
	"mattwach/rpngo/window/varwin"
)

type WindowCommands struct {
	root   *window.WindowRoot
	screen window.Screen
}

func InitWindowCommands(r *rpn.RPN, root *window.WindowRoot, screen window.Screen) *WindowCommands {
	conceptHelp := map[string]string{
		"window.layout": "Windows are arranged with window groups.  There\n" +
			"is always a window group named 'root' which is the parent of all \n" +
			"windows and groups.\n" +
			"- Add a new window group to the root window with w.new.group.\n" +
			"- Move a window or group to a different window group with w.move.beg and w.move.end\n" +
			"- Change the weight of a window or group with w.weight (default weight is 100).\n" +
			"- Change the layout mode of a window group to columns with w.columns.\n" +
			"- Print info on all existing windows and groups with w.dump.\n" +
			"See Also: windows, window.props",

		"window.props": "Each window supports properties that changes how the window operates\n" +
			"- Print all properties and values for a window with w.listp\n" +
			"- Get a single property with w.getp\n" +
			"- Set a single property with w.setp\n" +
			"See Also: windows, window.layout, plotting",

		"windows": "The display can be customized with different windows\n" +
			"- Add a window with a w.new.<type> command. Example: w.new.stack\n" +
			"- Reset to a single window with w.reset.\n" +
			"See Also: window.layout, window.props",
	}
	r.RegisterConceptHelp(conceptHelp)
	wc := WindowCommands{root: root, screen: screen}
	r.Register("w.columns", wc.WColumns, rpn.CatWindow, WColumnsHelp)
	r.Register("w.del", wc.WDelete, rpn.CatWindow, WDeleteHelp)
	r.Register("w.dump", wc.WDump, rpn.CatWindow, WDumpHelp)
	r.Register("w.move.beg", wc.WMoveBeg, rpn.CatWindow, WMoveBegHelp)
	r.Register("w.move.end", wc.WMoveEnd, rpn.CatWindow, WMoveEndHelp)
	r.Register("w.new.group", wc.WNewGroup, rpn.CatWindow, WNewGroupHelp)
	r.Register("w.new.plot", wc.WNewPlot, rpn.CatWindow, WNewPlotHelp)
	r.Register("w.new.stack", wc.WNewStack, rpn.CatWindow, WNewStackHelp)
	r.Register("w.new.var", wc.WNewVar, rpn.CatWindow, WNewVarHelp)
	r.Register("w.listp", wc.WListP, rpn.CatWindow, WListPHelp)
	r.Register("w.getp", wc.WGetP, rpn.CatWindow, WGetPHelp)
	r.Register("w.setp", wc.WSetP, rpn.CatWindow, WSetPHelp)
	r.Register("w.reset", wc.WReset, rpn.CatWindow, WResetHelp)
	r.Register("w.update", wc.WUpdate, rpn.CatWindow, WUpdateHelp)
	r.Register("w.weight", wc.WWeight, rpn.CatWindow, WWeightHelp)
	return &wc
}

const WUpdateHelp = "Updates the given window or window group"

func (wc *WindowCommands) WUpdate(r *rpn.RPN) error {
	name, err := r.PopString()
	if err != nil {
		return err
	}
	return wc.root.UpdateByName(r, name)
}

const WDumpHelp = "Dump the state of all created windows and groups"

func (wc *WindowCommands) WDump(r *rpn.RPN) error {
	wc.root.Dump(r)
	return nil
}

const WResetHelp = "Resets window configuration to just a single input window"

func (wc *WindowCommands) WReset(r *rpn.RPN) error {
	iw := wc.root.FindWindow("i")
	if iw == nil {
		return rpn.ErrInputWindowNotFound
	}
	wc.root.RemoveAllChildren()
	wc.root.UseColumnLayout("root", false)
	wc.root.AddWindowChild(iw, "i", 100)
	return nil
}

const WColumnsHelp = "Sets a window group layout to column mode\n" +
	"Example: 'g1' w.columns"

func (wc *WindowCommands) WColumns(r *rpn.RPN) error {
	name, err := r.PopString()
	if err != nil {
		return err
	}
	if err := wc.root.UseColumnLayout(name, true); err != nil {
		return err
	}
	return nil
}

const WDeleteHelp = "Deletes a window or window group\n" +
	"Example: 'p1' w.del"

func (wc *WindowCommands) WDelete(r *rpn.RPN) error {
	name, err := r.PopString()
	if err != nil {
		return err
	}
	return wc.root.DeleteWindowOrGroup(name)
}

const WMoveBegHelp = "Moves a window or group to the beginning of a window group\n" +
	"Example: 's1' 'root' w.move.beg"

func (wc *WindowCommands) WMoveBeg(r *rpn.RPN) error {
	src, dst, err := r.Pop2Strings()
	if err != nil {
		return err
	}
	return wc.root.MoveWindowOrGroup(src, dst, true)
}

const WMoveEndHelp = "Moves a window or group to the end of a window group\n" +
	"Example: 's1' 'root' w.move.end"

func (wc *WindowCommands) WMoveEnd(r *rpn.RPN) error {
	src, dst, err := r.Pop2Strings()
	if err != nil {
		return err
	}
	return wc.root.MoveWindowOrGroup(src, dst, false)
}

const WNewGroupHelp = "Creates a new window group with the given name and\n" +
	"adds it to the root window. Example: 'g1' w.new.group"

func (wc *WindowCommands) WNewGroup(r *rpn.RPN) error {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return err
	}
	wc.root.AddNewWindowGroupChild(name, 100)
	return nil
}

const WNewStackHelp = "Creates a new stack window with the given name and\n" +
	"adds it to the root window. Example: 's1' w.new.stack"

func (wc *WindowCommands) WNewStack(r *rpn.RPN) error {
	txtw, name, err := wc.newTextWindow(r)
	if err != nil {
		return err
	}
	sw, err := stackwin.Init(txtw)
	if err != nil {
		return err
	}
	wc.root.AddWindowChild(sw, name, 100)
	return nil
}

const WNewPlotHelp = "Creates a new plot window with the given name and\n" +
	"adds it to the root window. Example: 'p1' w.new.plot"

func (wc *WindowCommands) WNewPlot(r *rpn.RPN) error {
	txtw, name, err := wc.newTextWindow(r)
	if err != nil {
		return err
	}
	var pw plotwin.TxtPlotWindow
	pw.Init(txtw)
	wc.root.AddWindowChild(&pw, name, 100)
	return nil
}

const WNewVarHelp = "Creates a new variable window with the given name and\n" +
	"adds it to the root window. Example: 'v1' w.new.var"

func (wc *WindowCommands) WNewVar(r *rpn.RPN) error {
	txtw, name, err := wc.newTextWindow(r)
	if err != nil {
		return err
	}
	sw, err := varwin.Init(txtw)
	if err != nil {
		return err
	}
	wc.root.AddWindowChild(sw, name, 100)
	return nil
}

func (wc *WindowCommands) newTextWindow(r *rpn.RPN) (window.TextWindow, string, error) {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return nil, "", err
	}
	txtw, err := wc.screen.NewTextWindow(0, 0, 10, 5)
	return txtw, name, err
}

func (wc *WindowCommands) newWindowNameFromStack(r *rpn.RPN) (string, error) {
	name, err := r.PopString()
	if err != nil {
		return "", err
	}
	existing := wc.root.FindWindow(name)
	if existing != nil {
		return "", rpn.ErrWindowAlreadyExists
	}
	return name, nil
}

const WWeightHelp = "Changes the weight of a window or window group causing it\n" +
	"to take more or less screen space. The default value is 100.\n" +
	"Example: 's1' 20 w.weight"

func (wc *WindowCommands) WWeight(r *rpn.RPN) error {
	cw, err := r.PopComplex()
	if err != nil {
		return err
	}
	w := int(real(cw))
	name, err := r.PopString()
	if err != nil {
		r.PushComplex(cw)
		return err
	}
	if err := wc.root.SetWindowWeight(name, w); err != nil {
		r.PushString(name)
		r.PushComplex(cw)
		return err
	}
	return nil
}
