// Package commands is creates window management commands
package commands

import (
	"errors"
	"fmt"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/io/window/stackwin"
	"mattwach/rpngo/rpn"
)

type WindowCommands struct {
	root   *window.WindowGroup
	screen window.Screen
}

func InitWindowCommands(root *window.WindowGroup, screen window.Screen) *WindowCommands {
	return &WindowCommands{root: root, screen: screen}
}

func (wc *WindowCommands) Register(r *rpn.RPN) {
	r.Register("wnewstack", wc.WNewStack, WNewStackHelp)
	r.Register("wreset", wc.WReset, WResetHelp)
	r.Register("wweight", wc.WWeight, WWeightHelp)
}

const WResetHelp = "Resets window configuration to just a single input window"

func (wc *WindowCommands) WReset(r *rpn.Stack) error {
	iw := wc.root.FindWindow("i")
	if iw == nil {
		return errors.New("internal error: no window 'i' was found")
	}
	wc.root.RemoveAllChildren()
	wc.root.AddWindowChild(iw, "i", 100)
	return nil
}

const WNewStackHelp = "Creates a new stack window with the given name and\n" +
	"adds it to the root window. Example: 's1' wnewstack"

func (wc *WindowCommands) WNewStack(r *rpn.Stack) error {
	name, err := r.PopString()
	if err != nil {
		return err
	}
	existing := wc.root.FindWindow(name)
	if existing != nil {
		return fmt.Errorf("window already exits: %s", name)
	}
	txtw, err := wc.screen.NewTextWindow(0, 0, 10, 5)
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

const WWeightHelp = "Changes the weight of a window or window group causing it\n" +
	"to take more or less screen space. The default value is 100.\n" +
	"Example: 's1' 20 wweight"

func (wc *WindowCommands) WWeight(r *rpn.Stack) error {
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
