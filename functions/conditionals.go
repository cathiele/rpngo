package functions

import (
	"mattwach/rpngo/rpn"
)

const GreaterThanHelp = "Returns true if a > b, false otherwise"

func GreaterThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!af.IsLessThanOrEqual(bf)))
}

const GreaterThanEqualHelp = "Returns true if a >= b, false otherwise"

func GreaterThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!af.IsLessThan(bf)))
}

const LessThanHelp = "Returns true if a < b, false otherwise"

func LessThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(af.IsLessThan(bf)))
}

const LessThanEqualHelp = "Returns true if a <= b, false otherwise"

func LessThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(af.IsLessThanOrEqual(bf)))
}

const EqualHelp = "Returns true if a = b, false otherwise (approximate)"

func Equal(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(af.IsEqual(bf)))
}

const NotEqualHelp = "Returns true if a != b, false otherwise (approximate)"

func NotEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!af.IsEqual(bf)))
}

const MinHelp = "Pops two frames and repushes the minimum value.  Pushes $1 if the frames were equal."

func Min(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if af.IsLessThanOrEqual(bf) {
		return r.PushFrame(af)
	}
	return r.PushFrame(bf)
}

const MaxHelp = "Pops two frames and repushes the maximum value.  Pushes $1 if the frames were equal."

func Max(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if af.IsLessThan(bf) {
		return r.PushFrame(bf)
	}
	return r.PushFrame(af)
}
