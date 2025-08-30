package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

// Exec executes a single instruction
func (rpn *RPN) exec(arg string) error {
	if fn := rpn.functions[arg]; fn != nil {
		return fn(rpn)
	}
	if len(arg) > 1 {
		if arg[len(arg)-1] == '=' {
			return rpn.setVariable(arg[:len(arg)-1])
		}
		if arg[0] == '$' {
			f, ok := rpn.getVariable(arg[1:])
			if !ok {
				return fmt.Errorf("variable not found: %s", arg[1:])
			}
			rpn.PushFrame(f)
			return nil
		}
		if arg[0] == '@' {
			return rpn.execVariableAsMacro(arg[1:])
		}
	}
	if len(arg) >= 2 {
		if (arg[0] == '"') && (arg[len(arg)-1] == '"') {
			return rpn.PushString(arg[1 : len(arg)-1])
		}
		if (arg[0] == '\'') && (arg[len(arg)-1] == '\'') {
			return rpn.PushString(arg[1 : len(arg)-1])
		}
	}
	if arg == "true" {
		return rpn.PushBool(true)
	}
	if arg == "false" {
		return rpn.PushBool(false)
	}
	return rpn.pushComplex(arg)
}

func (rpn *RPN) Exec(args []string) error {
	for _, arg := range args {
		if err := rpn.exec(arg); err != nil {
			return err
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
	return rpn.PushComplex(v)
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
