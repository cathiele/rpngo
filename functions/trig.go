package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

const SinHelp = "takes the sine of a number"

func Sin(r *rpn.RPN) error {
	a, err := r.PopComplex()
	if err != nil {
		return err
	}
	return r.PushComplex(cmplx.Sin(a))
}

const CosHelp = "takes the cosine of a number"

func Cos(r *rpn.RPN) error {
	a, err := r.PopComplex()
	if err != nil {
		return err
	}
	return r.PushComplex(cmplx.Cos(a))
}

const TanHelp = "takes the tangent of a number"

func Tan(r *rpn.RPN) error {
	a, err := r.PopComplex()
	if err != nil {
		return err
	}
	return r.PushComplex(cmplx.Tan(a))
}
