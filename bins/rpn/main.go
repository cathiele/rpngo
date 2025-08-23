// A simple console demonstration
package main

import (
	"fmt"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/io"
	"mattwach/rpngo/io/curses"
	"mattwach/rpngo/rpn"
	"os"
)

func run() error {
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
	c, err := curses.Init()
	if err != nil {
		return err
	}
	defer c.End()
	return io.Loop(r, c, c)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
