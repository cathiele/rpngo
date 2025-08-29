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

const SwapHelp = "Swaps two elements at the top of the stack"

func Swap(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if err := r.PushFrame(b); err != nil {
		return err
	}
	return r.PushFrame(a)
}

const DropHelp = "Drops the element at the top of the stack"

func Drop(r *rpn.RPN) error {
	_, err := r.PopFrame()
	if err != nil {
		return err
	}
	return nil
}
