package rpn

import "strconv"

type FrameType uint8

const (
	EMPTY_FRAME FrameType = iota
	STRING_FRAME
	COMPLEX_FRAME
	BOOL_FRAME
	INTEGER_FRAME
	HEXIDECIMAL_FRAME
	OCTAL_FRAME
	BINARY_FRAME
)

// Frame Defines a single stack frame
type Frame struct {
	ftype FrameType
	str   string
	cmplx complex128
	intv  int64
}

func (f *Frame) IsInt() bool {
	switch f.ftype {
	case INTEGER_FRAME:
		return true
	case HEXIDECIMAL_FRAME:
		return true
	case OCTAL_FRAME:
		return true
	case BINARY_FRAME:
		return true
	default:
		return false
	}
}

func (f *Frame) IsComplex() bool {
	return f.ftype == COMPLEX_FRAME
}

func (f *Frame) IsNumber() bool {
	return f.IsComplex() || f.IsInt()
}

func (f *Frame) Complex() (complex128, error) {
	if f.ftype == COMPLEX_FRAME {
		return f.cmplx, nil
	}
	if f.IsInt() {
		return complex(float64(f.intv), 0), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) Real() (float64, error) {
	if f.ftype == COMPLEX_FRAME {
		if imag(f.cmplx) != 0 {
			return 0, ErrComplexNumberNotSupported
		}
		return real(f.cmplx), nil
	}
	if f.IsInt() {
		return float64(f.intv), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) Int() (int64, error) {
	if f.IsInt() {
		return f.intv, nil
	}
	if f.ftype == COMPLEX_FRAME {
		return int64(real(f.cmplx)), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) Bool() (bool, error) {
	if f.ftype == BOOL_FRAME {
		return f.intv != 0, nil
	}
	return false, ErrExpectedABoolean
}

func (f *Frame) String(quote bool) string {
	switch f.ftype {
	case EMPTY_FRAME:
		return "nil"
	case STRING_FRAME:
		if quote {
			return "\"" + f.str + "\""
		}
		return f.str
	case COMPLEX_FRAME:
		return f.complexString()
	case BOOL_FRAME:
		if f.intv != 0 {
			return "true"
		}
		return "false"
	case INTEGER_FRAME:
		return strconv.FormatInt(f.intv, 10) + "d"
	case HEXIDECIMAL_FRAME:
		return strconv.FormatInt(f.intv, 16) + "x"
	case OCTAL_FRAME:
		return strconv.FormatInt(f.intv, 8) + "o"
	case BINARY_FRAME:
		return strconv.FormatInt(f.intv, 2) + "b"
	default:
		return "BAD_TYPE"
	}
}

func BoolFrame(v bool) Frame {
	if v {
		return Frame{intv: 1, ftype: BOOL_FRAME}
	}
	return Frame{intv: 0, ftype: BOOL_FRAME}
}

func IntFrame(v int64, t FrameType) Frame {
	return Frame{intv: v, ftype: t}
}

func ComplexFrame(v complex128) Frame {
	return Frame{cmplx: v, ftype: COMPLEX_FRAME}
}

func RealFrame(v float64) Frame {
	return Frame{cmplx: complex(v, 0), ftype: COMPLEX_FRAME}
}

func StringFrame(v string) Frame {
	return Frame{str: v, ftype: STRING_FRAME}
}

func EmptyFrame() Frame {
	return Frame{ftype: EMPTY_FRAME}
}

func (f *Frame) complexString() string {
	if imag(f.cmplx) == 0 {
		return strconv.FormatFloat(real(f.cmplx), 'g', 10, 64)
	}
	if real(f.cmplx) == 0 {
		return complexString(imag(f.cmplx))
	}
	r := strconv.FormatFloat(real(f.cmplx), 'g', 10, 64)
	if imag(f.cmplx) < 0 {
		return r + complexString(imag(f.cmplx))
	}
	return r + "+" + complexString(imag(f.cmplx))
}

func complexString(v float64) string {
	if v == 1 {
		return "i"
	}
	if v == -1 {
		return "-i"
	}
	return strconv.FormatFloat(v, 'g', 10, 64) + "i"
}
