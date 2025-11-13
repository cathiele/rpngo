package rpn

import (
	"math/cmplx"
	"strconv"
)

type FrameType uint8

const (
	EMPTY_FRAME FrameType = iota
	STRING_FRAME
	COMPLEX_FRAME
	POLAR_FRAME
	BOOL_FRAME
	INTEGER_FRAME
	HEXIDECIMAL_FRAME
	OCTAL_FRAME
	BINARY_FRAME
)

const (
	STRING_SINGLE_QUOTE = iota
	STRING_DOUBLE_QUOTE
	STRING_BRACES
)

// Frame Defines a single stack frame
type Frame struct {
	ftype FrameType
	str   string
	cmplx complex128
	// If ftype == BOOL_FRAME, intv holds 1 or 0
	// if ftype == STRING_FRAME, intv holds the quote type
	// If ftype == POLAR_FRAME, intv holds AngleUnit
	intv int64
}

// Annotates a frame.  Don't call this on string frames
// or the string will be replaced.
func (f *Frame) Annotate(s string) {
	f.str = s
}

func (f *Frame) IsInt() bool {
	switch f.ftype {
	case INTEGER_FRAME, HEXIDECIMAL_FRAME, OCTAL_FRAME, BINARY_FRAME:
		return true
	default:
		return false
	}
}

func (f *Frame) IsComplex() bool {
	return (f.ftype == COMPLEX_FRAME) || (f.ftype == POLAR_FRAME)
}

func (f *Frame) IsPolarComplex() bool {
	return f.ftype == POLAR_FRAME
}

func (f *Frame) IsNumber() bool {
	return f.IsComplex() || f.IsInt()
}

func (f *Frame) IsBool() bool {
	return f.ftype == BOOL_FRAME
}

func (f *Frame) IsString() bool {
	return f.ftype == STRING_FRAME
}

func (f *Frame) QuoteType() int {
	return int(f.intv)
}

func (f *Frame) Complex() (complex128, error) {
	if (f.ftype == COMPLEX_FRAME) || (f.ftype == POLAR_FRAME) {
		return f.cmplx, nil
	}
	if f.IsInt() {
		return complex(float64(f.intv), 0), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) UnsafeComplex() complex128 {
	if (f.ftype == COMPLEX_FRAME) || (f.ftype == POLAR_FRAME) {
		return f.cmplx
	}
	return complex(float64(f.intv), 0)
}

func (f *Frame) Real() (float64, error) {
	if (f.ftype == COMPLEX_FRAME) || (f.ftype == POLAR_FRAME) {
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
	if (f.ftype == COMPLEX_FRAME) || (f.ftype == POLAR_FRAME) {
		if imag(f.cmplx) != 0 {
			return 0, ErrComplexNumberNotSupported
		}
		return int64(real(f.cmplx)), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) UnsafeInt() int64 {
	return f.intv
}

func (f *Frame) BoundedInt(min, max int64) (int64, error) {
	v, err := f.Int()
	if err != nil {
		return 0, err
	}
	if (v < min) || (v > max) {
		return 0, ErrIllegalValue
	}
	return v, nil
}

func (f *Frame) Bool() (bool, error) {
	if f.ftype == BOOL_FRAME {
		return f.intv != 0, nil
	}
	return false, ErrExpectedABoolean
}

func (f *Frame) UnsafeBool() bool {
	return f.intv != 0
}

func (f *Frame) String(quote bool) string {
	var s string
	switch f.ftype {
	case EMPTY_FRAME:
		return "nil"
	case STRING_FRAME:
		if quote {
			switch f.intv {
			case STRING_SINGLE_QUOTE:
				return "'" + f.str + "'"
			case STRING_DOUBLE_QUOTE:
				return "\"" + f.str + "\""
			case STRING_BRACES:
				return "{" + f.str + "}"
			}
		}
		return f.str
	case COMPLEX_FRAME:
		s = f.complexString()
	case POLAR_FRAME:
		s = f.polarString()
	case BOOL_FRAME:
		if f.intv != 0 {
			s = "true"
		} else {
			s = "false"
		}
	case INTEGER_FRAME:
		s = strconv.FormatInt(f.intv, 10) + "d"
	case HEXIDECIMAL_FRAME:
		s = strconv.FormatInt(f.intv, 16) + "x"
	case OCTAL_FRAME:
		s = strconv.FormatInt(f.intv, 8) + "o"
	case BINARY_FRAME:
		s = strconv.FormatInt(f.intv, 2) + "b"
	default:
		return "BAD_TYPE"
	}
	if len(f.str) > 0 {
		s += " " + f.str
	}
	return s
}

func (f *Frame) UnsafeString() string {
	return f.str
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

func IntFrameCloneType(v int64, f Frame) Frame {
	return Frame{intv: v, ftype: f.ftype}
}

func ComplexFrame(v complex128) Frame {
	return Frame{cmplx: v, ftype: COMPLEX_FRAME}
}

func ComplexFrameCloneType(v complex128, f Frame) Frame {
	return Frame{cmplx: v, ftype: f.ftype, intv: f.intv}
}

func PolarFrame2(r, a float64, u AngleUnit) Frame {
	return Frame{cmplx: cmplx.Rect(r, a), ftype: POLAR_FRAME, intv: int(u)}
}

func RealFrame(v float64) Frame {
	return Frame{cmplx: complex(v, 0), ftype: COMPLEX_FRAME}
}

func StringFrame(v string, quoteType int) Frame {
	return Frame{str: v, ftype: STRING_FRAME, intv: int64(quoteType)}
}

func EmptyFrame() Frame {
	return Frame{ftype: EMPTY_FRAME}
}

func (f *Frame) complexString() string {
	if imag(f.cmplx) == 0 {
		return strconv.FormatFloat(real(f.cmplx), 'g', 16, 64)
	}
	if real(f.cmplx) == 0 {
		return complexString(imag(f.cmplx))
	}
	r := strconv.FormatFloat(real(f.cmplx), 'g', 16, 64)
	if imag(f.cmplx) < 0 {
		return r + complexString(imag(f.cmplx))
	}
	return r + "+" + complexString(imag(f.cmplx))
}

func (f *Frame) polarString() string {
	r, a := cmplx.Polar(f.cmplx)
	return strconv.FormatFloat(r, 'g', 16, 64) + "<" + strconv.FormatFloat(FromRadiansBase(a), 'g', 16, 64)
}

func complexString(v float64) string {
	if v == 1 {
		return "i"
	}
	if v == -1 {
		return "-i"
	}
	return strconv.FormatFloat(v, 'g', 16, 64) + "i"
}
