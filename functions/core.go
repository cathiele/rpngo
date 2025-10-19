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
	if a.IsString() || b.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false) + b.String(false)))
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
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		ac, _ := a.Complex()
		bc, _ := a.Complex()
		return r.PushFrame(rpn.ComplexFrame(ac - bc))
	}
	ai, _ := a.Int()
	bi, _ := b.Int()
	return r.PushFrame(rpn.IntFrameCloneType(ai-bi, a))
}

const MultiplyHelp = "Multiplies two numbers"

func Multiply(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		ac, _ := a.Complex()
		bc, _ := a.Complex()
		return r.PushFrame(rpn.ComplexFrame(ac * bc))
	}
	ai, _ := a.Int()
	bi, _ := b.Int()
	return r.PushFrame(rpn.IntFrameCloneType(ai*bi, a))
}

const DivideHelp = "Divides two numbers"

func Divide(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		ac, _ := a.Complex()
		bc, _ := a.Complex()
		if bc == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushFrame(rpn.ComplexFrame(ac / bc))
	}
	ai, _ := a.Int()
	bi, _ := b.Int()
	if bi == 0 {
		return rpn.ErrDivideByZero
	}
	return r.PushFrame(rpn.IntFrameCloneType(ai/bi, a))
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
	fields := make([]string, 32) // object allocated on the heap: escapes at line 124 (OK)
	fields, err = parse.Fields(f.String(false), fields)
	if err != nil {
		return err
	}
	return r.Exec(fields)
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
	return r.PushFrame(rpn.RealFrame(imag(c)))
}

const TrueHelp = "Pushes a boolean true"

func True(r *rpn.RPN) error {
	return r.PushBool(true)
}

const FalseHelp = "Pushes a boolean false"

func False(r *rpn.RPN) error {
	return r.PushBool(false)
}

const RoundHelp = "Rounds a number to the given number of places"

func Round(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	origb := b
	if err != nil {
		return err
	}
	if a.Type != rpn.COMPLEX_FRAME {
		a.Type = rpn.COMPLEX_FRAME
		a.Complex = complex(float64(a.Int), 0)
	}
	if b.Type == rpn.COMPLEX_FRAME {
		if imag(b.Complex) != 0 {
			r.PushFrame(a)
			r.PushFrame(origb)
			return rpn.ErrComplexNumberNotSupported
		}
		if real(b.Complex) != math.Round(real(b.Complex)) {
			r.PushFrame(a)
			r.PushFrame(origb)
			return rpn.ErrIllegalValue
		}
		b.Int = int64(real(b.Complex))
	}
	if (b.Int < 0) || (b.Int > 16) {
		r.PushFrame(a)
		r.PushFrame(origb)
		return rpn.ErrIllegalValue
	}
	rl := real(a.Complex)
	im := imag(a.Complex)
	for i := 0; i < int(b.Int); i++ {
		rl *= 10
		im *= 10
	}
	rl = math.Round(rl)
	im = math.Round(im)
	for i := 0; i < int(b.Int); i++ {
		rl /= 10
		im /= 10
	}
	return r.PushComplex(complex(rl, im))
}
