package functions

import "mattwach/rpngo/rpn"

// RegisterAll resisters functions that are core to the RPN engine. There
// are other functions, such as window-specific functions, that will be
// added by their respective owners as the functions module should not
// know about these.
func RegisterAll(rpn *rpn.RPN) {
	rpn.Register(".", Duplicate)
	rpn.Register("a", Add)
	rpn.Register("+", Add)
	rpn.Register("s", Subtract)
	rpn.Register("-", Subtract)
	rpn.Register("m", Multiply)
	rpn.Register("*", Multiply)
	rpn.Register("d", Divide)
	rpn.Register("/", Divide)
	rpn.Register("sq", Square)
	rpn.Register("sqrt", SquareRoot)
	rpn.Register("cos", Cos)
	rpn.Register("sin", Sin)
	rpn.Register("tan", Tan)
}
