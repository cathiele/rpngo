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
	cf, err := r.PopFrame()
	if err != nil {
		return err
	}
	cond, err := cf.Bool()
	if err != nil {
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
		return err
	}
	cf, err := r.PopFrame()
	if err != nil {
		return err
	}
	cond, err := cf.Bool()
	if err != nil {
		return err
	}
	if cond {
		return r.PushFrame(ifv)
	}
	return r.PushFrame(elsev)
}

const ForHelp = "Executes the head of the stack in a loop until a value < is found\n" +
	"Example: 1 'c 1 + c 50 <' for # put 1 to 50 on the stack"

// avoid allocating fields on the stack (for can be nested though)
var forFields []string
var forFieldsStart []int

func For(r *rpn.RPN) error {
	forFieldsStart = append(forFieldsStart, len(forFields))
	defer func() { forFieldsStart = forFieldsStart[:len(forFieldsStart)-1] }()
	mf, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !mf.IsString() {
		return rpn.ErrExpectedAString
	}
	macro := mf.UnsafeString()
	addField := func(t string) error {
		forFields = append(forFields, t)
		return nil
	}
	if err := parse.Fields(macro, addField); err != nil {
		return err
	}
	for {
		if err := r.ExecSlice(forFields[forFieldsStart[len(forFieldsStart)-1]:]); err != nil {
			return err
		}
		cf, err := r.PopFrame()
		if err != nil {
			return err
		}
		cond, err := cf.Bool()
		if err != nil {
			return err
		}
		if !cond {
			break
		}
	}
	return nil
}
