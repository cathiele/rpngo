package window

type WindowGroupEntry struct {
	weight int
	// Only one of the following should be non-nil
	group *WindowGroup
	txtw  *TextWindow
}

type WindowGroup struct {
	isVertical bool
	x          int
	y          int
	w          int
	h          int
	children   []WindowGroupEntry
}

func NewWindowGroup(x, y, w, h int) *WindowGroup {
	wg := &WindowGroup{}
	wg.Resize(x, y, w, h)
	return wg
}

func (wg *WindowGroup) Resize(x, y, w, h int) {
	// TODO: implement me
}
