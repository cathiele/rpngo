package functions

import "mattwach/rpngo/rpn"

// RegisterAll resisters functions that are core to the RPN engine. There
// are other functions, such as window-specific functions, that will be
// added by their respective owners as the functions module should not
// know about these.
func RegisterAll(rpn *rpn.RPN) {
	rpn.Register(".", Duplicate, DuplicateHelp)
	rpn.Register("a", Add, AddHelp)
	rpn.Register("+", Add, AddHelp)
	rpn.Register("s", Subtract, SubtractHelp)
	rpn.Register("-", Subtract, SubtractHelp)
	rpn.Register("m", Multiply, MultiplyHelp)
	rpn.Register("*", Multiply, MultiplyHelp)
	rpn.Register("d", Divide, DivideHelp)
	rpn.Register("/", Divide, DivideHelp)
	rpn.Register("sq", Square, SquareHelp)
	rpn.Register("sqrt", SquareRoot, SquareRootHelp)
	rpn.Register("cos", Cos, CosHelp)
	rpn.Register("sin", Sin, SinHelp)
	rpn.Register("tan", Tan, TanHelp)
}
