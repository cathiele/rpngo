// package functions defines core functions
package functions

import (
	"errors"
	"math/cmplx"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

var (
	errDivideByZero = errors.New("divide by zero")
)

const AddHelp = "Adds two numbers"

func Add(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushComplex(a + b)
}

const SubtractHelp = "Subtracts two numbers"

func Subtract(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushComplex(a - b)
}

const MultiplyHelp = "Multiplies two numbers"

func Multiply(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushComplex(a * b)
}

const DivideHelp = "Divides two numbers"

func Divide(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	if b == 0 {
		return errDivideByZero
	}
	return r.PushComplex(a / b)
}

const NegateHelp = "Negates the top number"

func Negate(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(-f.Complex)
	}
	if f.Type == rpn.BOOL_FRAME {
		if f.Int == 0 {
			return r.PushBool(true)
		}
		return r.PushBool(false)
	}
	return errors.New("expected number or boolean")
}

const SquareHelp = "executes v * v"

func Square(r *rpn.RPN) error {
	a, err := r.PopComplex()
	if err != nil {
		return err
	}
	return r.PushComplex(a * a)
}

const SquareRootHelp = "takes the square root of a complex number"

func SquareRoot(r *rpn.RPN) error {
	a, err := r.PopComplex()
	if err != nil {
		return err
	}
	return r.PushComplex(cmplx.Sqrt(a))
}

const ExecHelp = "Executes a string\n" +
	"Example: '4 5 +' @"

func Exec(r *rpn.RPN) error {
	s, err := r.PopString()
	if err != nil {
		return err
	}
	fields, err := parse.Fields(s)
	if err != nil {
		return err
	}
	return r.Exec(fields)
}
