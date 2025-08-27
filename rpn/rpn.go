package rpn

import (
	"fmt"
	"mattwach/rpngo/parse"
	"strconv"
	"strings"
)

// RPN is the main structure
type RPN struct {
	Stack     Stack
	Variables map[string]Frame
	functions map[string]func(*Stack) error
}

// Init initializes an RPNCalc object
func (rpn *RPN) Init() {
	rpn.Stack.Clear()
	rpn.functions = make(map[string]func(*Stack) error)
	rpn.Variables = make(map[string]Frame)
}

// Exec executes a single instruction
func (rpn *RPN) Exec(arg string) error {
	if fn := rpn.functions[arg]; fn != nil {
		return fn(&rpn.Stack)
	}
	if len(arg) > 1 {
		if arg[len(arg)-1] == '=' {
			return rpn.setVariable(arg[:len(arg)-1])
		}
		if arg[0] == '$' {
			return rpn.getVariable(arg[1:])
		}
		if arg[0] == '@' {
			return rpn.execVariableAsMacro(arg[1:])
		}
	}
	if len(arg) >= 2 {
		if (arg[0] == '"') && (arg[len(arg)-1] == '"') {
			return rpn.Stack.PushString(arg[1 : len(arg)-1])
		}
		if (arg[0] == '\'') && (arg[len(arg)-1] == '\'') {
			return rpn.Stack.PushString(arg[1 : len(arg)-1])
		}
	}
	return rpn.pushComplex(arg)
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *Stack) error) {
	rpn.functions[name] = fn
}

// Sets a variable
func (rpn *RPN) setVariable(name string) error {
	f, err := rpn.Stack.PopFrame()
	if err != nil {
		return err
	}
	rpn.Variables[name] = f
	return nil
}

// Gets a variable
func (rpn *RPN) getVariable(name string) error {
	f, ok := rpn.Variables[name]
	if !ok {
		return fmt.Errorf("unknown variable: $%s", name)
	}
	return rpn.Stack.PushFrame(f)
}

// Executes a Variables as a macro
func (rpn *RPN) execVariableAsMacro(name string) error {
	f, ok := rpn.Variables[name]
	if !ok {
		return fmt.Errorf("unknown variable: @%s", name)
	}
	if f.Type == COMPLEX_FRAME {
		// Just push the frame
		return rpn.Stack.PushFrame(f)
	}
	fields, err := parse.Fields(f.Str)
	if err != nil {
		return err
	}
	for _, f := range fields {
		if err := rpn.Exec(f); err != nil {
			return fmt.Errorf("@%s(%s): %v", name, f, err)
		}
	}
	return nil
}

// Pushes a float onto the stack
func (rpn *RPN) pushComplex(arg string) error {
	var v complex128
	var err error

	if strings.HasSuffix(arg, "i") {
		v, err = parseComplexWithI(arg)
		if err != nil {
			return err
		}
	} else {
		fv, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return err
		}
		v = complex(fv, 0)
	}
	return rpn.Stack.PushComplex(v)
}

// parses a complex string that contains an i
func parseComplexWithI(arg string) (complex128, error) {
	// a is the "real" part and b is the "imag" part: a + bi
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
		// no real part was given.  e.g. 5i
		a = "0"
		b = arg[:len(arg)-1]
	}
	b = strings.TrimPrefix(b, "+")
	if b == "" {
		// the user specified just i
		b = "1"
	} else if b == "-" {
		// the user specified just -i
		b = "-1"
	}
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	fb, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return complex(fa, fb), nil
}
