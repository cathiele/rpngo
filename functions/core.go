// package functions defines core functions
package functions

import (
	"math"
	"math/cmplx"
	"math/rand"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const NoOpHelp = "No operation. e.g. 'noop' plot will plot y = x"

func NoOp(r *rpn.RPN) error {
	return nil
}

const AddHelp = "Adds two numbers"

func Add(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false)+b.String(false), a.Type()))
	}
	if b.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false)+b.String(false), b.Type()))
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()+bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac+b.UnsafeComplex(), b.Type()))
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
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()-bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac-b.UnsafeComplex(), b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab-bb, a))
}

const MultiplyHelp = "Multiplies two numbers"

func Multiply(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()*bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac*b.UnsafeComplex(), b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab*bb, a))
}

const DivideHelp = "Divides two numbers"

func Divide(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		if bc == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()/bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		bc := b.UnsafeComplex()
		if bc == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac/bc, b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	if bb == 0 {
		return rpn.ErrDivideByZero
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab/bb, a))
}

const NegateHelp = "Negates the top number"

func Negate(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsComplex() {
		c, _ := f.Complex()
		return r.PushFrame(rpn.ComplexFrameWithType(-c, f.Type()))
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

const PolarHelp = "Converts head element to a complex polar"

func Polar(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsBool() {
		if f.UnsafeBool() {
			return r.PushFrame(rpn.PolarFrame(1, 0, r.AngleUnit))
		} else {
			return r.PushFrame(rpn.PolarFrame(0, 0, r.AngleUnit))
		}
	}
	if f.IsString() {
		err := r.Exec(f.UnsafeString())
		if err != nil {
			return err
		}
		f, err = r.PopFrame()
		if err != nil {
			return err
		}
	}
	v, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(v, r.AngleUnit))
}

const FloatHelp = "Converts head element to a complex float"

func Float(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsBool() {
		if f.UnsafeBool() {
			return r.PushFrame(rpn.RealFrame(1))
		} else {
			return r.PushFrame(rpn.RealFrame(0))
		}
	}
	if f.IsString() {
		err := r.Exec(f.UnsafeString())
		if err != nil {
			return err
		}
		f, err = r.PopFrame()
		if err != nil {
			return err
		}
	}
	v, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(v))
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

const ImagHelp = "Takes the imaginary portion of a complex number (as a real number)"

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
	if af.IsComplex() && (af.Type() != rpn.COMPLEX_FRAME) {
		rl, an := cmplx.Polar(a)
		for i := 0; i < int(b); i++ {
			rl *= 10
			an *= 10
		}
		rl = math.Round(rl)
		an = math.Round(an)
		for i := 0; i < int(b); i++ {
			rl /= 10
			an /= 10
		}
		return r.PushFrame(rpn.PolarFrame(rl, an, af.Type()))
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
