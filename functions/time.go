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
