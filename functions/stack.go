package functions

import (
	"fmt"
	"mattwach/rpngo/rpn"
)

const NoOpHelp = "No operation. e.g. 'noop' plot will plot y = x"

func NoOp(r *rpn.RPN) error {
	return nil
}

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

const PrintHelp = "Prints the head element of the stack to the output window"

func Print(r *rpn.RPN) error {
	f, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	return nil
}

const PrintXHelp = "Pops head element of the stack and prints it"

func PrintX(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	return nil
}

const PrintSHelp = "Prints the head element of the stack plus a space"

func PrintS(r *rpn.RPN) error {
	f, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print(" ")
	return nil
}

const PrintSXHelp = "Pops head element of the stack and prints it and a space"

func PrintSX(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print(" ")
	return nil
}

const PrintlnHelp = "Prints the head element of the stack plus a newline"

func Println(r *rpn.RPN) error {
	f, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print("\n")
	return nil
}

const PrintlnXHelp = "Pops head element of the stack and prints it and a newline"

func PrintlnX(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print("\n")
	return nil
}

const PrintAllHelp = "Prints the whole stack"

func PrintAll(r *rpn.RPN) error {
	i := r.Size()
	r.IterFrames(func(f rpn.Frame) {
		i--
		r.Print(fmt.Sprintf("%d: %s\n", i, f.String(true)))
	})
	return nil
}

const InputHelp = "Pauses for user input and pushes the result to the stack as a string"

func Input(r *rpn.RPN) error {
	str, err := r.Input(r)
	if err != nil {
		return err
	}
	return r.PushString(str)
}
