package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

const PowerHelp = "executes a to the power of b"

func Power(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		b, err := bf.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Pow(af.UnsafeComplex(), b), af.Type()))
	}
	if bf.IsComplex() {
		a, err := af.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Pow(a, bf.UnsafeComplex()), bf.Type()))
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(powInts(a, b), af))
}

func powInts(x, n int64) int64 {
	if n < 0 {
		return 0
	}
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	y := powInts(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return x * y * y
}

const SquareRootHelp = "takes the square root of a number"

func SquareRoot(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	a = cmplx.Sqrt(a)
	if af.IsComplex() {
		return r.PushFrame(rpn.ComplexFrameWithType(a, af.Type()))
	}
	return r.PushFrame(rpn.IntFrameCloneType(int64(real(a)), af))
}

const AbsHelp = "Takes the absolute value"

func Abs(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		a, _ := af.Complex()
		return r.PushFrame(rpn.RealFrame(cmplx.Abs(a)))
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	if a < 0 {
		a = -a
	}
	return r.PushFrame(rpn.IntFrameCloneType(a, af))
}

const SquareHelp = "executes v * v"

func Square(r *rpn.RPN) error {
	a, err := r.PopFrame()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		ac := a.UnsafeComplex()
		return r.PushFrame(rpn.ComplexFrameWithType(ac*ac, a.Type()))
	}
	ai, err := a.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ai*ai, a))
}

const LogHelp = "executes natural logrithm"

func Log(r *rpn.RPN) error {
	a, err := r.PopFrame()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Log(a.UnsafeComplex()), a.Type()))
	}
	ac, err := a.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Log(ac)))
}

const Log10Help = "executes base 10 logrithm"

func Log10(r *rpn.RPN) error {
	a, err := r.PopFrame()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Log10(a.UnsafeComplex()), a.Type()))
	}
	ac, err := a.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Log10(ac)))
}
