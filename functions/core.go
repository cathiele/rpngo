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

const AddHelp = "Adds two numbers"

func Add(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a + b)
}

const SubtractHelp = "Subtracts two numbers"

func Subtract(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a - b)
}

const MultiplyHelp = "Multiplies two numbers"

func Multiply(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a * b)
}

const DivideHelp = "Divides two numbers"

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

const SquareHelp = "executes v * v"

func Square(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(a * a)
}

const SquareRootHelp = "takes the square root of a complex number"

func SquareRoot(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Sqrt(a))
}
