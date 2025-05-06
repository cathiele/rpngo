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
	for _, arg := range os.Args[1:] {
		if err := r.Exec(arg); err != nil {
			return err
		}
	}

	r.Stack.IterFrames(func(sf rpn.Frame) {
		fmt.Printf("%f\n", sf.Float)
	})

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
