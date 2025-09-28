// A simple console demonstration
package main

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/drivers/curses"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/startup"
	"mattwach/rpngo/window"
	"mattwach/rpngo/window/commands"
	"mattwach/rpngo/window/input"
	"mattwach/rpngo/window/plotwin"
	"os"
	"os/signal"
)

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
	r.Init()
	functions.RegisterAll(&r)

	if len(os.Args) > 1 {
		return cli(&r)
	}

	return interactive(&r)
}

func cli(r *rpn.RPN) error {
	if err := r.Exec(os.Args[1:]); err != nil {
		return err
	}

	r.IterFrames(func(sf rpn.Frame) {
		fmt.Println(sf.String(true))
	})

	return nil
}

func interactive(r *rpn.RPN) error {
	r.Interrupt = setupSignals()
	screen, err := curses.Init()
	if err != nil {
		return err
	}
	defer screen.End()
	root, err := buildUI(screen, r)
	if err != nil {
		return err
	}
	_ = commands.InitWindowCommands(r, root, screen)
	_ = plotwin.InitPlotCommands(r, root, screen)
	if err := startup.OSStartup(r); err != nil {
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

func setupSignals() chan bool {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	interrupt := make(chan bool, 1)
	go func() {
		for {
			<-sigc
			interrupt <- true
		}
	}()
	return interrupt

}

func buildUI(screen *curses.Curses, r *rpn.RPN) (*window.WindowRoot, error) {
	w, h := screen.ScreenSize()
	root := window.NewWindowRoot(w, h)
	if err := addInputWindow(screen, root, r); err != nil {
		return nil, err
	}
	return root, nil
}

func addInputWindow(screen window.Screen, root *window.WindowRoot, r *rpn.RPN) error {
	txtw, err := screen.NewTextWindow(0, 0, 10, 5)
	if err != nil {
		return err
	}
	iw, err := input.Init(txtw.(*curses.Curses), txtw, r)
	if err != nil {
		return err
	}
	root.AddWindowChild(iw, "i", 100)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
