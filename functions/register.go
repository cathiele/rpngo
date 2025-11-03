package functions

import "mattwach/rpngo/rpn"

// RegisterAll resisters functions that are core to the RPN engine. There
// are other functions, such as window-specific functions, that will be
// added by their respective owners as the functions module should not
// know about these.
func RegisterAll(r *rpn.RPN) {
	r.Register("!=", NotEqual, rpn.CatCompare, NotEqualHelp)
	r.Register("<", LessThan, rpn.CatCompare, LessThanHelp)
	r.Register("<=", LessThanEqual, rpn.CatCompare, LessThanEqualHelp)
	r.Register("=", Equal, rpn.CatCompare, EqualHelp)
	r.Register(">", GreaterThan, rpn.CatCompare, GreaterThanHelp)
	r.Register(">=", GreaterThanEqual, rpn.CatCompare, GreaterThanEqualHelp)

	r.Register("&", And, rpn.CatBitwise, AndHelp)
	r.Register("|", Or, rpn.CatBitwise, OrHelp)
	r.Register("^", XOr, rpn.CatBitwise, XOrHelp)
	r.Register("<<", ShiftLeft, rpn.CatBitwise, ShiftLeftHelp)
	r.Register(">>", ShiftRight, rpn.CatBitwise, ShiftRightHelp)

	r.Register("-", Subtract, rpn.CatCore, SubtractHelp)
	r.Register("*", Multiply, rpn.CatCore, MultiplyHelp)
	r.Register("/", Divide, rpn.CatCore, DivideHelp)
	r.Register("+", Add, rpn.CatCore, AddHelp)
	r.Register("false", False, rpn.CatCore, FalseHelp)
	r.Register("neg", Negate, rpn.CatCore, NegateHelp)
	r.Register("round", Round, rpn.CatCore, RoundHelp)
	r.Register("true", True, rpn.CatCore, TrueHelp)

	r.Register("**", Power, rpn.CatEng, PowerHelp)
	r.Register("abs", Abs, rpn.CatEng, AbsHelp)
	r.Register("cos", Cos, rpn.CatEng, CosHelp)
	r.Register("log", Log, rpn.CatEng, LogHelp)
	r.Register("log10", Log10, rpn.CatEng, Log10Help)
	r.Register("rand", Rand, rpn.CatEng, RandHelp)
	r.Register("sin", Sin, rpn.CatEng, SinHelp)
	r.Register("sq", Square, rpn.CatEng, SquareHelp)
	r.Register("sqrt", SquareRoot, rpn.CatEng, SquareRootHelp)
	r.Register("tan", Tan, rpn.CatEng, TanHelp)

	r.Register("print", Print, rpn.CatIO, PrintHelp)
	r.Register("printall", PrintAll, rpn.CatIO, PrintAllHelp)
	r.Register("println", Println, rpn.CatIO, PrintlnHelp)
	r.Register("printlnx", PrintlnX, rpn.CatIO, PrintlnXHelp)
	r.Register("prints", PrintS, rpn.CatIO, PrintSHelp)
	r.Register("printsx", PrintSX, rpn.CatIO, PrintSXHelp)
	r.Register("printx", PrintX, rpn.CatIO, PrintXHelp)

	r.Register("@", Exec, rpn.CatProg, ExecHelp)
	r.Register("delay", Delay, rpn.CatProg, DelayHelp)
	r.Register("filter", Filter, rpn.CatProg, FilterHelp)
	r.Register("filterm", FilterM, rpn.CatProg, FilterMHelp)
	r.Register("filtermn", FilterMN, rpn.CatProg, FilterMNHelp)
	r.Register("filtern", FilterN, rpn.CatProg, FilterNHelp)
	r.Register("for", For, rpn.CatProg, ForHelp)
	r.Register("if", If, rpn.CatProg, IfHelp)
	r.Register("ifelse", IfElse, rpn.CatProg, IfElseHelp)
	r.Register("input", Input, rpn.CatProg, InputHelp)
	r.Register("noop", NoOp, rpn.CatProg, NoOpHelp)
	r.Register("fields", Fields, rpn.CatProg, FieldsHelp)
	r.Register("time", Time, rpn.CatProg, TimeHelp)

	r.Register("X", DropAll, rpn.CatStack, DropAllHelp)

	r.Register("bin", Bin, rpn.CatType, BinHelp)
	r.Register("float", Float, rpn.CatType, FloatHelp)
	r.Register("hex", Hex, rpn.CatType, HexHelp)
	r.Register("imag", Imag, rpn.CatType, ImagHelp)
	r.Register("int", Int, rpn.CatType, IntHelp)
	r.Register("oct", Oct, rpn.CatType, OctHelp)
	r.Register("real", Real, rpn.CatType, RealHelp)
	r.Register("str", Str, rpn.CatType, StrHelp)
}
