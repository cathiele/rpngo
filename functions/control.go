package functions

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const IfHelp = "Pops action, then val. Keeps val if cond is true.\n" +
	"Example: 4 3 > '\"a is greater than b\" printlnx' if @"

func If(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	cond, err := r.PopBool()
	if err != nil {
		r.PushFrame(f)
		return err
	}
	if cond {
		return r.PushFrame(f)
	}
	return nil
}

const IfElseHelp = "Pops elsev, ifv, cond. Pushes ifv if cond=true, else pushes elsev.\n" +
	"Example: 4 3 > 'a is greater than b' 'a is less than b' ifelse printlnx"

func IfElse(r *rpn.RPN) error {
	elsev, err := r.PopFrame()
	if err != nil {
		return err
	}
	ifv, err := r.PopFrame()
	if err != nil {
		r.PushFrame(elsev)
		return err
	}
	cond, err := r.PopBool()
	if err != nil {
		r.PushFrame(ifv)
		r.PushFrame(elsev)
		return err
	}
	if cond {
		return r.PushFrame(ifv)
	}
	return r.PushFrame(elsev)
}

const ForHelp = "Executes the head of the stack in a loop until a value < is found\n" +
	"Example: 1 'c 1 + c 50 <' for # put 1 to 50 on the stack"

func For(r *rpn.RPN) error {
	macro, err := r.PopString()
	if err != nil {
		return err
	}
	fields, err := parse.Fields(macro)
	if err != nil {
		r.PushString(macro)
		return err
	}
	for {
		if err := r.Exec(fields); err != nil {
			return err
		}
		cond, err := r.PopBool()
		if err != nil {
			return err
		}
		if !cond {
			break
		}
	}
	return nil
}
