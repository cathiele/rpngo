// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"errors"
	"log"
	"mattwach/rpngo/drivers/pixelwinbuffer"
	"mattwach/rpngo/drivers/tinygo/picocalc/i2ckbd"
	"mattwach/rpngo/drivers/tinygo/picocalc/ili948x"
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
	var screen ili948x.Ili948xScreen // object allocated on the heap (OK)
	screen.Init()
	// This only seeems to work for panics I throw and not errors
	// like array out of bounds.
	defer func() {
		if r := recover(); r != nil {
			handlePanic(screen.Device, r)
		}
	}()
	if err := run(&screen); err != nil {
		panic(err)
	}
}

func run(screen *ili948x.Ili948xScreen) error { // object allocated on the heap (OK)
	time.Sleep(200 * time.Millisecond)

	log.SetOutput(os.Stdout)
	log.Println("Started") // object allocated on the heap (OK)
	var r rpn.RPN          // object allocated on the heap: (OK)
	r.Init()
	functions.RegisterAll(&r)

	root, err := buildUI(screen, &r)
	if err != nil {
		return err
	}
	newPixelPlotWindow := func() (window.WindowWithProps, error) {
		var ppw plotwin.PixelPlotWindow // object allocated on the heap (OK)
		pw, err := screen.NewPixelWindow()
		if err != nil {
			return nil, err
		}
		var pb pixelwinbuffer.PixelBuffer // object allocated on the heap (OK)
		pb.Init(pw)
		ppw.Init(&pb)
		return &ppw, nil
	}
	_ = commands.InitWindowCommands(&r, root, screen, newPixelPlotWindow)
	_ = plotwin.InitPlotCommands(&r, root, screen, plotwin.AddPixelPlotFn)
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
	gi := &getInput{} // object allocated on the heap (OK)
	gi.Init()
	iw, err := input.Init(gi, txtw, r)
	gi.lcd = txtw.(*ili948x.Ili948xTxtW)
	if err != nil {
		return err
	}
	root.AddWindowChildToRoot(iw, "i", 100)
	return nil
}

type getInput struct {
	lcd      *ili948x.Ili948xTxtW
	serialc  serialconsole.SerialConsole
	keyboard i2ckbd.I2CKbd
}

func (gi *getInput) Init() {
	if err := gi.keyboard.Init(); err != nil {
		log.Printf("failed to init keyboard: %v", err) // object allocated on the heap (OK)
	}
}

func (gi *getInput) GetChar() (key.Key, error) {
	for {
		time.Sleep(20 * time.Millisecond)
		gi.lcd.ShowCursorIfEnabled(true)
		k := gi.serialc.GetChar()
		if k != 0 {
			return k, nil
		}
		var err error
		k, err = gi.keyboard.GetChar()
		if err != nil {
			log.Printf("keyboard error: %v", err) // object allocated on the heap (OK)
			time.Sleep(1 * time.Second)
		} else if k != 0 {
			return k, nil
		}
	}
}
