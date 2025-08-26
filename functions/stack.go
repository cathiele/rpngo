package functions

import "mattwach/rpngo/rpn"

func Duplicate(s *rpn.Stack) error {
	a, err := s.Peek(0)
	if err != nil {
		return err
	}
	return s.Push(rpn.Frame{Complex: a.Complex})
}
