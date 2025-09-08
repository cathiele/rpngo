// RPNStack holds stack information for an RPNCalc
package rpn

import (
	"fmt"
)

func (f *Frame) IsInt() bool {
	switch f.Type {
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

func (f *Frame) Bool() bool {
	return f.Int != 0
}

func BoolFrame(v bool) Frame {
	if v {
		return Frame{Int: 1, Type: BOOL_FRAME}
	}
	return Frame{Int: 0, Type: BOOL_FRAME}
}

func (f *Frame) String(quote bool) string {
	switch f.Type {
	case EMPTY_FRAME:
		return "nil"
	case STRING_FRAME:
		if quote {
			return "\"" + f.Str + "\""
		}
		return f.Str
	case COMPLEX_FRAME:
		return f.complexString()
	case BOOL_FRAME:
		if f.Int != 0 {
			return "true"
		}
		return "false"
	case INTEGER_FRAME:
		return fmt.Sprintf("%vd", f.Int)
	case HEXIDECIMAL_FRAME:
		return fmt.Sprintf("%xx", f.Int)
	case OCTAL_FRAME:
		return fmt.Sprintf("%oo", f.Int)
	case BINARY_FRAME:
		return fmt.Sprintf("%bb", f.Int)
	default:
		return "BAD_TYPE"
	}
}

func (f *Frame) complexString() string {
	if imag(f.Complex) == 0 {
		return fmt.Sprintf("%g", real(f.Complex))
	}
	if real(f.Complex) == 0 {
		return complexString(imag(f.Complex))
	}
	if imag(f.Complex) < 0 {
		return fmt.Sprintf("%g%s", real(f.Complex), complexString(imag(f.Complex)))
	}
	return fmt.Sprintf("%g+%s", real(f.Complex), complexString(imag(f.Complex)))
}

func complexString(v float64) string {
	if v == 1 {
		return "i"
	}
	if v == -1 {
		return "-i"
	}
	return fmt.Sprintf("%gi", v)
}

func (r *RPN) Clear() {
	r.frames = r.frames[:0]
}

func (r *RPN) StackLen() int {
	return len(r.frames)
}

func (r *RPN) PushFrame(f Frame) error {
	if len(r.frames) >= MaxStackDepth {
		return ErrStackFull
	}
	r.frames = append(r.frames, f)
	return nil
}

func (r *RPN) PushComplex(v complex128) error {
	return r.PushFrame(Frame{Type: COMPLEX_FRAME, Complex: v})
}

func (r *RPN) PushString(v string) error {
	return r.PushFrame(Frame{Type: STRING_FRAME, Str: v})
}

func (r *RPN) PushBool(v bool) error {
	var val int64
	if v {
		val = 1
	}
	return r.PushFrame(Frame{Type: BOOL_FRAME, Int: val})
}

func (r *RPN) PushInt(v int64, t FrameType) error {
	return r.PushFrame(Frame{Type: t, Int: v})
}

func (r *RPN) PopFrame() (sf Frame, err error) {
	if len(r.frames) == 0 {
		err = ErrStackEmpty
		return
	}
	sf = r.frames[len(r.frames)-1]
	r.frames = r.frames[:len(r.frames)-1]
	return
}

func (r *RPN) Pop2Frames() (a Frame, b Frame, err error) {
	if len(r.frames) < 2 {
		err = ErrNotEnoughStackFrames
		return
	}
	a = r.frames[len(r.frames)-2]
	b = r.frames[len(r.frames)-1]
	r.frames = r.frames[:len(r.frames)-2]
	return
}

func (r *RPN) PopString() (str string, err error) {
	f, err := r.PopFrame()
	if err != nil {
		return
	}
	if f.Type != STRING_FRAME {
		r.PushFrame(f)
		err = ErrExpectedAString
		return
	}
	str = f.Str
	return
}

func (r *RPN) PopBool() (v bool, err error) {
	f, err := r.PopFrame()
	if err != nil {
		return
	}
	if f.Type != BOOL_FRAME {
		r.PushFrame(f)
		err = ErrExpectedABoolean
		return
	}
	v = f.Int != 0
	return
}

func (r *RPN) Pop2Strings() (a string, b string, err error) {
	as, bs, err := r.Pop2Frames()
	if err != nil {
		return
	}
	if as.Type != STRING_FRAME || bs.Type != STRING_FRAME {
		r.PushFrame(as)
		r.PushFrame(bs)
		err = ErrExpectedAString
		return
	}
	a = as.Str
	b = bs.Str
	return
}

func (r *RPN) PopNumber() (f Frame, err error) {
	f, err = r.PopFrame()
	if err != nil {
		return
	}
	if (f.Type != COMPLEX_FRAME) && !f.IsInt() {
		r.PushFrame(f)
		err = ErrExpectedANumber
		return
	}
	return
}

func (r *RPN) PopComplex() (v complex128, err error) {
	f, err := r.PopFrame()
	if err != nil {
		return
	}
	if f.Type == COMPLEX_FRAME {
		v = f.Complex
		return
	}
	if f.IsInt() {
		v = complex(float64(f.Int), 0)
		return
	}
	r.PushFrame(f)
	err = ErrExpectedANumber
	return
}

func (r *RPN) PopReal() (v float64, err error) {
	f, err := r.PopFrame()
	if err != nil {
		return
	}
	if f.Type == COMPLEX_FRAME {
		if imag(f.Complex) != 0 {
			err = ErrComplexNumberNotSupported
			r.PushFrame(f)
			return
		}
		v = real(f.Complex)
		return
	}
	if f.IsInt() {
		v = float64(f.Int)
		return
	}
	r.PushFrame(f)
	err = ErrExpectedANumber
	return
}

// Pops 2 numbers.
//
// If either is a non-number, an error is returned
// If either number is a complex, both are converted to complex
// Ohterwise both number are integers and they are returned.
func (r *RPN) Pop2Numbers() (a Frame, b Frame, err error) {
	a, b, err = r.Pop2Frames()
	if err != nil {
		return
	}
	if (a.Type == COMPLEX_FRAME) && (b.Type == COMPLEX_FRAME) {
		return
	}
	if (a.Type == COMPLEX_FRAME) && b.IsInt() {
		b.Type = COMPLEX_FRAME
		b.Complex = complex(float64(b.Int), 0)
		return
	}
	if (b.Type == COMPLEX_FRAME) && a.IsInt() {
		a.Type = COMPLEX_FRAME
		a.Complex = complex(float64(a.Int), 0)
		return
	}
	if a.IsInt() && b.IsInt() {
		return
	}
	r.PushFrame(a)
	r.PushFrame(b)
	err = ErrExpectedANumber
	return
}

func (r *RPN) PopStackIndex() (i int, err error) {
	var f Frame
	f, err = r.PopNumber()
	if err != nil {
		return
	}
	i = int(f.Int)
	if f.Type == COMPLEX_FRAME {
		if imag(f.Complex) != 0 {
			err = ErrComplexNumberNotSupported
			return
		}
		i = int(real(f.Complex))
	}
	if i < 0 {
		err = ErrIllegalValue
		return
	}
	if i >= len(r.frames) {
		err = ErrIllegalValue
		return
	}
	return
}

func (r *RPN) PeekFrame(framesBack int) (sf Frame, err error) {
	if len(r.frames)-framesBack <= 0 {
		err = ErrStackEmpty
		return
	}
	sf = r.frames[len(r.frames)-1-framesBack]
	return
}

func (r *RPN) DeleteFrame(framesBack int) (sf Frame, err error) {
	sf, err = r.PeekFrame(framesBack)
	if err != nil {
		return
	}
	idx := len(r.frames) - 1 - framesBack
	r.frames = append(r.frames[:idx], r.frames[idx+1:]...)
	return
}

func (r *RPN) IterFrames(fn func(Frame)) {
	for _, sf := range r.frames {
		fn(sf)
	}
}

func (r *RPN) Size() int {
	return len(r.frames)
}
