// RPNStack holds stack information for an RPNCalc
package rpn

import (
	"errors"
	"fmt"
)

func (f *Frame) String() string {
	switch f.Type {
	case STRING_FRAME:
		return "\"" + f.Str + "\""
	case COMPLEX_FRAME:
		return f.complexString()
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

func (f *Frame) Copy() Frame {
	return Frame{f.Type, f.Str, f.Complex}
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

func (r *RPN) PushComplex(v complex128) error {
	r.frames = append(r.frames, Frame{Type: COMPLEX_FRAME, Complex: v})
	return nil
}

func (r *RPN) PushString(v string) error {
	r.frames = append(r.frames, Frame{Type: STRING_FRAME, Str: v})
	return nil
}

func (r *RPN) PushFrame(f Frame) error {
	r.frames = append(r.frames, Frame{f.Type, f.Str, f.Complex})
	return nil
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
		err = errNotEnoughStackFrames
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
		err = errExpectedAString
		return
	}
	str = f.Str
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
		err = errExpectedAString
		return
	}
	a = as.Str
	b = bs.Str
	return
}

func (r *RPN) PopComplex() (v complex128, err error) {
	f, err := r.PopFrame()
	if err != nil {
		return
	}
	if f.Type != COMPLEX_FRAME {
		r.PushFrame(f)
		err = errExpectedANumber
		return
	}
	v = f.Complex
	return
}

func (r *RPN) Pop2Complex() (a complex128, b complex128, err error) {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return
	}
	if af.Type != COMPLEX_FRAME || bf.Type != COMPLEX_FRAME {
		r.PushFrame(af)
		r.PushFrame(bf)
		err = errExpectedANumber
		return
	}
	a = af.Complex
	b = bf.Complex
	return
}

func (r *RPN) PopStackIndex() (i int, err error) {
	var v complex128
	v, err = r.PopComplex()
	if err != nil {
		return
	}
	if imag(v) != 0 {
		err = errors.New("real number required")
		return
	}
	i = int(real(v))
	if i < 0 {
		err = errors.New("index must be >= 0")
		return
	}
	if i >= len(r.frames) {
		err = fmt.Errorf("index too high: %d", i)
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
