package functions

import "mattwach/rpngo/rpn"

func RegisterAll(rpn *rpn.RPN) {
	rpn.Register("+", Add)
	rpn.Register("-", Subtract)
	rpn.Register("*", Multiply)
	rpn.Register("/", Divide)
}
