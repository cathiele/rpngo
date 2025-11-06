package functions

import "mattwach/rpngo/rpn"

const ReverseHelp = "Reverses the stack"

func Reverse(r *rpn.RPN) error {
	if len(r.Frames) <= 1 {
		return nil
	}
	mid := len(r.Frames) / 2
	for i := 0; i < mid; i++ {
		end := len(r.Frames) - i - 1
		r.Frames[i], r.Frames[end] = r.Frames[end], r.Frames[i]
	}
	return nil
}
