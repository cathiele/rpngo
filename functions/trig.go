package functions

import (
	"errors"
	"math/cmplx"
	"mattwach/rpngo/rpn"
)

var errChooseDegRadOGrad = errors.New("choose 'deg', 'rad', or 'grad'")

const SetAngleHelp = "sets angle units to 'rad', 'deg', or 'grads'"

func SetAngle(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	switch f.UnsafeString() {
	case "rad":
		r.AngleUnits = rpn.ANGLE_RADIANS
	case "deg":
		r.AngleUnits = rpn.ANGLE_DEGREES
	case "grad":
		r.AngleUnits = rpn.ANGLE_GRADS
	default:
		return errChooseDegRadOGrad
	}
	return nil
}

const GetAngleHelp = "returns currently-set angle units"

func GetAngle(r *rpn.RPN) error {
	switch r.AngleUnits {
	case rpn.ANGLE_RADIANS:
		return r.PushFrame(rpn.StringFrame("rad", rpn.STRING_SINGLE_QUOTE))
	case rpn.ANGLE_DEGREES:
		return r.PushFrame(rpn.StringFrame("deg", rpn.STRING_SINGLE_QUOTE))
	case rpn.ANGLE_GRADS:
		return r.PushFrame(rpn.StringFrame("grad", rpn.STRING_SINGLE_QUOTE))
	}
	return rpn.ErrIllegalValue
}

const RadHelp = "sets trig / polar units to radians (calls 'rad' setangle)"

func Rad(r *rpn.RPN) error {
	if err := r.PushFrame(rpn.StringFrame("rad", rpn.STRING_SINGLE_QUOTE)); err != nil {
		return err
	}
	return SetAngle(r)
}

const DegHelp = "sets trig / polar units to degrees (calls 'deg' setangle)"

func Deg(r *rpn.RPN) error {
	if err := r.PushFrame(rpn.StringFrame("deg", rpn.STRING_SINGLE_QUOTE)); err != nil {
		return err
	}
	return SetAngle(r)
}

const GradHelp = "sets trig / polar units to grads (calls 'grad' setangle)"

func Grad(r *rpn.RPN) error {
	if err := r.PushFrame(rpn.StringFrame("grad", rpn.STRING_SINGLE_QUOTE)); err != nil {
		return err
	}
	return SetAngle(r)
}

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
	return r.PushFrame(rpn.ComplexFrameCloneType(cmplx.Sin(r.ToRadians(a)), af))
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
	return r.PushFrame(rpn.ComplexFrameCloneType(cmplx.Cos(r.ToRadians(a)), af))
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
	return r.PushFrame(rpn.ComplexFrameCloneType(cmplx.Tan(r.ToRadians(a)), af))
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
