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
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: a.Complex + b.Complex})
}

func Subtract(s *rpn.Stack) error {
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: a.Complex - b.Complex})
}

func Multiply(s *rpn.Stack) error {
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: a.Complex * b.Complex})
}

func Divide(s *rpn.Stack) error {
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	if b.Complex == 0 {
		return errDivideByZero
	}
	return s.Push(rpn.Frame{Complex: a.Complex / b.Complex})
}

func Square(s *rpn.Stack) error {
	a, err := s.Pop()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: a.Complex * a.Complex})
}

func SquareRoot(s *rpn.Stack) error {
	a, err := s.Pop()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: cmplx.Sqrt(a.Complex)})
}
