package functions

import (
	"mattwach/rpngo/rpn"
)

func convert(r *rpn.RPN, t rpn.FrameType) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsBool() {
		if f.UnsafeBool() {
			return r.PushFrame(rpn.IntFrame(1, t))
		} else {
			return r.PushFrame(rpn.IntFrame(0, t))
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
	v, err := f.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrame(v, t))
}

const IntHelp = "Converts head element to an integer number"

func Int(r *rpn.RPN) error {
	return convert(r, rpn.INTEGER_FRAME)
}

const HexHelp = "Converts head element to a hexidecimal number"

func Hex(r *rpn.RPN) error {
	return convert(r, rpn.HEXIDECIMAL_FRAME)
}

const OctHelp = "Converts head element to an octal number"

func Oct(r *rpn.RPN) error {
	return convert(r, rpn.OCTAL_FRAME)
}

const BinHelp = "Converts head element to a binary number"

func Bin(r *rpn.RPN) error {
	return convert(r, rpn.BINARY_FRAME)
}

const StrHelp = "Converts head element to a string"

func Str(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(f.String(false), rpn.STRING_DOUBLE_QUOTE))
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
	return r.PushFrame(rpn.ComplexFrame(v, rpn.COMPLEX_FRAME))
}
