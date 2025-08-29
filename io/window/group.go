package window

import (
	"errors"
	"fmt"
	"mattwach/rpngo/rpn"
)

type Window interface {
	Update(*rpn.RPN) error
	Resize(x, y, w, h int) error
	ShowBorder(screenw, screenh int) error
}

type windowGroupEntry struct {
	name   string
	weight int
	// Only one of the following should be non-nil
	group  *WindowGroup
	window Window
}

func (wge *windowGroupEntry) resize(x, y, w, h int) {
	if wge.group != nil {
		wge.group.Resize(x, y, w, h)
	} else if wge.window != nil {
		wge.window.Resize(x, y, w, h)
	}
}

func (wge *windowGroupEntry) showBorder(screenw, screenh int) error {
	if wge.group != nil {
		for _, c := range wge.group.children {
			if err := c.showBorder(screenw, screenh); err != nil {
				return err
			}
		}
	} else if wge.window != nil {
		if err := wge.window.ShowBorder(screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}

type WindowGroup struct {
	isRoot       bool
	isColumn     bool
	adjustNeeded bool
	// Coordinates are in global screen coordinates
	x        int
	y        int
	w        int
	h        int
	children []*windowGroupEntry
}

func NewWindowGroup(isRoot bool) *WindowGroup {
	return &WindowGroup{isRoot: isRoot}
}

func (wg *WindowGroup) removeChild(c *windowGroupEntry) {
	idx := -1
	for i, rc := range wg.children {
		if rc == c {
			idx = i
			break
		}
	}
	// idx may be -1 here, which will crash, but that is what we want
	wg.children = append(wg.children[:idx], wg.children[idx+1:]...)
}

func (wg *WindowGroup) findWindowGroupEntry(name string) *windowGroupEntry {
	for _, c := range wg.children {
		if c.name == name {
			return c
		}
		if c.group != nil {
			wge := c.group.findWindowGroupEntry(name)
			if wge != nil {
				return wge
			}
		}
	}
	return nil
}

func (wg *WindowGroup) findWindowGroupEntryAndParent(name string) (*WindowGroup, *windowGroupEntry) {
	for _, c := range wg.children {
		if c.name == name {
			return wg, c
		}
		if c.group != nil {
			g, wge := c.group.findWindowGroupEntryAndParent(name)
			if wge != nil {
				return g, wge
			}
		}
	}
	return nil, nil
}

func (wg *WindowGroup) FindWindow(name string) Window {
	wge := wg.findWindowGroupEntry(name)
	if wge == nil {
		return nil
	}
	return wge.window
}

func (wg *WindowGroup) FindWindowGroup(name string) (*WindowGroup, error) {
	if (name == "root") && wg.isRoot {
		return wg, nil
	}
	wge := wg.findWindowGroupEntry(name)
	if wge == nil {
		return nil, fmt.Errorf("window group not found: %s", name)
	}
	if wge.group == nil {
		return nil, fmt.Errorf("not a window group: %s", name)
	}
	return wge.group, nil
}

func (wg *WindowGroup) RemoveAllChildren() {
	for i, c := range wg.children {
		if c.group != nil {
			c.group.RemoveAllChildren()
		}
		wg.children[i] = nil
	}
	wg.children = make([]*windowGroupEntry, 0)
}

func (wg *WindowGroup) MoveWindowOrGroup(src string, dst string, beginning bool) error {
	if src == "root" {
		return errors.New("can not move root window")
	}
	srcpg, srcwge := wg.findWindowGroupEntryAndParent(src)
	if srcwge == nil {
		return fmt.Errorf("source window not found: %s", src)
	}
	if srcwge.group != nil {
		check := srcwge.group.findWindowGroupEntry(dst)
		if check != nil {
			return fmt.Errorf("moving %s to %s would detach from root", src, dst)
		}
	}
	dstpg := wg
	if dst != "root" {
		var err error
		dstpg, err = wg.FindWindowGroup(dst)
		if err != nil {
			return err
		}
	}
	srcpg.removeChild(srcwge)
	if beginning {
		dstpg.children = append([]*windowGroupEntry{srcwge}, dstpg.children...)
	} else {
		dstpg.children = append(dstpg.children, srcwge)
	}
	wg.adjustNeeded = true
	return nil
}

func (wg *WindowGroup) SetWindowWeight(name string, w int) error {
	if w < 1 {
		return fmt.Errorf("weight must be >= 1: %d", w)
	}
	if w > 10000 {
		return fmt.Errorf("weight must be <= 10000: %d", w)
	}
	wge := wg.findWindowGroupEntry(name)
	if wge == nil {
		return fmt.Errorf("window not found: %s", name)
	}
	wge.weight = w
	wg.adjustNeeded = true
	return nil
}

func (wg *WindowGroup) AddWindowGroupChild(group *WindowGroup, name string, weight int) {
	wg.children = append(wg.children, &windowGroupEntry{name: name, weight: weight, group: group})
	wg.adjustNeeded = true
}

func (wg *WindowGroup) AddWindowChild(window Window, name string, weight int) {
	wg.children = append(wg.children, &windowGroupEntry{name: name, weight: weight, window: window})
	wg.adjustNeeded = true
}

func (wg *WindowGroup) UseColumnLayout(v bool) {
	wg.isColumn = v
	wg.adjustNeeded = true
}

func (wg *WindowGroup) Resize(x, y, w, h int) {
	wg.x = x
	wg.y = y
	wg.w = w
	wg.h = h
	wg.adjustNeeded = true
}

func (wg *WindowGroup) adjustChildren(screenw, screenh int) error {
	totalWeight := 0
	for _, c := range wg.children {
		totalWeight += c.weight
	}
	if wg.isColumn {
		wg.adjustChildrenColumn(totalWeight)
	} else {
		wg.adjustChildrenRow(totalWeight)
	}
	return wg.redrawChildBorders(screenw, screenh)
}

func (wg *WindowGroup) adjustChildrenColumn(totalWeight int) {
	x1 := wg.x
	for i, c := range wg.children {
		x2 := (i + 1) * (wg.w * c.weight / totalWeight)
		c.resize(x1, wg.y, x2-x1, wg.h)
		x1 = x2
	}
}

func (wg *WindowGroup) adjustChildrenRow(totalWeight int) {
	y1 := wg.y
	for i, c := range wg.children {
		y2 := (i + 1) * (wg.h * c.weight / totalWeight)
		c.resize(wg.x, y1, wg.w, y2-y1)
		y1 = y2
	}
}

func (wg *WindowGroup) redrawChildBorders(screenw, screenh int) error {
	for _, c := range wg.children {
		if err := c.showBorder(screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}

// Calls update on all contained windows
func (wg *WindowGroup) Update(rpn *rpn.RPN, updateInput bool, screenw, screenh int) error {
	if wg.adjustNeeded {
		if err := wg.adjustChildren(screenw, screenh); err != nil {
			return err
		}
		wg.adjustNeeded = false
		// We want to give screens other than the input screen a chance to
		// redraw before getting locked into the input screen
		if err := wg.Update(rpn, false, screenw, screenh); err != nil {
			return err
		}
	}
	if wg.isRoot && updateInput {
		// Update the input window first
		input := wg.FindWindow("i")
		if input == nil {
			return errors.New("could not find window 'i' for input")
		}
		if err := input.Update(rpn); err != nil {
			return err
		}
	}

	for _, c := range wg.children {
		if c.name == "i" {
			continue
		}
		if c.window != nil {
			if err := c.window.Update(rpn); err != nil {
				return err
			}
			continue
		}
		if err := c.group.Update(rpn, false, screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}
