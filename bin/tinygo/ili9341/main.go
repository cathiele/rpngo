// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"errors"
	"log"
	"mattwach/rpngo/drivers/pixelwinbuffer"
	"mattwach/rpngo/drivers/tinygo/ili9341"
	"mattwach/rpngo/drivers/tinygo/serialconsole"
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

	var screen ili9341.Ili9341Screen
	screen.Init()
	root, err := buildUI(&screen, &r)
	if err != nil {
		return err
	}
	newPixelPlotWindow := func() (window.WindowWithProps, error) {
		var ppw plotwin.PixelPlotWindow
		pw, err := screen.NewPixelWindow()
		if err != nil {
			return nil, err
		}
		var pb pixelwinbuffer.PixelBuffer
		pb.Init(pw)
		ppw.Init(&pb)
		return &ppw, nil
	}
	_ = commands.InitWindowCommands(&r, root, &screen, newPixelPlotWindow)
	_ = plotwin.InitPlotCommands(&r, root, &screen, plotwin.AddPixelPlotFn)
	if err := startup.LCD320Startup(&r); err != nil {
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
	txtw, err := screen.NewTextWindow()
	if err != nil {
		return err
	}
	gi := &getInput{}
	iw, err := input.Init(gi, txtw, r)
	gi.lcd = txtw.(*ili9341.Ili9341TxtW)
	if err != nil {
		return err
	}
	root.AddWindowChildToRoot(iw, "i", 100)
	return nil
}

type getInput struct {
	lcd     *ili9341.Ili9341TxtW
	serialc serialconsole.SerialConsole
}

func (gi *getInput) GetChar() (key.Key, error) {
	for {
		time.Sleep(20 * time.Millisecond)
		gi.lcd.ShowCursorIfEnabled(true)
		k := gi.serialc.GetChar()
		if k != 0 {
			return k, nil
		}
	}
}
