// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"log"
	"machine"
	"mattwach/rpngo/functions"
	"mattwach/rpngo/io/drivers/tinygo/ili9341tw"
	"mattwach/rpngo/io/key"
	"mattwach/rpngo/io/window"
	"mattwach/rpngo/io/window/input"
	"mattwach/rpngo/rpn"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	time.Sleep(2 * time.Second)

	log.SetOutput(os.Stdout)
	log.Println("Started")
	var r rpn.RPN
	r.Init()
	functions.RegisterAll(&r)

	var screen ili9341tw.Ili9341Screen
	screen.Init()
	input, err := buildUI(&screen, &r)
	if err != nil {
		return err
	}

	for {
		if err := input.Update(&r); err != nil {
			return err
		}
	}
}

func buildUI(screen window.Screen, r *rpn.RPN) (*input.InputWindow, error) {
	w, h := screen.Size()
	txtw, err := screen.NewTextWindow(0, 0, w, h)
	if err != nil {
		return nil, err
	}
	iw, err := input.Init(getInput{}, txtw, r)
	if err != nil {
		return nil, err
	}
	return iw, nil
}

type getInput struct{}

func (getInput) GetChar() (key.Key, error) {
	for {
		c, err := machine.Serial.ReadByte()
		machine.Serial.WriteByte(c)
		if err == nil {
			return key.Key(c), nil
		}

		time.Sleep(time.Millisecond * 10)
	}
}
