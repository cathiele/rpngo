package functions

import (
	"mattwach/rpngo/rpn"
	"time"
)

const DelayHelp = "Pauses for the given number of seconds"

func Delay(r *rpn.RPN) error {
	cv, err := r.PopComplex()
	if err != nil {
		return err
	}
	v := real(cv)
	if v <= 0 {
		return nil
	}
	time.Sleep(time.Duration(v*1000) * time.Millisecond)
	return nil
}

const TimeHelp = "Returns unix epoch time, assuming the clock on the hardware is calibrated."

func Time(r *rpn.RPN) error {
	t := time.Now()
	return r.PushComplex(complex(float64(t.UnixMicro())/1000000, 0))
}
