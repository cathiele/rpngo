package main

import (
	"machine"
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
	// avoid using the UART by-default becuase it has a 15% perf penalty
	gi.serial.Init(nil)
	xmodemCommands.InitAndRegister(r, nil)
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
	if gi.serial.Serial != nil {
		for _, c := range str {
			_ = gi.serial.WriteByte(byte(c))
		}
	}
	gi.originalPrint(str)
}

const SerialHelp = "If true, enables serial communications with a host PC."

func (gi *picoCalcIO) Serial(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	enabled, err := f.Bool()
	if err != nil {
		return err
	}
	if !enabled {
		if gi.serial.Serial != nil {
			return rpn.ErrNotSupported
		}
		return nil
	}
	if gi.serial.Serial != nil {
		return nil
	}
	gi.serial.Init(machine.Serial)
	gi.originalPrint = r.Print
	r.Print = gi.Print
	xmodemCommands.SetSerial(&gi.serial)
	return nil
}
