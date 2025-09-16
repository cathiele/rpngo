// A simple console demonstration
package main

import (
	"fmt"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/rpn"
	"os"
)

func run() error {
	var r rpn.RPN
	r.Init()
	functions.RegisterAll(&r)

	if err := r.Exec(os.Args[1:]); err != nil {
		return err
	}

	r.IterFrames(func(sf rpn.Frame) {
		fmt.Println(sf.String(true))
	})

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
