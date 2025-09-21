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
	if a.Type == rpn.STRING_FRAME || b.Type == rpn.STRING_FRAME {
		return r.PushString(a.String(false) + b.String(false))
	}
	if a.Type == rpn.BOOL_FRAME || b.Type == rpn.BOOL_FRAME {
		r.PushFrame(a)
		r.PushFrame(b)
		return rpn.ErrIllegalValue
	}
	intMask := 0
	if a.IsInt() {
		intMask |= 1
	}
	if b.IsInt() {
		intMask |= 2
	}
	switch intMask {
	case 0:
		return r.PushComplex(a.Complex + b.Complex)
	case 1:
		return r.PushComplex(complex(float64(a.Int), 0) + b.Complex)
	case 2:
		return r.PushComplex(a.Complex + complex(float64(b.Int), 0))
	case 3:
		return r.PushInt(a.Int+b.Int, a.Type)
	}
	return nil
}

const SubtractHelp = "Subtracts two numbers"

func Subtract(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(a.Complex - b.Complex)
	}
	return r.PushInt(a.Int-b.Int, a.Type)
}

const MultiplyHelp = "Multiplies two numbers"

func Multiply(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(a.Complex * b.Complex)
	}
	return r.PushInt(a.Int*b.Int, a.Type)
}

const DivideHelp = "Divides two numbers"

func Divide(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		if b.Complex == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushComplex(a.Complex / b.Complex)
	}
	if b.Int == 0 {
		return rpn.ErrDivideByZero
	}
	return r.PushInt(a.Int/b.Int, a.Type)
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
	if f.IsInt() {
		f.Int = -f.Int
		return r.PushFrame(f)
	}
	r.PushFrame(f)
	return rpn.ErrIllegalValue
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

const RandHelp = "Pushes a random number between 0 and 1"

func Rand(r *rpn.RPN) error {
	return r.PushComplex(complex(rand.Float64(), 0))
}

const RealHelp = "Takes the real portion of a complex number"

func Real(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	if a.Type != rpn.COMPLEX_FRAME {
		r.PushFrame(a)
		return rpn.ErrExpectedAComplexNumber
	}
	a.Complex = complex(real(a.Complex), 0)
	return r.PushFrame(a)
}

const ImagHelp = "Takes the imaginary portion of a complex number"

func Imag(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	if a.Type != rpn.COMPLEX_FRAME {
		r.PushFrame(a)
		return rpn.ErrExpectedAComplexNumber
	}
	a.Complex = complex(0, imag(a.Complex))
	return r.PushFrame(a)
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
