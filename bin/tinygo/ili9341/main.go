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
	gi := &getInput{}
	iw, err := input.Init(gi, txtw, r)
	gi.lcd = txtw.(*ili9341tw.Ili9341TW)
	if err != nil {
		return nil, err
	}
	return iw, nil
}

type getInput struct {
	frame uint8
	lcd   *ili9341tw.Ili9341TW
}

func (g *getInput) GetChar() (key.Key, error) {
	for {
		c, err := machine.Serial.ReadByte()
		machine.Serial.WriteByte(c)
		if err == nil {
			if c == 13 {
				return '\n', nil
			}
			return key.Key(c), nil
		}

		time.Sleep(time.Millisecond * 10)
		g.frame++
		g.lcd.ShowCursorIfEnabled((g.frame & 0xC0) != 0)
	}
}
