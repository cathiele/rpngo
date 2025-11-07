package functions

import (
	"mattwach/rpngo/rpn"
	"strconv"
)

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
	i := len(r.Frames)
	for _, f := range r.Frames {
		i--
		r.Print(strconv.Itoa(i) + ": " + f.String(true) + "\n")
	}
	return nil
}

const InputHelp = "Pauses for user input and pushes the result to the stack as a string"

func Input(r *rpn.RPN) error {
	str, err := r.Input(r)
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(str, rpn.STRING_DOUBLE_QUOTE))
}
