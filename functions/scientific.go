package functions

import (
	"math"
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

const PowerHelp = "executes a to the power of b"

func Power(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(cmplx.Pow(a.Complex, b.Complex))
	}
	return r.PushInt(int64(math.Pow(float64(a.Int), float64(b.Int))), a.Type)
}

const SquareRootHelp = "takes the square root of a complex number"

func SquareRoot(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	c := a.Complex
	if a.Type != rpn.COMPLEX_FRAME {
		c = complex(float64(a.Int), 0)
	}
	return r.PushComplex(cmplx.Sqrt(c))
}

const AbsHelp = "Takes the absolute value"

func Abs(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(complex(cmplx.Abs(a.Complex), 0))
	}
	iv := a.Int
	if iv < 0 {
		iv = -iv
	}
	return r.PushInt(iv, a.Type)
}

const SquareHelp = "executes v * v"

func Square(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(a.Complex * a.Complex)
	}
	return r.PushInt(a.Int*a.Int, a.Type)
}

const LogHelp = "executes natural logrithm"

func Log(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(cmplx.Log(a.Complex))
	}
	return r.PushComplex(complex(math.Log(float64(a.Int)), 0))
}

const Log10Help = "executes base 10 logrithm"

func Log10(r *rpn.RPN) error {
	a, err := r.PopNumber()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushComplex(cmplx.Log10(a.Complex))
	}
	return r.PushComplex(complex(math.Log10(float64(a.Int)), 0))
}
