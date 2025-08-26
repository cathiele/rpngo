package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

func Sin(s *rpn.Stack) error {
	a, err := s.Pop()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: cmplx.Sin(a.Complex)})
}

func Cos(s *rpn.Stack) error {
	a, err := s.Pop()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: cmplx.Cos(a.Complex)})
}

func Tan(s *rpn.Stack) error {
	a, err := s.Pop()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: cmplx.Tan(a.Complex)})
}
