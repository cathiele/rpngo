// RPNStack holds stack information for an RPNCalc
package rpn

import "errors"

var (
	errStackEmpty           = errors.New("stack empty")
	errNotEnoughStackFrames = errors.New("not enough stack frames")
)

// Frame Defines a single stack frame
type Frame struct {
	Float float64
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

func (s *Stack) IterFrames(fn func(Frame)) {
	for _, sf := range s.frames {
		fn(sf)
	}
}
