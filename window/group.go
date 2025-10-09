package window

import (
	"fmt"
	"mattwach/rpngo/rpn"
	"strings"
)

type WindowWithProps interface {
	WindowBase
	Update(*rpn.RPN) error
	Type() string
	SetProp(name string, val rpn.Frame) error
	GetProp(name string) (rpn.Frame, error)
	ListProps() []string
}

type windowGroupEntry struct {
	name   string
	weight int
	// Only one of the following should be non-nil
	group  *windowGroup
	window WindowWithProps
}

func (wge *windowGroupEntry) resize(x, y, w, h int) {
	if wge.group != nil {
		wge.group.resize(x, y, w, h)
	} else if wge.window != nil {
		wge.window.ResizeWindow(x, y, w, h)
	}
}

func (wge *windowGroupEntry) showBorder(screenw, screenh int) error {
	//log.Printf("showBorder: name=%s sw=%d, sh=%d", wge.name, screenw, screenh)
	if wge.group != nil {
		wge.group.showBorder(screenw, screenh)
	} else if wge.window != nil {
		if err := wge.window.ShowBorder(screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}

type windowGroup struct {
	isColumn bool
	// Coordinates are in global screen coordinates
	x        int
	y        int
	w        int
	h        int
	children []*windowGroupEntry
}

func (wg *windowGroup) removeChild(c *windowGroupEntry) {
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

func (wg *windowGroup) findwindowGroupEntry(name string) *windowGroupEntry {
	for _, c := range wg.children {
		if c.name == name {
			return c
		}
		if c.group != nil {
			wge := c.group.findwindowGroupEntry(name)
			if wge != nil {
				return wge
			}
		}
	}
	return nil
}

func (wg *windowGroup) findwindowGroupEntryAndParent(name string) (*windowGroup, *windowGroupEntry) {
	for _, c := range wg.children {
		if c.name == name {
			return wg, c
		}
		if c.group != nil {
			g, wge := c.group.findwindowGroupEntryAndParent(name)
			if wge != nil {
				return g, wge
			}
		}
	}
	return nil, nil
}

func (wg *windowGroup) resize(x, y, w, h int) {
	wg.x = x
	wg.y = y
	wg.w = w
	wg.h = h
}

func (wg *windowGroup) dump(lines []string, name string, indent int, weight int) []string {  // object allocated on the heap: escapes at line 110
	pad := strings.Repeat("  ", indent)  // object allocated on the heap: escapes at line 109
	line := fmt.Sprintf(
		"%s%s(x=%d, y=%d, w=%d, h=%d, cols=%v, weight=%d):",
		pad,
		name,
		wg.x,
		wg.y,
		wg.w,
		wg.h,
		wg.isColumn,
		weight,
	)
	lines = append(lines, line)
	pad = strings.Repeat("  ", indent+1)  // object allocated on the heap: escapes at line 131
	for _, c := range wg.children {
		if c.group != nil {
			lines = c.group.dump(lines, c.name, indent+1, c.weight)
		}
		if c.window != nil {
			x, y := c.window.WindowXY()
			w, h := c.window.WindowSize()
			lines = append(
				lines,
				fmt.Sprintf(
					"%s%s(type=%s, x=%d, y=%d, w=%d, h=%d, weight=%d)",
					pad,
					c.name,  // object allocated on the heap: escapes at line 132
					c.window.Type(),  // object allocated on the heap: escapes at line 133
					x,
					y,
					w,
					h,
					c.weight,
				))
		}
	}
	return lines
}

func (wg *windowGroup) removeAllChildren() {
	for i, c := range wg.children {
		if c.group != nil {
			c.group.removeAllChildren()
		}
		wg.children[i] = nil
	}
	wg.children = make([]*windowGroupEntry, 0)
}

func (wg *windowGroup) adjustChildren(screenw, screenh int) error {
	//log.Printf("adjustChildren: wg=%v sw=%d sh=%d", wg, screenw, screenh)
	totalWeight := 0
	for _, c := range wg.children {
		totalWeight += c.weight
	}
	if wg.isColumn {
		wg.adjustChildrenColumn(screenw, screenh, totalWeight)
	} else {
		wg.adjustChildrenRow(screenw, screenh, totalWeight)
	}
	var err error
	return err
}

func (wg *windowGroup) adjustChildrenColumn(screenw, screenh, totalWeight int) {
	x1 := wg.x
	weightSum := 0
	for _, c := range wg.children {
		weightSum += c.weight
		x2 := wg.x + wg.w*weightSum/totalWeight
		c.resize(x1, wg.y, x2-x1, wg.h)
		if c.group != nil {
			c.group.adjustChildren(screenw, screenh)
		}
		x1 = x2 - 1
	}
}

func (wg *windowGroup) adjustChildrenRow(screenw, screenh, totalWeight int) {
	y1 := wg.y
	weightSum := 0
	for _, c := range wg.children {
		weightSum += c.weight
		y2 := wg.y + wg.h*weightSum/totalWeight
		//log.Printf("adjustChildrenRow: y1=%d y2=%d weightSum=%d totalWeight=%d", y1, y2, weightSum, totalWeight)
		c.resize(wg.x, y1, wg.w, y2-y1)
		if c.group != nil {
			c.group.adjustChildren(screenw, screenh)
		}
		y1 = y2 - 1
	}
}

func (wg *windowGroup) showBorder(screenw, screenh int) error {
	for _, c := range wg.children {
		if err := c.showBorder(screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}

// Calls update on all contained windows
func (wg *windowGroup) update(rpn *rpn.RPN, screenw, screenh int) error {
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
		if err := c.group.update(rpn, screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}