package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

const SinHelp = "takes the sine of a number"

func Sin(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Sin(a))
}

const CosHelp = "takes the cosine of a number"

func Cos(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Cos(a))
}

const TanHelp = "takes the tangent of a number"

func Tan(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Tan(a))
}
