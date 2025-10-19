package functions

import (
	"mattwach/rpngo/rpn"
)

const tolerance = 1e-9

const GreaterThanHelp = "Returns true if a > b, false otherwise"

func GreaterThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(af, bf)
	if err != nil {
		return err
	}
	if bothints {
		return r.PushFrame(rpn.BoolFrame(af.UnsafeInt() > bf.UnsafeInt()))
	}
	a, err := af.Real()
	if err != nil {
		return err
	}
	b, err := bf.Real()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(a > b))
}

const GreaterThanEqualHelp = "Returns true if a >= b, false otherwise"

func GreaterThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(af, bf)
	if err != nil {
		return err
	}
	if bothints {
		return r.PushFrame(rpn.BoolFrame(af.UnsafeInt() >= bf.UnsafeInt()))
	}
	a, err := af.Real()
	if err != nil {
		return err
	}
	b, err := bf.Real()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(a >= b))
}

const LessThanHelp = "Returns true if a < b, false otherwise"

func LessThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(af, bf)
	if err != nil {
		return err
	}
	if bothints {
		return r.PushFrame(rpn.BoolFrame(af.UnsafeInt() < bf.UnsafeInt()))
	}
	a, err := af.Real()
	if err != nil {
		return err
	}
	b, err := bf.Real()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(a < b))
}

const LessThanEqualHelp = "Returns true if a <= b, false otherwise"

func LessThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	bothints, err := rpn.CheckIfNumbers(af, bf)
	if err != nil {
		return err
	}
	if bothints {
		return r.PushFrame(rpn.BoolFrame(af.UnsafeInt() <= bf.UnsafeInt()))
	}
	a, err := af.Real()
	if err != nil {
		return err
	}
	b, err := bf.Real()
	if err != nil {
		return err
	}
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
		b, err := bf.Complex()
		if err != nil {
			return false, nil
		}
		return checkFloatEqual(a, b), nil
	} else if af.IsInt() && bf.IsInt() {
		return af.UnsafeInt() == bf.UnsafeInt(), nil
	} else if af.IsBool() && bf.IsBool() {
		return af.UnsafeBool() == bf.UnsafeBool(), nil
	} else if af.IsString() && bf.IsString() {
		return af.UnsafeString() == bf.UnsafeString(), nil
	}
	return false, nil
}
