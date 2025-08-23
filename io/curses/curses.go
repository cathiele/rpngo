// Package ncurses implementation for IO
package curses

import (
	"strings"

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
	return &Curses{window: window}, nil
}

func (c *Curses) End() {
	goncurses.End()
}

func (c *Curses) ReadByte() (byte, error) {
	return byte(c.window.GetChar()), nil
}

func (c *Curses) Clear() error {
	if err := c.window.Clear(); err != nil {
		return err
	}
	c.window.Refresh()
	return nil
}

func (c *Curses) Write(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	// filter out \n so that we emulate a dumb screen
	idx := strings.Index(s, "\n")
	if idx >= 0 {
		if err := c.Write(b[:idx]); err != nil {
			return err
		}
		return c.Write(b[idx+1:])
	}
	c.window.Print(s)
	c.window.Refresh()
	return nil
}

func (c *Curses) Width() uint {
	_, x := c.window.MaxYX()
	return uint(x)
}

func (c *Curses) Height() uint {
	y, _ := c.window.MaxYX()
	return uint(y)
}

func (c *Curses) X() uint {
	_, x := c.window.CursorYX()
	return uint(x)
}

func (c *Curses) Y() uint {
	y, _ := c.window.CursorYX()
	return uint(y)
}

func (c *Curses) SetX(x uint) {
	c.window.Move(int(c.Y()), int(x))
}

func (c *Curses) SetY(y uint) {
	c.window.Move(int(y), int(c.X()))
}
