package window

import (
	"errors"
	"mattwach/rpngo/rpn"
)

type Window interface {
	Update(*rpn.RPN) error
	Resize(x, y, w, h int)
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
		return
	}
	if wge.window != nil {
		wge.window.Resize(x, y, w, h)
	}
}

type WindowGroup struct {
	isRoot     bool
	isVertical bool
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

func (wg *WindowGroup) setVertical(v bool) {
	wg.isVertical = v
	wg.adjustChildren()
}

func (wg *WindowGroup) Resize(x, y, w, h int) {
	wg.x = x
	wg.y = y
	wg.w = w
	wg.h = h
	wg.adjustChildren()
}

func (wg *WindowGroup) adjustChildren() {
	totalWeight := 0
	for _, c := range wg.children {
		totalWeight += c.weight
	}
	if wg.isVertical {
		wg.adjustChildrenVertical(totalWeight)
	} else {
		wg.adjustChildrenHorizontal(totalWeight)
	}
}

func (wg *WindowGroup) adjustChildrenVertical(totalWeight int) {
	x1 := wg.x
	for _, c := range wg.children {
		x2 := x1 + (wg.w * c.weight / totalWeight)
		c.resize(wg.x, wg.y, x2-x1, wg.h)
		x1 = x2
	}
}

func (wg *WindowGroup) adjustChildrenHorizontal(totalWeight int) {
	y1 := wg.y
	for _, c := range wg.children {
		y2 := y1 + (wg.h * c.weight / totalWeight)
		c.resize(wg.x, wg.y, wg.w, y2-y1)
		y1 = y2
	}
}

// Calls update on all contained windows
func (wg *WindowGroup) Update(rpn *rpn.RPN) error {
	if wg.isRoot {
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
		if err := c.group.Update(rpn); err != nil {
			return err
		}
	}
	return nil
}
