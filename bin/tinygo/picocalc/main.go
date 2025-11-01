// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"errors"
	"mattwach/rpngo/drivers/pixelwinbuffer"
	"mattwach/rpngo/drivers/tinygo/picocalc/i2ckbd"
	"mattwach/rpngo/drivers/tinygo/picocalc/ili948x"
	"mattwach/rpngo/drivers/tinygo/serial"
	"mattwach/rpngo/drivers/tinygo/tinyfs"
	"mattwach/rpngo/elog"
	"mattwach/rpngo/fileops"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/startup"
	"mattwach/rpngo/window"
	"mattwach/rpngo/window/commands"
	"mattwach/rpngo/window/input"
	"mattwach/rpngo/window/plotwin"
	"time"
)

const scrollbytes = 8 * 1024
const maxStackDepth = 256

// Persistant globals.  Tinygo curently can't move heap pointers,
// which leads to heap fragmentation and thus the heap should
// be used sparingly
var rpnInst rpn.RPN
var getInputInst getInput
var interruptCheckInst interruptCheck
var inputWin input.InputWindow
var screen ili948x.Ili948xScreen
var root window.WindowRoot
var fileOps fileops.FileOps
var fileOpsDriver tinyfs.FileOpsDriver

// We need a way to implement break, but it's complicated by the fact
// that reading the keyboard via i2c takes a long time compared to
// executing a command.
//
// I tried using a time.Time based check, but that had a ~35% performance
// hit as-measured by @benchmark.  So Instead, we go with a simple counter.
//
// The problem with a simple counter is that the delay function
// can make the counts happen very slowly, thus we need a hook to handle that
// case as well.
const breakExecCountTrigger = 8192

type interruptCheck struct {
	calls       uint32
	origSleepFn func(float64)
}

// Initialize this after the keyboard is ready to read, just in case.
func (ic *interruptCheck) Init() {
	ic.origSleepFn = functions.DelaySleepFn
	functions.DelaySleepFn = ic.delaySleepFn
	rpnInst.Interrupt = ic.checkForInterrupt
}

func (ic *interruptCheck) checkForInterrupt() bool {
	ic.calls++
	if ic.calls < breakExecCountTrigger {
		return false
	}
	ic.calls = 0
	k, _ := getInputInst.keyboard.GetChar()
	return k == key.KEY_BREAK
}

func (ic *interruptCheck) delaySleepFn(t float64) {
	ic.calls += uint32(t * float64(breakExecCountTrigger*4))
	ic.origSleepFn(t)
}

type getInput struct {
	lcd      *ili948x.Ili948xTxtW
	serial   serial.Serial
	keyboard i2ckbd.I2CKbd
}

func main() {
	screen.Init()
	// This only seeems to work for panics I throw and not errors
	// like array out of bounds.
	defer func() {
		if r := recover(); r != nil {
			handlePanic(screen.Device, r)
		}
	}()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	time.Sleep(200 * time.Millisecond)

	elog.Print("Started")
	rpnInst.Init(maxStackDepth)
	functions.RegisterAll(&rpnInst)
	_ = fileOpsDriver.Init()
	fileOps.InitAndRegister(&rpnInst, 65536, &fileOpsDriver)
	err := buildUI()
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
	_ = commands.InitWindowCommands(&rpnInst, &root, &screen, newPixelPlotWindow)
	_ = plotwin.InitPlotCommands(&rpnInst, &root, &screen)
	if err := startup.LCD320Startup(&rpnInst); err != nil {
		return err
	}
	interruptCheckInst.Init()
	w, h := screen.ScreenSize()
	if err := root.Update(&rpnInst, w, h, true); err != nil {
		return err
	}
	for {
		w, h = screen.ScreenSize()
		if err := root.Update(&rpnInst, w, h, true); err != nil {
			if errors.Is(err, input.ErrExit) {
				return nil
			}
			return err
		}
	}
}

func buildUI() error {
	w, h := screen.ScreenSize()
	root.Init(w, h)
	txtw, err := screen.NewTextWindow()
	if err != nil {
		return err
	}
	getInputInst.Init()
	inputWin.Init(&getInputInst, txtw, &rpnInst, scrollbytes)
	getInputInst.lcd = txtw.(*ili948x.Ili948xTxtW)
	root.AddWindowChildToRoot(&inputWin, "i", 100)
	return nil
}

func (gi *getInput) Init() {
	if err := gi.keyboard.Init(); err != nil {
		elog.Print("failed to init keyboard: ", err.Error())
	}
	gi.serial.Init()
}

func (gi *getInput) GetChar() (key.Key, error) {
	for {
		time.Sleep(20 * time.Millisecond)
		k := gi.serial.GetChar()
		if k != 0 {
			return k, nil
		}
		var err error
		k, err = gi.keyboard.GetChar()
		if err != nil {
			elog.Print("keyboard error: ", err.Error())
			time.Sleep(1 * time.Second)
		} else if k != 0 {
			return k, nil
		}
	}
}
