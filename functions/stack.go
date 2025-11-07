package functions

import (
	"mattwach/rpngo/rpn"
)

const DropAllHelp = "Clears the stack"

func DropAll(r *rpn.RPN) error {
	r.Clear()
	return nil
}
