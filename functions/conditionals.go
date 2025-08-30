package functions

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const tolerance = 1e-9

const GreaterThanHelp = "Returns true if a > b, false otherwise"

func GreaterThan(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushBool(real(a) > real(b))
}

const GreaterThanEqualHelp = "Returns true if a >= b, false otherwise"

func GreaterThanEqual(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushBool(real(a) >= real(b))
}

const LessThanHelp = "Returns true if a < b, false otherwise"

func LessThan(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushBool(real(a) < real(b))
}

const LessThanEqualHelp = "Returns true if a <= b, false otherwise"

func LessThanEqual(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushBool(real(a) <= real(b))
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
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushBool(checkEqual(a, b))
}

const NotEqualHelp = "Returns true if a != b, false otherwise (approximate)"

func NotEqual(r *rpn.RPN) error {
	a, b, err := r.Pop2Complex()
	if err != nil {
		return err
	}
	return r.PushBool(!checkEqual(a, b))
}

const ForHelp = "Executes the head of the stack in a loop until a value < is found\n" +
	"Example: 1 'c 1 + c 50 <' for # put 1 to 50 on the stack"

func For(r *rpn.RPN) error {
	macro, err := r.PopString()
	if err != nil {
		return err
	}
	fields, err := parse.Fields(macro)
	if err != nil {
		return err
	}
	for {
		for _, f := range fields {
			if err := r.Exec(f); err != nil {
				return err
			}
		}
		cond, err := r.PopBool()
		if err != nil {
			return err
		}
		if !cond {
			break
		}
	}
	return nil
}
