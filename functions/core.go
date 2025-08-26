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

// Add adds 2 complex numbers
func Add(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a + b)
}

// Subtract subtracts 2 complex numbers
func Subtract(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a - b)
}

// Multiply multiplies 2 complex numbers
func Multiply(s *rpn.Stack) error {
	a, b, err := s.Pop2Complex()
	if err != nil {
		return err
	}
	return s.PushComplex(a * b)
}

// Divide divides 2 complex numbers
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

// Square executes v * v
func Square(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(a * a)
}

// SquareRoot takes the square root of a complex number
func SquareRoot(s *rpn.Stack) error {
	a, err := s.PopComplex()
	if err != nil {
		return err
	}
	return s.PushComplex(cmplx.Sqrt(a))
}
