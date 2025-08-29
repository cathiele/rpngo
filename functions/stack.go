package functions

import "mattwach/rpngo/rpn"

const CopyHelp = "Duplicates the element at the head of the stack"

func Copy(r *rpn.RPN) error {
	a, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	return r.PushFrame(a)
}

const CopyIndexHelp = "Duplicates the element at index n on the stack\n" +
	"Example: 11 22 1 ci # will copy 11"

func CopyIndex(r *rpn.RPN) error {
	i, err := r.PopStackIndex()
	if err != nil {
		return err
	}
	a, err := r.PeekFrame(i)
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

const MoveIndexHelp = "Moves the element at index n on the stack\n" +
	"Example: 11 22 1 mi # same result as s"

func MoveIndex(r *rpn.RPN) error {
	i, err := r.PopStackIndex()
	if err != nil {
		return err
	}
	a, err := r.DeleteFrame(i)
	if err != nil {
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

const DropAllHelp = "Drops the element at the top of the stack"

func DropAll(r *rpn.RPN) error {
	r.Clear()
	return nil
}

const DropIndexHelp = "Removes the element at index n on the stack\n" +
	"Example: 11 22 1 di # removes the 11"

func DropIndex(r *rpn.RPN) error {
	i, err := r.PopStackIndex()
	if err != nil {
		return err
	}
	_, err = r.DeleteFrame(i)
	if err != nil {
		return err
	}
	return nil
}
