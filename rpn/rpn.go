package rpn

import (
	"mattwach/rpngo/convert"
	"sort"
)

// RPN is the main structure
type RPN struct {
	Frames    []Frame
	pushed    [][]Frame
	variables []map[string]Frame
	functions map[string]func(*RPN) error
	// maps are category -> command -> help
	help      map[string]map[string]string
	Interrupt func() bool
	Print     func(string)
	Input     func(*RPN) (string, error)
	TextWidth int
	conv      *convert.Conversion
}

// Init initializes an RPNCalc object
func (r *RPN) Init(maxStackDepth int) {
	r.Clear()
	r.Frames = make([]Frame, 0, maxStackDepth) // object allocated on the heap: (OK)
	r.functions = make(map[string]func(*RPN) error)
	r.variables = []map[string]Frame{make(map[string]Frame)} // object allocated on the heap (OK)
	r.initHelp()
	r.Register("ssize", stackSize, CatStack, stackSizeHelp)
	r.Register("spush", pushStack, CatStack, pushStackHelp)
	r.Register("spop", popStack, CatStack, popStackHelp)
	r.Register("vlist", listVariables, CatVariables, listVariablesHelp)
	r.Register("vpush", pushVariableFrame, CatVariables, pushVariableFrameHelp)
	r.Register("vpop", popVariableFrame, CatVariables, popVariableFrameHelp)
	r.Print = DefaultPrint
	r.Interrupt = DefaultInterrupt
	r.conv = convert.Init()
	r.TextWidth = 80
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *RPN) error, helpcat, helptxt string) {
	rpn.functions[name] = fn
	cat := rpn.help[helpcat]
	if cat == nil {
		rpn.help[helpcat] = map[string]string{name: helptxt}
	} else {
		cat[name] = helptxt
	}
}

func (rpn *RPN) AllFunctionNames() []string {
	var names []string
	for name := range rpn.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func DefaultPrint(msg string) {
	print(msg)
}

func DefaultInterrupt() bool {
	return false
}

func (r *RPN) Println(msg string) {
	r.Print(msg)
	r.Print("\n")
}
