// A simple console demonstration
package main

import (
	"fmt"
	"log"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/io/drivers/curses"
	"mattwach/rpngo/io/input"
	"mattwach/rpngo/io/window"
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
	root := window.NewWindowGroup()
	w, h := screen.Size()
	root.Resize(0, 0, w, h)
	inw := screen.NewTextWindow(0, 0, w, h)
	root.AddTextWindowChild(inw, "i", 100)
	return input.Loop(r, inw.(*curses.Curses), root, screen)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
