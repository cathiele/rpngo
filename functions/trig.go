package functions

import (
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

const SinHelp = "takes the sine of a number"

func Sin(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Sin(r.ToRadians(a)), af.Type()))
}

const ASinHelp = "takes the inverse sine of a number"

func ASin(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(r.FromRadians(cmplx.Asin(a), af))
}

const CosHelp = "takes the cosine of a number"

func Cos(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Cos(r.ToRadians(a)), af.Type()))
}

const ACosHelp = "takes the inverse cosine of a number"

func ACos(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(r.FromRadians(cmplx.Acos(a), af))
}

const TanHelp = "takes the tangent of a number"

func Tan(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Tan(r.ToRadians(a)), af.Type()))
}

const ATanHelp = "takes the inverse tangent of a number"

func ATan(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(r.FromRadians(cmplx.Atan(a), af))
}
