// A simple console demonstration
package main

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/drivers/curses"
	"mattwach/rpngo/drivers/posix/fs"
	"mattwach/rpngo/drivers/posix/serial"
	"mattwach/rpngo/fileops"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/startup"
	"mattwach/rpngo/window"
	"mattwach/rpngo/window/commands"
	"mattwach/rpngo/window/input"
	"mattwach/rpngo/window/plotwin"
	"mattwach/rpngo/xmodem"
	"os"
	"os/signal"
)

const scrollbytes = 256 * 1024
const maxStackDepth = 65536

func run() error {
	os.RemoveAll("/tmp/rpngo.log")
	logFile, err := os.Create("/tmp/rpngo.log")
	if err != nil {
		return err
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Application started")
	var r rpn.RPN
	r.Init(maxStackDepth)
	functions.RegisterAll(&r)
	var fo fileops.FileOps
	fo.InitAndRegister(&r, 65536, &fs.FileOpsDriver{})
	var xm xmodem.XmodemCommands
	xm.InitAndRegister(&r, &serial.Serial{})

	if len(os.Args) > 1 {
		return cli(&r)
	}

	return interactive(&r)
}

func cli(r *rpn.RPN) error {
	if err := r.ExecSlice(os.Args[1:]); err != nil {
		return err
	}

	for _, f := range r.Frames {
		fmt.Println(f.String(true))
	}

	return nil
}

func interactive(r *rpn.RPN) error {
	var inter interrupt
	inter.init()
	r.Interrupt = inter.interrupt
	screen, err := curses.Init()
	if err != nil {
		return err
	}
	defer screen.End()
	var root window.WindowRoot
	err = buildUI(&root, screen, r)
	if err != nil {
		return err
	}
	newTextPlotWindow := func() (window.WindowWithProps, error) {
		var tpw plotwin.TxtPlotWindow
		pw, err := screen.NewTextWindow()
		if err != nil {
			return nil, err
		}
		tpw.Init(pw)
		return &tpw, nil
	}
	_ = commands.InitWindowCommands(r, &root, screen, newTextPlotWindow)
	_ = plotwin.InitPlotCommands(r, &root, screen)
	if err := startup.Startup(r, &fs.FileOpsDriver{}); err != nil {
		return err
	}
	w, h := screen.ScreenSize()
	if err := root.Update(r, w, h, true); err != nil {
		return err
	}
	for {
		w, h = screen.ScreenSize()
		if err := root.Update(r, w, h, true); err != nil {
			if errors.Is(err, input.ErrExit) {
				return nil
			}
			return err
		}
	}
}

type interrupt struct {
	sigc chan os.Signal
}

func (i *interrupt) init() {
	i.sigc = make(chan os.Signal, 1)
	signal.Notify(i.sigc, os.Interrupt)
}

func (i *interrupt) interrupt() bool {
	select {
	case <-i.sigc:
		return true
	default:
		return false
	}
}

func buildUI(root *window.WindowRoot, screen *curses.Curses, r *rpn.RPN) error {
	w, h := screen.ScreenSize()
	root.Init(w, h)
	if err := addInputWindow(screen, root, r); err != nil {
		return err
	}
	return nil
}

func addInputWindow(screen window.Screen, root *window.WindowRoot, r *rpn.RPN) error {
	txtw, err := screen.NewTextWindow()
	if err != nil {
		return err
	}
	var iw input.InputWindow
	iw.Init(txtw.(*curses.Curses), txtw, r, &fs.FileOpsDriver{}, scrollbytes)
	if err != nil {
		return err
	}
	root.AddWindowChildToRoot(&iw, "i", 100)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
