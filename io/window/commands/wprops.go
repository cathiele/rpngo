package commands

import (
	"fmt"
	"mattwach/rpngo/rpn"
)

const WListPHelp = "Prints all properties / values for a window\n" +
	"Example 'p1' w.listp"

func (wc *WindowCommands) WListP(r *rpn.RPN) error {
	wname, err := r.PopString()
	if err != nil {
		return err
	}
	w := wc.root.FindWindow(wname)
	if w == nil {
		return fmt.Errorf("window not found: %s", wname)
	}
	for _, p := range w.ListProps() {
		f, err := w.GetProp(p)
		if err != nil {
			// unexpected
			return err
		}
		r.Print(fmt.Sprintf("%s: %s\n", p, f.String(true)))
	}
	return nil
}

const WGetPHelp = "Pushes the value of the given property to the stack.\n" +
	"Example: 'p1' 'minx' w.getp"

func (wc *WindowCommands) WGetP(r *rpn.RPN) error {
	wname, pname, err := r.Pop2Strings()
	if err != nil {
		return err
	}
	w := wc.root.FindWindow(wname)
	if w == nil {
		return fmt.Errorf("window not found: %s", wname)
	}
	f, err := w.GetProp(pname)
	if err != nil {
		return err
	}
	return r.PushFrame(f)
}

const WSetPHelp = "Sets a property on a window.\n" +
	"Example: 'p1' 'minx' -1 w.setp"

func (wc *WindowCommands) WSetP(r *rpn.RPN) error {
	wname, pname, err := r.Pop2Strings()
	if err != nil {
		return err
	}
	w := wc.root.FindWindow(wname)
	if w == nil {
		return fmt.Errorf("window not found: %s", wname)
	}
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	return w.SetProp(pname, f)
}
