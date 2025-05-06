// package functions defines core functions
package functions

import (
	"errors"
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
	return s.Push(rpn.Frame{Float: a.Float + b.Float})
}

func Subtract(s *rpn.Stack) error {
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Float: a.Float - b.Float})
}

func Multiply(s *rpn.Stack) error {
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Float: a.Float * b.Float})
}

func Divide(s *rpn.Stack) error {
	a, b, err := s.Pop2()
	if err != nil {
		return err
	}
	if b.Float == 0 {
		return errDivideByZero
	}
	return s.Push(rpn.Frame{Float: a.Float / b.Float})
}
