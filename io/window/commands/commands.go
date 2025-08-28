// Package commands is creates window management commands
package commands

import (
	"errors"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/rpn"
)

type WindowCommands struct {
	root *window.WindowGroup
}

func InitWindowCommands(root *window.WindowGroup) *WindowCommands {
	return &WindowCommands{root: root}
}

func (wc *WindowCommands) Register(r *rpn.RPN) {
	r.Register("wreset", wc.WReset, WResetHelp)
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
