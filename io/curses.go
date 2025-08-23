// ncurses implementation for IO
package curses

import "github.com/gbin/goncurses"

type Curses struct {
	window *goncurses.Window
}

func Init() (*Curses, error) {
	window, err := goncurses.Init()
	if err != nil {
		return nil, err
	}
	return &Curses{window: window}, nil
}

func (c *Curses) End() {
	c.End()
}

func (c *Curses) ReadByte() (byte, error) {
	return byte(c.window.GetChar()), nil
}

func (c *Curses) Clear() error {
	return c.window.Clear()
}

func (c *Curses) Write(b []byte) error {
	c.window.Print(string(b))
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
	_, x := c.window.YX()
	return uint(x)
}

func (c *Curses) Y() uint {
	y, _ := c.window.YX()
	return uint(y)
}

func (c *Curses) SetX(x uint) {
	c.window.Move(int(c.Y()), int(x))
}

func (c *Curses) SetY(y uint) {
	c.window.Move(int(y), int(c.X()))
}
