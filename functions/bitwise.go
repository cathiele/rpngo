package functions

import (
	"mattwach/rpngo/rpn"
)

const AndHelp = "Performs a logical AND operation"

func And(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a & b })
}

const OrHelp = "Performs a logical OR operation"

func Or(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a | b })
}

const XOrHelp = "Performs a logical XOR operation"

func XOr(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a ^ b })
}

const ShiftLeftHelp = "Performs a logical shift left operation"

func ShiftLeft(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a << b })
}

const ShiftRightHelp = "Performs a logical shift right operation"

func ShiftRight(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a >> b })
}

func binaryOp(r *rpn.RPN, fn func(a, b int64) int64) error {
	af, bf, err := r.Pop2Frames()
	if af.IsBool() {
		return binaryBoolOp(r, fn, af, bf)
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(fn(a, b), af))
}

func binaryBoolOp(r *rpn.RPN, fn func(a, b int64) int64, af, bf rpn.Frame) error {
	ab, err := af.Bool()
	if err != nil {
		return err
	}
	bb, err := bf.Bool()
	if err != nil {
		return err
	}
	var a int64 = 0
	var b int64 = 0
	if ab {
		a = 1
	}
	if bb {
		b = 1
	}
	v := fn(a, b)
	if v != 0 {
		return r.PushFrame(rpn.BoolFrame(true))
	}
	return r.PushFrame(rpn.BoolFrame(false))
}
