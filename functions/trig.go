package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

// Sin takes the sine of a complex number
func Sin(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Sin(a))
}

// Cos takes the cosine of a complex number
func Cos(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Cos(a))
}

// Tan takes the tangent of a complex number
func Tan(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Tan(a))
}
