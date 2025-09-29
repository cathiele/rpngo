package window

import (
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

func (wr *WindowRoot) FindWindow(name string) WindowWithProps {
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
		return nil, rpn.ErrNotFound
	}
	if wge.group == nil {
		return nil, rpn.ErrNotAWindowGroup
	}
	return wge.group, nil
}

func (wr *WindowRoot) DeleteWindowOrGroup(name string) error {
	if name == "root" {
		return rpn.ErrCanNotDeleteRootWindow
	}
	if name == "i" {
		return rpn.ErrCanNotDeleteInputWindow
	}
	pwge, wge := wr.group.findwindowGroupEntryAndParent(name)
	if wge == nil {
		return rpn.ErrNotFound
	}
	pwge.removeChild(wge)
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) MoveWindowOrGroup(src string, dst string, beginning bool) error {
	if src == "root" {
		return rpn.ErrIllegalWindowOperation
	}
	srcpg, srcwge := wr.group.findwindowGroupEntryAndParent(src)
	if srcwge == nil {
		return rpn.ErrNotFound
	}
	if srcwge.group != nil {
		check := srcwge.group.findwindowGroupEntry(dst)
		if check != nil {
			return rpn.ErrIllegalWindowOperation
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
		return rpn.ErrIllegalWindowOperation
	}
	if w > 10000 {
		return rpn.ErrIllegalWindowOperation
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return rpn.ErrNotFound
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

func (wr *WindowRoot) AddWindowChild(window WindowWithProps, name string, weight int) {
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
func (wr *WindowRoot) Update(r *rpn.RPN, screenw, screenh int, updateInput bool) error {
	if wr.adjustNeeded {
		if err := wr.group.adjustChildren(screenw, screenh); err != nil {
			return err
		}
		wr.group.showBorder(screenw, screenh)
		wr.adjustNeeded = false
		// We want to give screens other than the input screen a chance to
		// redraw before getting locked into the input screen
		if err := wr.group.update(r, screenw, screenh); err != nil {
			return err
		}
	}
	if updateInput {
		// Update the input window first
		input := wr.FindWindow("i")
		if input == nil {
			return rpn.ErrNotFound
		}
		if err := input.Update(r); err != nil {
			return err
		}
	}
	return wr.group.update(r, screenw, screenh)
}

func (wr *WindowRoot) UpdateByName(r *rpn.RPN, name string) error {
	if name == "root" {
		if err := wr.Update(r, wr.group.w, wr.group.y, false); err != nil {
			return err
		}
		return wr.group.showBorder(wr.group.w, wr.group.h)
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return rpn.ErrNotFound
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
