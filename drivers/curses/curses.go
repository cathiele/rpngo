// Package ncurses implementation for IO
package curses

import (
	"mattwach/rpngo/key"
	"mattwach/rpngo/window"

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
	goncurses.Cursor(0)
	tw := &Curses{
		window:    window,
		rgbToPair: make(map[uint32]int16),
	}
	return tw, nil
}

func (c *Curses) NewTextWindow() (window.TextWindow, error) {
	tw := &Curses{
		rgbToPair: c.rgbToPair,
	}
	if err := tw.ResizeWindow(0, 0, 8, 8); err != nil {
		return nil, err
	}
	return tw, nil
}

func (c *Curses) NewPixelWindow() (window.PixelWindow, error) {
	panic("ncurses does not support pixel-based windows")
}

func (c *Curses) ShowBorder(screenw, screenh int) error {
	ch, err := c.colorPairFor(window.Magenta)
	if err != nil {
		return err
	}
	c.border.AttrSet(ch)
	if err := c.border.Border('|', '|', '-', '-', '+', '+', '+', '+'); err != nil {
		return err
	}
	c.border.Refresh()
	ch, err = c.colorPairFor(window.White)
	if err != nil {
		return err
	}
	c.border.AttrSet(ch)
	return nil
}

func (c *Curses) Refresh() {
	c.window.Refresh()
}

func (c *Curses) End() {
	goncurses.End()
}

func (c *Curses) Cursor(on bool) {
	var val byte = 0
	if on {
		val = 1
	}
	goncurses.Cursor(val)
}

func (c *Curses) ResizeWindow(x, y, w, h int) error {
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
	goncurses.KEY_F1:        key.KEY_F1,
	goncurses.KEY_F2:        key.KEY_F2,
	goncurses.KEY_F3:        key.KEY_F3,
	goncurses.KEY_F4:        key.KEY_F4,
	goncurses.KEY_F5:        key.KEY_F5,
	goncurses.KEY_F6:        key.KEY_F6,
	goncurses.KEY_F7:        key.KEY_F7,
	goncurses.KEY_F8:        key.KEY_F8,
	goncurses.KEY_F9:        key.KEY_F9,
	goncurses.KEY_F10:       key.KEY_F10,
	goncurses.KEY_F11:       key.KEY_F11,
	goncurses.KEY_F12:       key.KEY_F11,
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
	y := c.CursorY()
	h := c.TextHeight()
	if y < (h - 1) {
		y++
	} else {
		y = h - 1
		c.Scroll(1)
	}
	c.SetCursorXY(0, y)
	return nil
}

func (c *Curses) TextWidth() int {
	_, x := c.window.MaxYX()
	return x
}

func (c *Curses) TextHeight() int {
	y, _ := c.window.MaxYX()
	return y
}

func (c *Curses) TextSize() (int, int) {
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) WindowSize() (int, int) {
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) ScreenSize() (int, int) {
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) WindowXY() (int, int) {
	y, x := c.window.YX()
	return x, y
}

func (c *Curses) CursorX() int {
	_, x := c.window.CursorYX()
	return x
}

func (c *Curses) CursorY() int {
	y, _ := c.window.CursorYX()
	return y
}

func (c *Curses) CursorXY() (int, int) {
	y, x := c.window.CursorYX()
	return x, y
}

func (c *Curses) SetCursorX(x int) {
	c.window.Move(c.CursorY(), x)
}

func (c *Curses) SetCursorY(y int) {
	c.window.Move(y, c.CursorX())
}

func (c *Curses) SetCursorXY(x int, y int) {
	c.window.Move(y, x)
}

func (c *Curses) Scroll(n int) {
	c.window.ScrollOk(true)
	c.window.Scroll(n)
}

func (c *Curses) TextColor(col window.ColorChar) {
	ch, err := c.colorPairFor(col)
	if err != nil {
		return
	}
	c.window.AttrSet(ch)
}

func (c *Curses) colorPairFor(col window.ColorChar) (goncurses.Char, error) {
	if !goncurses.HasColors() {
		return 0, nil
	}
	fc := idxToCol[uint8(col>>12)]
	bc := idxToCol[uint8((col&0x0F00)>>8)]
	pc := (uint32(fc) << 15) | uint32(bc)
	pidx, ok := c.rgbToPair[pc]
	if !ok {
		pidx = int16(len(c.rgbToPair) + 1) // zero is the default so start at 1
		if err := goncurses.InitPair(pidx, fc, bc); err != nil {
			return 0, err
		}
		c.rgbToPair[pc] = pidx
	}
	return goncurses.ColorPair(pidx), nil
}

var idxToCol = map[uint8]int16{
	0b0000: goncurses.C_BLACK,
	0b0001: goncurses.C_BLUE,
	0b0010: goncurses.C_GREEN,
	0b0011: goncurses.C_BLUE,
	0b0100: goncurses.C_GREEN,
	0b0101: goncurses.C_CYAN,
	0b0110: goncurses.C_GREEN,
	0b0111: goncurses.C_CYAN,
	0b1000: goncurses.C_RED,
	0b1001: goncurses.C_MAGENTA,
	0b1010: goncurses.C_RED,
	0b1011: goncurses.C_MAGENTA,
	0b1100: goncurses.C_YELLOW,
	0b1101: goncurses.C_WHITE,
	0b1110: goncurses.C_YELLOW,
	0b1111: goncurses.C_WHITE,
}
