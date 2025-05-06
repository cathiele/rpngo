package rpn

import (
	"strconv"
)

// RPN is the main structure
type RPN struct {
	Stack     Stack
	functions map[string]func(*Stack) error
}

// Init initializes an RPNCalc object
func (rpn *RPN) Init() {
	rpn.Stack.Clear()
	rpn.functions = make(map[string]func(*Stack) error)
}

// Exec executes a single instruction
func (rpn *RPN) Exec(arg string) error {
	if fn := rpn.functions[arg]; fn != nil {
		return fn(&rpn.Stack)
	}
	return rpn.pushFloat(arg)
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *Stack) error) {
	rpn.functions[name] = fn
}

// Pushes a float onto the stack
func (rpn *RPN) pushFloat(arg string) error {
	v, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return err
	}
	return rpn.Stack.Push(Frame{Float: v})
}
