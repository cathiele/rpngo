package functions

import (
	"mattwach/rpngo/rpn"
)

const tolerance = 1e-9

const GreaterThanHelp = "Returns true if a > b, false otherwise"

func GreaterThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		a, _ := af.Complex()
		b, _ := af.Complex()
		return r.PushFrame(rpn.BoolFrame(real(a) > real(b)))
	}
	a, _ := af.Int()
	b, _ := bf.Int()
	return r.PushFrame(rpn.BoolFrame(a > b))
}

const GreaterThanEqualHelp = "Returns true if a >= b, false otherwise"

func GreaterThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		a, _ := af.Complex()
		b, _ := af.Complex()
		return r.PushFrame(rpn.BoolFrame(real(a) >= real(b)))
	}
	a, _ := af.Int()
	b, _ := bf.Int()
	return r.PushFrame(rpn.BoolFrame(a >= b))
}

const LessThanHelp = "Returns true if a < b, false otherwise"

func LessThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		a, _ := af.Complex()
		b, _ := af.Complex()
		return r.PushFrame(rpn.BoolFrame(real(a) < real(b)))
	}
	a, _ := af.Int()
	b, _ := bf.Int()
	return r.PushFrame(rpn.BoolFrame(a < b))
}

const LessThanEqualHelp = "Returns true if a <= b, false otherwise"

func LessThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		a, _ := af.Complex()
		b, _ := af.Complex()
		return r.PushFrame(rpn.BoolFrame(real(a) <= real(b)))
	}
	a, _ := af.Int()
	b, _ := bf.Int()
	return r.PushFrame(rpn.BoolFrame(a <= b))
}

func checkFloatEqual(a, b complex128) bool {
	d := real(a) - real(b)
	if d < -tolerance {
		return false
	}
	if d > tolerance {
		return false
	}
	d = imag(a) - imag(b)
	if d < -tolerance {
		return false
	}
	return d <= tolerance
}

const EqualHelp = "Returns true if a = b, false otherwise (approximate)"

func Equal(r *rpn.RPN) error {
	eq, err := commonEqual(r)
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(eq))
}

const NotEqualHelp = "Returns true if a != b, false otherwise (approximate)"

func NotEqual(r *rpn.RPN) error {
	eq, err := commonEqual(r)
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!eq))
}

func commonEqual(r *rpn.RPN) (bool, error) {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return false, err
	}
	if af.IsComplex() || bf.IsComplex() {
		a, err := af.Complex()
		if err != nil {
			return false, nil
		}
		b, err := af.Complex()
		if err != nil {
			return false, nil
		}
		return checkFloatEqual(a, b), nil
	} else if af.IsInt() && bf.IsInt() {
		a, _ := af.Int()
		b, _ := bf.Int()
		return a == b, nil
	} else if af.IsBool() && bf.IsBool() {
		a, _ := af.Bool()
		b, _ := bf.Bool()
		return a == b, nil
	} else if af.IsString() && bf.IsString() {
		return af.String(false) == bf.String(false), nil
	}
	return false, nil
}
