// Package ncurses implementation for IO
package curses

import (
	"errors"
	"mattwach/rpngo/io/key"
	"mattwach/rpngo/io/window"

	"github.com/gbin/goncurses"
)

type Curses struct {
	border    *goncurses.Window
	window    *goncurses.Window
	rgbToPair map[uint32]int16 // maps color,color values to a pair.
}

func Init() (*Curses, error) {
	window, err := goncurses.Init()
	if err != nil {
		return nil, err
	}
	if goncurses.HasColors() {
		goncurses.StartColor()
	}
	goncurses.Echo(false)
	tw := &Curses{
		window:    window,
		rgbToPair: make(map[uint32]int16),
	}
	return tw, nil
}

func (c *Curses) NewTextWindow(x, y, w, h int) (window.TextWindow, error) {
	tw := &Curses{
		rgbToPair: c.rgbToPair,
	}
	if err := tw.Resize(x, y, w, h); err != nil {
		return nil, err
	}
	return tw, nil
}

func (c *Curses) ShowBorder(top, bottom, left, right bool) error {
	ls := ' '
	rs := ' '
	ts := ' '
	bs := ' '
	tl := ' '
	tr := ' '
	bl := ' '
	br := ' '
	if top {
		ts = '-'
		tr = '-'
		tl = '-'
	}
	if bottom {
		bs = '-'
		br = '-'
		bl = '-'
	}
	if left {
		ls = '|'
		tl = '|'
		bl = '|'
	}
	if right {
		rs = '|'
		tr = '|'
		br = '|'
	}
	if top && left {
		tl = '+'
	}
	if top && right {
		tr = '+'
	}
	if bottom && left {
		bl = '+'
	}
	if bottom && right {
		br = '+'
	}
	if err := c.border.Border(
		goncurses.Char(ls),
		goncurses.Char(rs),
		goncurses.Char(ts),
		goncurses.Char(bs),
		goncurses.Char(tl),
		goncurses.Char(tr),
		goncurses.Char(bl),
		goncurses.Char(br),
	); err != nil {
		return err
	}
	c.border.Refresh()
	return nil
}

func (c *Curses) Refresh() {
	c.window.Refresh()
}

func (c *Curses) End() {
	goncurses.End()
}

func (c *Curses) Resize(x, y, w, h int) error {
	if c.border != nil {
		// erase the contents so that artifacts do no collect on the screen
		c.border.Erase()
		c.border.Refresh()
		if err := c.border.Delete(); err != nil {
			return err
		}
	}
	if c.window != nil {
		if err := c.window.Delete(); err != nil {
			return err
		}
	}
	var err error
	c.border, err = goncurses.NewWindow(h, w, y, x)
	if err != nil {
		return err
	}
	c.window, err = goncurses.NewWindow(h-2, w-2, y+1, x+1)
	if err != nil {
		return err
	}
	if err := c.window.Keypad(true); err != nil {
		return err
	}
	return nil
}

var charMap = map[goncurses.Key]key.Key{
	goncurses.KEY_LEFT:      key.KEY_LEFT,
	goncurses.KEY_RIGHT:     key.KEY_RIGHT,
	goncurses.KEY_UP:        key.KEY_UP,
	goncurses.KEY_DOWN:      key.KEY_DOWN,
	goncurses.KEY_BACKSPACE: key.KEY_BACKSPACE,
	goncurses.KEY_DC:        key.KEY_DEL,
	goncurses.KEY_IC:        key.KEY_INS,
	goncurses.KEY_END:       key.KEY_END,
	goncurses.KEY_HOME:      key.KEY_HOME,
	4:                       key.KEY_EOF,
}

func (c *Curses) GetChar() (key.Key, error) {
	ch := c.window.GetChar()
	k, ok := charMap[ch]
	if ok {
		return k, nil
	}
	return key.Key(ch), nil
}

func (c *Curses) Erase() {
	c.window.Erase()
}

func (c *Curses) Write(b byte) error {
	if b == '\n' {
		return c.newLine()
	}
	c.window.AddChar(goncurses.Char(b))
	return nil
}

func (c *Curses) newLine() error {
	y := c.Y()
	h := c.Height()
	if y < (h - 1) {
		y++
	} else {
		y = h - 1
		c.Scroll(1)
	}
	c.SetXY(0, y)
	return nil
}

func (c *Curses) Width() int {
	_, x := c.window.MaxYX()
	return x
}

func (c *Curses) Height() int {
	y, _ := c.window.MaxYX()
	return y
}

func (c *Curses) Size() (int, int) {
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) X() int {
	_, x := c.window.CursorYX()
	return x
}

func (c *Curses) Y() int {
	y, _ := c.window.CursorYX()
	return y
}

func (c *Curses) XY() (int, int) {
	y, x := c.window.CursorYX()
	return x, y
}

func (c *Curses) SetX(x int) {
	c.window.Move(c.Y(), x)
}

func (c *Curses) SetY(y int) {
	c.window.Move(y, c.X())
}

func (c *Curses) SetXY(x int, y int) {
	c.window.Move(y, x)
}

func (c *Curses) Scroll(n int) {
	c.window.ScrollOk(true)
	c.window.Scroll(n)
}

func (c *Curses) Color(fr, fg, fb, br, bg, bb int) error {
	if err := checkColorRange(fr, fg, fb); err != nil {
		return err
	}
	if err := checkColorRange(br, bg, bb); err != nil {
		return err
	}
	if !goncurses.HasColors() {
		return nil
	}
	pairIdx, err := c.colorPairFor(fr, fg, fb, br, bg, bb)
	if err != nil {
		return err
	}
	ch := goncurses.ColorPair(pairIdx)
	c.window.AttrSet(ch)
	return nil
}

func (c *Curses) colorPairFor(fr, fg, fb, br, bg, bb int) (int16, error) {
	fc := colorIndexFor(fr, fg, fb)
	bc := colorIndexFor(br, bg, bb)
	pc := (uint32(fc) << 15) | uint32(bc)
	pidx, ok := c.rgbToPair[pc]
	if !ok {
		pidx = int16(len(c.rgbToPair) + 1) // zero is the default so start at 1
		if err := goncurses.InitPair(pidx, fc, bc); err != nil {
			return 0, err
		}
		c.rgbToPair[pc] = pidx
	}
	return pidx, nil
}

var idxToCol = map[uint8]int16{
	0: goncurses.C_BLACK,
	1: goncurses.C_BLUE,
	2: goncurses.C_GREEN,
	3: goncurses.C_CYAN,
	4: goncurses.C_RED,
	5: goncurses.C_MAGENTA,
	6: goncurses.C_YELLOW,
	7: goncurses.C_WHITE,
}

func colorIndexFor(r, g, b int) int16 {
	var v uint8 = 0
	if r > 15 {
		v |= 4
	}
	if g > 15 {
		v |= 2
	}
	if b > 15 {
		v |= 1
	}
	col := idxToCol[v]
	return col
}

func checkColorRange(r, g, b int) error {
	if r < 0 {
		return errors.New("red value < 0")
	}
	if r > 31 {
		return errors.New("red value > 31")
	}
	if g < 0 {
		return errors.New("green value < 0")
	}
	if g > 31 {
		return errors.New("green value > 31")
	}
	if b < 0 {
		return errors.New("blue value < 0")
	}
	if b > 31 {
		return errors.New("blue value > 31")
	}
	return nil
}
