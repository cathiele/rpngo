package functions

import (
	"errors"
	"mattwach/rpngo/rpn"
)

func convert(r *rpn.RPN, t rpn.FrameType) error {
	v, err := r.PopFrame()
	if err != nil {
		return err
	}
	switch v.Type {
	case rpn.STRING_FRAME:
		err := r.Exec([]string{v.Str})
		if err != nil {
			return err
		}
		return convert(r, t)
	case rpn.COMPLEX_FRAME:
		v.Int = int64(real(v.Complex))
		v.Type = t
	case rpn.BINARY_FRAME:
		fallthrough
	case rpn.HEXIDECIMAL_FRAME:
		fallthrough
	case rpn.INTEGER_FRAME:
		fallthrough
	case rpn.BOOL_FRAME:
		fallthrough
	case rpn.OCTAL_FRAME:
		v.Type = t
	default:
		return errors.New("bad frame type")
	}
	return r.PushFrame(v)
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
	return r.PushString(f.String(false))
}

const FloatHelp = "Converts head element to a complex float"

func Float(r *rpn.RPN) error {
	v, err := r.PopFrame()
	if err != nil {
		return err
	}
	switch v.Type {
	case rpn.STRING_FRAME:
		err := r.Exec([]string{v.Str})
		if err != nil {
			return err
		}
		return Float(r)
	case rpn.COMPLEX_FRAME:
		break
	case rpn.BINARY_FRAME:
		fallthrough
	case rpn.HEXIDECIMAL_FRAME:
		fallthrough
	case rpn.INTEGER_FRAME:
		fallthrough
	case rpn.BOOL_FRAME:
		fallthrough
	case rpn.OCTAL_FRAME:
		v.Type = rpn.COMPLEX_FRAME
		v.Complex = complex(float64(v.Int), 0)
	default:
		return errors.New("bad frame type")
	}
	return r.PushFrame(v)
}
