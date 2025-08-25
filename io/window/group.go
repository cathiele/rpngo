package window

type windowGroupEntry struct {
	name   string
	weight int
	// Only one of the following should be non-nil
	group *WindowGroup
	txtw  TextWindow
}

func (wge *windowGroupEntry) resize(x, y, w, h int) {
	if wge.group != nil {
		wge.group.Resize(x, y, w, h)
		return
	}
	if wge.txtw != nil {
		wge.txtw.Resize(x, y, w, h)
	}
}

type WindowGroup struct {
	isVertical bool
	// Coordinates are in global screen coordinates
	x        int
	y        int
	w        int
	h        int
	children []*windowGroupEntry
}

func NewWindowGroup() *WindowGroup {
	return &WindowGroup{}
}

func (wg *WindowGroup) FindTextWindow(name string) TextWindow {
	for _, c := range wg.children {
		if c.name == name {
			return c.txtw
		}
		if c.group != nil {
			txtw := c.group.FindTextWindow(name)
			if txtw != nil {
				return txtw
			}
		}
	}
	return nil
}

func (wg *WindowGroup) AddWindowGroupChild(group *WindowGroup, name string, weight int) {
	wg.children = append(wg.children, &windowGroupEntry{name: name, weight: weight, group: group})
	wg.adjustChildren()
}

func (wg *WindowGroup) AddTextWindowChild(txtw TextWindow, name string, weight int) {
	wg.children = append(wg.children, &windowGroupEntry{name: name, weight: weight, txtw: txtw})
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
