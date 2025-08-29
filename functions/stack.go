package functions

import "mattwach/rpngo/rpn"

const DuplicateHelp = "Duplicates the element at the head of the stack"

func Duplicate(r *rpn.RPN) error {
	a, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	return r.PushFrame(a)
}
