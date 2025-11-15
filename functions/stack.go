package functions

import (
	"mattwach/rpngo/rpn"
)

const dropAllHelp = "Clears the stack"

func dropAll(r *rpn.RPN) error {
	r.Clear()
	return nil
}
