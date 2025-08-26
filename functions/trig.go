package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

func Sin(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Sin(a))
}

func Cos(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Cos(a))
}

func Tan(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Tan(a))
}
