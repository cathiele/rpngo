package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

const SinHelp = "takes the sine of a number"

func Sin(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Sin(a)))
}

const CosHelp = "takes the cosine of a number"

func Cos(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Cos(a)))
}

const TanHelp = "takes the tangent of a number"

func Tan(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Tan(a)))
}
