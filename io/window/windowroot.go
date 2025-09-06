package window

import (
	"errors"
	"fmt"
	"mattwach/rpngo/rpn"
	"strings"
)

type WindowRoot struct {
	group        *windowGroup
	adjustNeeded bool
}

func NewWindowRoot(w, h int) *WindowRoot {
	wr := &WindowRoot{
		group: &windowGroup{},
	}
	wr.group.resize(0, 0, w, h)
	return wr
}

func (wr *WindowRoot) FindWindow(name string) Window {
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return nil
	}
	return wge.window
}

func (wr *WindowRoot) FindwindowGroup(name string) (*windowGroup, error) {
	if name == "root" {
		return wr.group, nil
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return nil, fmt.Errorf("window group not found: %s", name)
	}
	if wge.group == nil {
		return nil, fmt.Errorf("not a window group: %s", name)
	}
	return wge.group, nil
}

func (wr *WindowRoot) DeleteWindowOrGroup(name string) error {
	if name == "root" {
		return errors.New("can not delete root window")
	}
	if name == "i" {
		return errors.New("can not delete input window")
	}
	pwge, wge := wr.group.findwindowGroupEntryAndParent(name)
	if wge == nil {
		return fmt.Errorf("window not found: %s", name)
	}
	pwge.removeChild(wge)
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) MoveWindowOrGroup(src string, dst string, beginning bool) error {
	if src == "root" {
		return errors.New("can not move root window")
	}
	srcpg, srcwge := wr.group.findwindowGroupEntryAndParent(src)
	if srcwge == nil {
		return fmt.Errorf("source window not found: %s", src)
	}
	if srcwge.group != nil {
		check := srcwge.group.findwindowGroupEntry(dst)
		if check != nil {
			return fmt.Errorf("moving %s to %s would detach from root", src, dst)
		}
	}
	dstpg, err := wr.FindwindowGroup(dst)
	if err != nil {
		return err
	}
	srcpg.removeChild(srcwge)
	if beginning {
		dstpg.children = append([]*windowGroupEntry{srcwge}, dstpg.children...)
	} else {
		dstpg.children = append(dstpg.children, srcwge)
	}
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) SetWindowWeight(name string, w int) error {
	if w < 1 {
		return fmt.Errorf("weight must be >= 1: %d", w)
	}
	if w > 10000 {
		return fmt.Errorf("weight must be <= 10000: %d", w)
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return fmt.Errorf("window not found: %s", name)
	}
	wge.weight = w
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) AddNewWindowGroupChild(name string, weight int) {
	group := &windowGroup{}
	wr.group.children = append(wr.group.children, &windowGroupEntry{name: name, weight: weight, group: group})
	wr.adjustNeeded = true
}

func (wr *WindowRoot) AddWindowChild(window Window, name string, weight int) {
	wr.group.children = append(wr.group.children, &windowGroupEntry{name: name, weight: weight, window: window})
	wr.adjustNeeded = true
}

func (wr *WindowRoot) UseColumnLayout(name string, v bool) error {
	wg, err := wr.FindwindowGroup(name)
	if err != nil {
		return err
	}
	wg.isColumn = v
	wr.adjustNeeded = true
	return nil
}

// Calls update on all contained windows
func (wr *WindowRoot) Update(rpn *rpn.RPN, screenw, screenh int, updateInput bool) error {
	if wr.adjustNeeded {
		if err := wr.group.adjustChildren(screenw, screenh); err != nil {
			return err
		}
		wr.group.showBorder(screenw, screenh)
		wr.adjustNeeded = false
		// We want to give screens other than the input screen a chance to
		// redraw before getting locked into the input screen
		if err := wr.group.update(rpn, screenw, screenh); err != nil {
			return err
		}
	}
	if updateInput {
		// Update the input window first
		input := wr.FindWindow("i")
		if input == nil {
			return errors.New("could not find window 'i' for input")
		}
		if err := input.Update(rpn); err != nil {
			return err
		}
	}
	return wr.group.update(rpn, screenw, screenh)
}

func (wr *WindowRoot) UpdateByName(r *rpn.RPN, name string) error {
	if name == "root" {
		return wr.Update(r, wr.group.w, wr.group.y, false)
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return fmt.Errorf("window not found: %s", name)
	}
	if wge.group != nil {
		return wge.group.update(r, wr.group.w, wr.group.h)
	}
	return wge.window.Update(r)
}

// Removes all children
func (wr *WindowRoot) RemoveAllChildren() {
	wr.group.removeAllChildren()
}

func (wr *WindowRoot) Dump(r *rpn.RPN) {
	lines := wr.group.dump(nil, "root", 0, 100)
	r.Println(strings.Join(lines, "\n"))
}
