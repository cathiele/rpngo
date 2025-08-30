// Package commands is creates window management commands
package commands

import (
	"errors"
	"fmt"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/io/window/stackwin"
	"mattwach/rpngo/io/window/varwin"
	"mattwach/rpngo/rpn"
)

type WindowCommands struct {
	root   *window.WindowRoot
	screen window.Screen
}

func InitWindowCommands(root *window.WindowRoot, screen window.Screen) *WindowCommands {
	return &WindowCommands{root: root, screen: screen}
}

func (wc *WindowCommands) Register(r *rpn.RPN) {
	r.Register("w.columns", wc.WColumns, WColumnsHelp)
	r.Register("w.dump", wc.WDump, WDumpHelp)
	r.Register("w.move.beg", wc.WMoveBeg, WMoveBegHelp)
	r.Register("w.move.end", wc.WMoveEnd, WMoveEndHelp)
	r.Register("w.new.group", wc.WNewGroup, WNewGroupHelp)
	r.Register("w.new.stack", wc.WNewStack, WNewStackHelp)
	r.Register("w.new.var", wc.WNewVar, WNewVarHelp)
	r.Register("w.reset", wc.WReset, WResetHelp)
	r.Register("w.weight", wc.WWeight, WWeightHelp)
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
		return errors.New("internal error: no window 'i' was found")
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
		return "", fmt.Errorf("window already exits: %s", name)
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
