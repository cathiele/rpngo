package functions

import "mattwach/rpngo/rpn"

// RegisterAll resisters functions that are core to the RPN engine. There
// are other functions, such as window-specific functions, that will be
// added by their respective owners as the functions module should not
// know about these.
func RegisterAll(rpn *rpn.RPN) {
	rpn.Register("-", Subtract, SubtractHelp)
	rpn.Register(".", Duplicate, DuplicateHelp)
	rpn.Register("*", Multiply, MultiplyHelp)
	rpn.Register("/", Divide, DivideHelp)
	rpn.Register("+", Add, AddHelp)
	rpn.Register("a", Add, AddHelp)
	rpn.Register("cos", Cos, CosHelp)
	rpn.Register("d", Divide, DivideHelp)
	rpn.Register("load", Load, LoadHelp)
	rpn.Register("m", Multiply, MultiplyHelp)
	rpn.Register("s", Subtract, SubtractHelp)
	rpn.Register("sin", Sin, SinHelp)
	rpn.Register("sq", Square, SquareHelp)
	rpn.Register("sqrt", SquareRoot, SquareRootHelp)
	rpn.Register("sw", Swap, SwapHelp)
	rpn.Register("tan", Tan, TanHelp)
	rpn.Register("x", Drop, DropHelp)
}
