package functions

import (
	"mattwach/rpngo/rpn"
)

const tolerance = 1e-9

const GreaterThanHelp = "Returns true if a > b, false otherwise"

func GreaterThan(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushBool(real(a.Complex) > real(b.Complex))
	}
	return r.PushBool(a.Int > b.Int)
}

const GreaterThanEqualHelp = "Returns true if a >= b, false otherwise"

func GreaterThanEqual(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushBool(real(a.Complex) >= real(b.Complex))
	}
	return r.PushBool(a.Int >= b.Int)
}

const LessThanHelp = "Returns true if a < b, false otherwise"

func LessThan(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushBool(real(a.Complex) < real(b.Complex))
	}
	return r.PushBool(a.Int < b.Int)
}

const LessThanEqualHelp = "Returns true if a <= b, false otherwise"

func LessThanEqual(r *rpn.RPN) error {
	a, b, err := r.Pop2Numbers()
	if err != nil {
		return err
	}
	if a.Type == rpn.COMPLEX_FRAME {
		return r.PushBool(real(a.Complex) <= real(b.Complex))
	}
	return r.PushBool(a.Int <= b.Int)
}

func checkEqual(a, b complex128) bool {
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
	return r.PushBool(eq)
}

const NotEqualHelp = "Returns true if a != b, false otherwise (approximate)"

func NotEqual(r *rpn.RPN) error {
	eq, err := commonEqual(r)
	if err != nil {
		return err
	}
	return r.PushBool(!eq)
}

func commonEqual(r *rpn.RPN) (bool, error) {
	a, b, err := r.Pop2Frames()
	if err != nil {
		r.PushFrame(a)
		r.PushFrame(b)
		return false, err
	}
	if a.Type == rpn.COMPLEX_FRAME && b.IsInt() {
		b.Type = rpn.COMPLEX_FRAME
		b.Complex = complex(float64(b.Int), 0)
	} else if b.Type == rpn.COMPLEX_FRAME && a.IsInt() {
		a.Type = rpn.COMPLEX_FRAME
		a.Complex = complex(float64(a.Int), 0)
	}
	if a.Type != b.Type {
		return false, nil
	}
	if a.IsInt() || a.Type == rpn.BOOL_FRAME {
		return a.Int == b.Int, nil
	}
	switch a.Type {
	case rpn.COMPLEX_FRAME:
		return checkEqual(a.Complex, b.Complex), nil
	case rpn.STRING_FRAME:
		return a.Str == b.Str, nil
	}
	return false, nil
}
