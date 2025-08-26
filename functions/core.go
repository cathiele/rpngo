// package functions defines core functions
package functions

import (
	"errors"
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

var (
	errDivideByZero = errors.New("divide by zero")
)

func Add(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a + b)
}

func Subtract(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a - b)
}

func Multiply(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a * b)
}

func Divide(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	if b == 0 {
		return errDivideByZero
	}
	return s.PushComplex(a / b)
}

func Square(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(a * a)
}

func SquareRoot(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Sqrt(a))
}
