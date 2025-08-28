package window

import (
	"errors"
	"mattwach/rpngo/rpn"
)

type Window interface {
	Update(*rpn.RPN) error
	Resize(x, y, w, h int) error
	ShowBorder(top, bottom, left, right bool) error
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

func (wge *windowGroupEntry) showBorder(isColumn, isFinal bool) error {
	if wge.group != nil {
		for i, c := range wge.group.children {
			if err := c.showBorder(
				isColumn, i == (len(wge.group.children)-1)); err != nil {
				return err
			}
		}
	} else if wge.window != nil {
		if err := wge.window.ShowBorder(
			false,
			!isColumn && !isFinal,
			false,
			isColumn && !isFinal); err != nil {
			return err
		}
	}
	return nil
}

type WindowGroup struct {
	isRoot   bool
	isColumn bool
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

func (wg *WindowGroup) FindWindow(name string) Window {
	for _, c := range wg.children {
		if c.name == name {
			return c.window
		}
		if c.group != nil {
			window := c.group.FindWindow(name)
			if window != nil {
				return window
			}
		}
	}
	return nil
}

func (wg *WindowGroup) AddWindowGroupChild(group *WindowGroup, name string, weight int) {
	wg.children = append(wg.children, &windowGroupEntry{name: name, weight: weight, group: group})
	wg.adjustChildren()
}

func (wg *WindowGroup) AddWindowChild(window Window, name string, weight int) {
	wg.children = append(wg.children, &windowGroupEntry{name: name, weight: weight, window: window})
	wg.adjustChildren()
}

func (wg *WindowGroup) UseColumnLayout(v bool) {
	wg.isColumn = v
	wg.adjustChildren()
}

func (wg *WindowGroup) Resize(x, y, w, h int) error {
	wg.x = x
	wg.y = y
	wg.w = w
	wg.h = h
	return wg.adjustChildren()
}

func (wg *WindowGroup) adjustChildren() error {
	totalWeight := 0
	for _, c := range wg.children {
		totalWeight += c.weight
	}
	if wg.isColumn {
		wg.adjustChildrenColumn(totalWeight)
	} else {
		wg.adjustChildrenRow(totalWeight)
	}
	return wg.redrawChildBorders()
}

func (wg *WindowGroup) adjustChildrenColumn(totalWeight int) {
	x1 := wg.x
	for _, c := range wg.children {
		x2 := x1 + (wg.w * c.weight / totalWeight)
		c.resize(x1, wg.y, x2-x1, wg.h)
		x1 = x2
	}
}

func (wg *WindowGroup) adjustChildrenRow(totalWeight int) {
	y1 := wg.y
	for _, c := range wg.children {
		y2 := y1 + (wg.h * c.weight / totalWeight)
		c.resize(wg.x, y1, wg.w, y2-y1)
		y1 = y2
	}
}

func (wg *WindowGroup) redrawChildBorders() error {
	for i, c := range wg.children {
		if err := c.showBorder(
			wg.isColumn, i == (len(wg.children)-1)); err != nil {
			return err
		}
	}
	return nil
}

// Calls update on all contained windows
func (wg *WindowGroup) Update(rpn *rpn.RPN, updateInput bool) error {
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
		if err := c.group.Update(rpn, false); err != nil {
			return err
		}
	}
	return nil
}
