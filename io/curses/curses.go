// Package ncurses implementation for IO
package curses

import (
	"mattwach/rpngo/io/key"

	"github.com/gbin/goncurses"
)

type Curses struct {
	window *goncurses.Window
}

func Init() (*Curses, error) {
	window, err := goncurses.Init()
	if err != nil {
		return nil, err
	}
	goncurses.Echo(false)
	if err := window.Keypad(true); err != nil {
		return nil, err
	}
	return &Curses{window: window}, nil
}

func (c *Curses) End() {
	goncurses.End()
}

func (c *Curses) GetChar() (key.Key, error) {
	ch := c.window.GetChar()
	switch ch {
	case goncurses.KEY_LEFT:
		return key.KEY_LEFT, nil
	case goncurses.KEY_RIGHT:
		return key.KEY_RIGHT, nil
	default:
		return key.Key(ch), nil
	}
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
	c.window.AddChar(goncurses.Char((b)))
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
