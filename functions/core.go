// package functions defines core functions
package functions

import (
	"math"
	"math/rand"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const AddHelp = "Adds two numbers"

func Add(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false)+b.String(false), a.QuoteType()))
	}
	if b.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false)+b.String(false), b.QuoteType()))
	}
	if a.IsComplex() || b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrame(ac + bc))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab+bb, a))
}

const SubtractHelp = "Subtracts two numbers"

func Subtract(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(a, b)
	if err != nil {
		return err
	}
	if bothints {
		return r.PushFrame(rpn.IntFrameCloneType(a.UnsafeInt()-b.UnsafeInt(), a))
	}
	return r.PushFrame(rpn.ComplexFrame(a.UnsafeComplex() - b.UnsafeComplex()))
}

const MultiplyHelp = "Multiplies two numbers"

func Multiply(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(a, b)
	if err != nil {
		return err
	}
	if bothints {
		return r.PushFrame(rpn.IntFrameCloneType(a.UnsafeInt()*b.UnsafeInt(), a))
	}
	return r.PushFrame(rpn.ComplexFrame(a.UnsafeComplex() * b.UnsafeComplex()))
}

const DivideHelp = "Divides two numbers"

func Divide(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(a, b)
	if err != nil {
		return err
	}
	if bothints {
		bi := b.UnsafeInt()
		if bi == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushFrame(rpn.IntFrameCloneType(a.UnsafeInt()/bi, a))
	}
	bc := b.UnsafeComplex()
	if bc == 0 {
		return rpn.ErrDivideByZero
	}
	return r.PushFrame(rpn.ComplexFrame(a.UnsafeComplex() / bc))
}

const NegateHelp = "Negates the top number"

func Negate(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsComplex() {
		c, _ := f.Complex()
		return r.PushFrame(rpn.ComplexFrame(-c))
	}
	if f.IsBool() {
		b, _ := f.Bool()
		return r.PushFrame(rpn.BoolFrame(!b))
	}
	if f.IsInt() {
		i, _ := f.Int()
		return r.PushFrame(rpn.IntFrameCloneType(-i, f))
	}
	return rpn.ErrIllegalValue
}

const ExecHelp = "Executes a string\n" +
	"Example: '4 5 +' @"

func Exec(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	return parse.Fields(f.UnsafeString(), r.Exec)
}

const RandHelp = "Pushes a random number between 0 and 1"

func Rand(r *rpn.RPN) error {
	return r.PushFrame(rpn.RealFrame(rand.Float64()))
}

const RealHelp = "Takes the real portion of a complex number"

func Real(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	c, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.RealFrame(real(c)))
}

const ImagHelp = "Takes the imaginary portion of a complex number"

func Imag(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	c, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(complex(0, imag(c))))
}

const TrueHelp = "Pushes a boolean true"

func True(r *rpn.RPN) error {
	return r.PushFrame(rpn.BoolFrame(true))
}

const FalseHelp = "Pushes a boolean false"

func False(r *rpn.RPN) error {
	return r.PushFrame(rpn.BoolFrame(false))
}

const RoundHelp = "Rounds a number to the given number of places"

func Round(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	if (b < 0) || (b > 16) {
		return rpn.ErrIllegalValue
	}
	rl := real(a)
	im := imag(a)
	for i := 0; i < int(b); i++ {
		rl *= 10
		im *= 10
	}
	rl = math.Round(rl)
	im = math.Round(im)
	for i := 0; i < int(b); i++ {
		rl /= 10
		im /= 10
	}
	return r.PushFrame(rpn.ComplexFrame(complex(rl, im)))
}
