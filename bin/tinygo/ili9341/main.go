// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"errors"
	"log"
	"machine"
	"mattwach/rpngo/drivers/tinygo/ili9341tw"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/startup"
	"mattwach/rpngo/window"
	"mattwach/rpngo/window/commands"
	"mattwach/rpngo/window/input"
	"mattwach/rpngo/window/plotwin"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	time.Sleep(2 * time.Second)

	log.SetOutput(os.Stdout)
	log.Println("Started")
	var r rpn.RPN
	r.Init()
	functions.RegisterAll(&r)

	var screen ili9341tw.Ili9341Screen
	screen.Init()
	root, err := buildUI(&screen, &r)
	if err != nil {
		return err
	}
	_ = commands.InitWindowCommands(&r, root, &screen)
	_ = plotwin.InitPlotCommands(&r, root, &screen)
	if err := startup.OSStartup(&r); err != nil {
		return err
	}
	w, h := screen.ScreenSize()
	if err := root.Update(&r, w, h, true); err != nil {
		return err
	}
	for {
		w, h = screen.ScreenSize()
		if err := root.Update(&r, w, h, true); err != nil {
			if errors.Is(err, input.ErrExit) {
				return nil
			}
			return err
		}
	}
}

func buildUI(screen window.Screen, r *rpn.RPN) (*window.WindowRoot, error) {
	w, h := screen.ScreenSize()
	root := window.NewWindowRoot(w, h)
	if err := addInputWindow(screen, root, r); err != nil {
		return nil, err
	}
	return root, nil
}

func addInputWindow(screen window.Screen, root *window.WindowRoot, r *rpn.RPN) error {
	w, h := screen.ScreenSize()
	txtw, err := screen.NewTextWindow(0, 0, w, h)
	if err != nil {
		return err
	}
	gi := &getInput{}
	iw, err := input.Init(gi, txtw, r)
	gi.lcd = txtw.(*ili9341tw.Ili9341TW)
	if err != nil {
		return err
	}
	root.AddWindowChild(iw, "i", 100)
	return nil
}

type getInput struct {
	lcd *ili9341tw.Ili9341TW
}

type TermState int

const (
	NORMAL TermState = iota
	ESC
	ARROW
)

func (g *getInput) GetChar() (key.Key, error) {
	var state TermState = NORMAL
	for {
		g.lcd.ShowCursorIfEnabled(true)
		c, err := machine.Serial.ReadByte()
		if err != nil {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		//log.Printf("got char: %v", c)
		switch state {
		case NORMAL:
			switch c {
			case 13:
				return '\n', nil
			case 27:
				state = ESC
			case 127:
				return key.KEY_BACKSPACE, nil
			default:
				return key.Key(c), nil
			}
		case ESC:
			switch c {
			case '[':
				state = ARROW
			default:
				state = NORMAL
			}
		case ARROW:
			state = NORMAL
			switch c {
			case 'A':
				return key.KEY_UP, nil
			case 'B':
				return key.KEY_DOWN, nil
			case 'C':
				return key.KEY_RIGHT, nil
			case 'D':
				return key.KEY_LEFT, nil
			}
		}
	}
}
