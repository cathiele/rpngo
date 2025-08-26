package functions

import "mattwach/rpngo/rpn"

func Duplicate(s *rpn.Stack) error {
	a, err := s.PeekFrame(0)
	if err != nil {
		return err
	}
	return s.PushFrame(a)
}
