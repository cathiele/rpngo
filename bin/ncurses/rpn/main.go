// A simple console demonstration
package main

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/io/drivers/curses"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/io/window/input"
	"mattwach/rpngo/io/window/stackwin"
	"mattwach/rpngo/rpn"
	"os"
)

func run() error {
	os.RemoveAll("/tmp/rpngo.log")
	logFile, err := os.OpenFile("/tmp/rpngo.log", os.O_CREATE|os.O_WRONLY, 0644)
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
	for _, arg := range os.Args[1:] {
		if err := r.Exec(arg); err != nil {
			return err
		}
	}

	r.Stack.IterFrames(func(sf rpn.Frame) {
		fmt.Println(sf.String())
	})

	return nil
}

func interactive(r *rpn.RPN) error {
	screen, err := curses.Init()
	if err != nil {
		return err
	}
	defer screen.End()
	root := window.NewWindowGroup(true)
	root.SetVertical(true)
	w, h := screen.Size()
	root.Resize(0, 0, w, h)
	txtw, err := screen.NewTextWindow(0, 0, w, h)
	if err != nil {
		return err
	}
	iw, err := input.Init(txtw.(*curses.Curses), txtw)
	if err != nil {
		return err
	}
	root.AddWindowChild(iw, "i", 100)
	stackw, err := screen.NewTextWindow(50, 0, w-50, h)
	if err != nil {
		return err
	}
	sw, err := stackwin.Init(stackw)
	if err != nil {
		return err
	}
	root.AddWindowChild(sw, "s1", 25)
	for {
		if err := root.Update(r); err != nil {
			if errors.Is(err, input.ErrExit) {
				return nil
			}
			return err
		}
		screen.Refresh()
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
