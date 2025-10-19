package commands

import (
	"mattwach/rpngo/rpn"
)

const WListPHelp = "Prints all properties / values for a window\n" +
	"Example 'p1' w.listp"

func (wc *WindowCommands) WListP(r *rpn.RPN) error {
	wname, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !wname.IsString() {
		return rpn.ErrExpectedAString
	}
	w := wc.root.FindWindow(wname.UnsafeString())
	if w == nil {
		return rpn.ErrNotFound
	}
	for _, p := range w.ListProps() {
		f, err := w.GetProp(p)
		if err != nil {
			// unexpected
			return err
		}
		r.Print(p)
		r.Print(": ")
		r.Println(f.String(true))
	}
	return nil
}

const WGetPHelp = "Pushes the value of the given property to the stack.\n" +
	"Example: 'p1' 'minx' w.getp"

func (wc *WindowCommands) WGetP(r *rpn.RPN) error {
	wname, pname, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !wname.IsString() || !pname.IsString() {
		return rpn.ErrExpectedAString
	}
	w := wc.root.FindWindow(wname.UnsafeString())
	if w == nil {
		return rpn.ErrNotFound
	}
	f, err := w.GetProp(pname.UnsafeString())
	if err != nil {
		return err
	}
	return r.PushFrame(f)
}

const WSetPHelp = "Sets a property on a window.\n" +
	"Example: 'p1' 'minx' -1 w.setp"

func (wc *WindowCommands) WSetP(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	wname, pname, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !wname.IsString() || !pname.IsString() {
		return rpn.ErrExpectedAString
	}
	w := wc.root.FindWindow(wname.UnsafeString())
	if w == nil {
		return rpn.ErrNotFound
	}
	return w.SetProp(pname.UnsafeString(), f)
}
