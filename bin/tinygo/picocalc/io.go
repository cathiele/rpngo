package main

import (
	"mattwach/rpngo/drivers/tinygo/picocalc/i2ckbd"
	"mattwach/rpngo/drivers/tinygo/picocalc/ili948x"
	"mattwach/rpngo/drivers/tinygo/serial"
	"mattwach/rpngo/elog"
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"time"
)

type picoCalcIO struct {
	serial        serial.Serial
	keyboard      i2ckbd.I2CKbd
	screen        ili948x.Ili948xScreen
	originalPrint func(string)
}

func (gi *picoCalcIO) Init() {
	gi.screen.Init()
	if err := gi.keyboard.Init(); err != nil {
		elog.Print("failed to init keyboard: ", err.Error())
	}
	// avoid using the UART by-default becuase it has a 15% perf penalty
	gi.serial.Init(nil)
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

func (gi *picoCalcIO) PatchInUARTPrint(r *rpn.RPN) {
	gi.originalPrint = r.Print
	r.Print = gi.Print
}

func (gi *picoCalcIO) Print(str string) {
	if gi.serial.Serial != nil {
		for _, c := range str {
			_ = gi.serial.WriteByte(byte(c))
		}
	}
	gi.originalPrint(str)
}
