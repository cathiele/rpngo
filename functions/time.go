package functions

import (
	"mattwach/rpngo/rpn"
	"time"
)

const DelayHelp = "Pauses for the given number of seconds"

var DelaySleepFn = func(t float64) {
	time.Sleep(time.Duration(t*1000) * time.Millisecond)
}

func Delay(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	v, err := f.Real()
	if err != nil {
		return err
	}
	if v <= 0 {
		return nil
	}
	DelaySleepFn(v)
	return nil
}

const TimeHelp = "Returns unix epoch time, assuming the clock on the hardware is calibrated."

func Time(r *rpn.RPN) error {
	t := time.Now()
	return r.PushFrame(rpn.RealFrame(float64(t.UnixMicro()) / 1000000))
}
