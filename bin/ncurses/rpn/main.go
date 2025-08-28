// A simple console demonstration
package main

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/io"
	"mattwach/rpngo/io/drivers/curses"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/io/window/commands"
	"mattwach/rpngo/io/window/input"
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
	root, err := buildUI(screen)
	if err != nil {
		return err
	}
	wc := commands.InitWindowCommands(root, screen)
	wc.Register(r)
	if err := io.OSStartup(r); err != nil {
		return err
	}
	if err := root.Update(r, false); err != nil {
		return err
	}
	for {
		if err := root.Update(r, true); err != nil {
			if errors.Is(err, input.ErrExit) {
				return nil
			}
			return err
		}
	}
}

func buildUI(screen *curses.Curses) (*window.WindowGroup, error) {
	root := window.NewWindowGroup(true)
	w, h := screen.Size()
	if err := root.Resize(0, 0, w, h); err != nil {
		return nil, err
	}

	if err := addInputWindow(screen, root); err != nil {
		return nil, err
	}
	return root, nil
}

func addInputWindow(screen window.Screen, root *window.WindowGroup) error {
	txtw, err := screen.NewTextWindow(0, 0, 10, 5)
	if err != nil {
		return err
	}
	iw, err := input.Init(txtw.(*curses.Curses), txtw)
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
