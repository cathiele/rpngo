package main

import (
	"mattwach/rpngo/drivers/tinygo/picocalc/i2ckbd"
	"mattwach/rpngo/drivers/tinygo/picocalc/ili948x"
	"mattwach/rpngo/drivers/tinygo/serial"
	"mattwach/rpngo/elog"
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/xmodem"
	"time"
)

var xmodemCommands xmodem.XmodemCommands

type picoCalcIO struct {
	serial        serial.Serial
	keyboard      i2ckbd.I2CKbd
	screen        ili948x.Ili948xScreen
	originalPrint func(string)
}

func (gi *picoCalcIO) Init(r *rpn.RPN) {
	gi.screen.Init()
	if err := gi.keyboard.Init(); err != nil {
		elog.Print("failed to init keyboard: ", err.Error())
	}
	gi.serial.Init(true)
	xmodemCommands.InitAndRegister(r, &gi.serial)
	rpnInst.Register("serial", gi.serialFn, rpn.CatIO, serialHelp)
}

func (gi *picoCalcIO) RegisterPrint(r *rpn.RPN) {
	gi.originalPrint = r.Print
	r.Print = gi.Print
}

func (gi *picoCalcIO) GetChar() (key.Key, error) {
	for {
		time.Sleep(20 * time.Millisecond)
		k := gi.serial.GetChar()
		if k != 0 {
			return k, nil
		}
		k, _ = gi.keyboard.GetChar()
		if k != 0 {
			return k, nil
		}
	}
}

func (gi *picoCalcIO) Print(str string) {
	if gi.serial.Enabled {
		for _, c := range str {
			_ = gi.serial.WriteByte(byte(c))
		}
	}
	gi.originalPrint(str)
}

const serialHelp = "If true, enables serial communications with a host PC."

func (gi *picoCalcIO) serialFn(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	enabled, err := f.Bool()
	if err != nil {
		return err
	}
	gi.serial.Enabled = enabled
	return nil
}
