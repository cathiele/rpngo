// RPNStack holds stack information for an RPNCalc
package rpn

import (
	"errors"
	"fmt"
)

var (
	errStackEmpty           = errors.New("stack empty")
	errNotEnoughStackFrames = errors.New("not enough stack frames")
)

// Frame Defines a single stack frame
type Frame struct {
	Complex complex128
}

func (f *Frame) String() string {
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

func (s *Stack) Push(f Frame) error {
	s.frames = append(s.frames, f)
	return nil
}

func (s *Stack) Pop() (sf Frame, err error) {
	if len(s.frames) == 0 {
		err = errStackEmpty
		return
	}
	sf = s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-1]
	return
}

func (s *Stack) Pop2() (a Frame, b Frame, err error) {
	if len(s.frames) < 2 {
		err = errStackEmpty
		return
	}
	a = s.frames[len(s.frames)-2]
	b = s.frames[len(s.frames)-1]
	s.frames = s.frames[:len(s.frames)-2]
	return
}

func (s *Stack) Peek() (sf Frame, err error) {
	if len(s.frames) == 0 {
		err = errStackEmpty
		return
	}
	sf = s.frames[len(s.frames)-1]
	return
}

func (s *Stack) IterFrames(fn func(Frame)) {
	for _, sf := range s.frames {
		fn(sf)
	}
}
