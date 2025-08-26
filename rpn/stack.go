// RPNStack holds stack information for an RPNCalc
package rpn

import (
	"errors"
	"fmt"
)

var (
	errExpectedANumber      = errors.New("expected a number")
	errStackEmpty           = errors.New("stack empty")
	errNotEnoughStackFrames = errors.New("not enough stack frames")
)

type FrameType uint8

const (
	STRING_FRAME FrameType = iota
	COMPLEX_FRAME
)

// Frame Defines a single stack frame
type Frame struct {
	Type    FrameType
	Str     string
	Complex complex128
}

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

func complexString(v float64) string {
	if v == 1 {
		return "i"
	}
	if v == -1 {
		return "-i"
	}
	return fmt.Sprintf("%gi", v)
}

// Stak defines a full stack frame
type Stack struct {
	frames []Frame
}

func (s *Stack) Clear() {
	s.frames = s.frames[:0]
}

func (s *Stack) PushComplex(v complex128) error {
	s.frames = append(s.frames, Frame{Type: COMPLEX_FRAME, Complex: v})
	return nil
}

func (s *Stack) PushString(v string) error {
	s.frames = append(s.frames, Frame{Type: STRING_FRAME, Str: v})
	return nil
}

func (s *Stack) PushFrame(f Frame) error {
	s.frames = append(s.frames, Frame{f.Type, f.Str, f.Complex})
	return nil
}

func (s *Stack) PopFrame() (sf Frame, err error) {
	if len(s.frames) == 0 {
		err = errStackEmpty
		return
	}
	sf = s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-1]
	return
}

func (s *Stack) Pop2Frames() (a Frame, b Frame, err error) {
	if len(s.frames) < 2 {
		err = errNotEnoughStackFrames
		return
	}
	a = s.frames[len(s.frames)-2]
	b = s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-2]
	return
}

func (s *Stack) PopComplex() (v complex128, err error) {
	f, err := s.PopFrame()
	if err != nil {
		return
	}
	if f.Type != COMPLEX_FRAME {
		s.PushFrame(f)
		err = errExpectedANumber
		return
	}
	v = f.Complex
	return
}

func (s *Stack) Pop2Complex() (a complex128, b complex128, err error) {
	af, bf, err := s.Pop2Frames()
	if err != nil {
		return
	}
	if af.Type != COMPLEX_FRAME || bf.Type != COMPLEX_FRAME {
		s.PushFrame(af)
		s.PushFrame(bf)
		err = errExpectedANumber
		return
	}
	a = af.Complex
	b = bf.Complex
	return
}

func (s *Stack) PeekFrame(framesBack int) (sf Frame, err error) {
	if len(s.frames)-framesBack <= 0 {
		err = errStackEmpty
		return
	}
	sf = s.frames[len(s.frames)-1-framesBack]
	return
}

func (s *Stack) IterFrames(fn func(Frame)) {
	for _, sf := range s.frames {
		fn(sf)
	}
}

func (s *Stack) Size() int {
	return len(s.frames)
}
