package rpn

import (
	"strconv"
	"strings"
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
	if (len(arg) >= 2) && (arg[0] == '"') && (arg[len(arg)-1] == '"') {
		return rpn.Stack.PushString(arg[1 : len(arg)-1])
	}
	if (len(arg) >= 2) && (arg[0] == '\'') && (arg[len(arg)-1] == '\'') {
		return rpn.Stack.PushString(arg[1 : len(arg)-1])
	}
	return rpn.pushComplex(arg)
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *Stack) error) {
	rpn.functions[name] = fn
}

// Pushes a float onto the stack
func (rpn *RPN) pushComplex(arg string) error {
	var v complex128

	if strings.HasSuffix(arg, "i") {
		var a string
		var b string
		for i, c := range arg {
			if c == '+' || (c == '-' && i > 0) {
				a = arg[:i]
				b = arg[i : len(arg)-1]
				break
			}
		}
		if a == "" {
			a = "0"
			b = arg[:len(arg)-1]
		}
		b = strings.TrimPrefix(b, "+")
		if b == "" {
			b = "1"
		} else if b == "-" {
			b = "-1"
		}
		fa, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return err
		}
		fb, err := strconv.ParseFloat(b, 64)
		if err != nil {
			return err
		}
		v = complex(fa, fb)
	} else {
		fv, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return err
		}
		v = complex(fv, 0)
	}
	return rpn.Stack.PushComplex(v)
}
