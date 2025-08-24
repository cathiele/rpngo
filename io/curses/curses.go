// Package ncurses implementation for IO
package curses

import (
	"errors"
	"log"
	"mattwach/rpngo/io/key"

	"github.com/gbin/goncurses"
)

type Curses struct {
	window     *goncurses.Window
	rgbToPair  map[uint32]int16 // maps rgbrgb values to a pair.
	rgbToColor map[uint16]int16 // maps rgb to a color index
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
	if err := window.Keypad(true); err != nil {
		return nil, err
	}
	return &Curses{
		window:     window,
		rgbToPair:  make(map[uint32]int16),
		rgbToColor: make(map[uint16]int16)}, nil
}

func (c *Curses) End() {
	goncurses.End()
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

func (c *Curses) Clear() error {
	return c.window.Clear()
}

func (c *Curses) Refresh() {
	c.window.Refresh()
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
	log.Printf("AttrSet: pairIdx=%v ch=%v", pairIdx, ch)
	return nil
}

func (c *Curses) colorPairFor(fr, fg, fb, br, bg, bb int) (int16, error) {
	fc := (uint32(fr) << 10) | (uint32(fg) << 5) | uint32(fb)
	bc := (uint32(br) << 10) | (uint32(bg) << 5) | uint32(bb)
	pc := (fc << 15) | bc
	pidx, ok := c.rgbToPair[pc]
	if !ok {
		var err error
		pidx, err = c.createNewPair(fr, fg, fb, br, bg, bb)
		if err != nil {
			return 0, err
		}
		c.rgbToPair[pc] = pidx
	}
	return pidx, nil
}

func (c *Curses) createNewPair(fr, fg, fb, br, bg, bb int) (int16, error) {
	fIdx, err := c.colorIndexFor(fr, fg, fb)
	if err != nil {
		return 0, err
	}
	bIdx, err := c.colorIndexFor(br, bg, bb)
	if err != nil {
		return 0, err
	}
	pIdx := int16(len(c.rgbToPair) + 1) // zero is the default so start at 1
	err = goncurses.InitPair(pIdx, fIdx, bIdx)
	log.Printf("initPair: pidx=%v, fidx=%v, bidx=%v", pIdx, fIdx, bIdx)
	if err != nil {
		return 0, err
	}
	return pIdx, nil
}

func (c *Curses) colorIndexFor(r, g, b int) (int16, error) {
	key := (uint16(r) << 10) | (uint16(g) << 5) | uint16(b)
	idx, ok := c.rgbToColor[key]
	if !ok {
		idx = int16(len(c.rgbToColor)) + 50
		err := goncurses.InitColor(
			idx,
			int16(r*1000/31),
			int16(g*1000/31),
			int16(b*1000/31),
		)
		log.Printf("InitColor: idx=%v, r=%v, g=%v, b=%v",
			idx,
			int16(r*1000/31),
			int16(g*1000/31),
			int16(b*1000/31),
		)
		if err != nil {
			return 0, err
		}
		c.rgbToColor[key] = idx
	}
	return idx, nil
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
