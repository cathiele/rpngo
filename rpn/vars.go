package rpn

import (
	"fmt"
	"mattwach/rpngo/parse"
)

// Sets a variable
func (rpn *RPN) setVariable(name string) error {
	f, err := rpn.PopFrame()
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
	return rpn.PushFrame(f)
}

// Executes a Variables as a macro
func (rpn *RPN) execVariableAsMacro(name string) error {
	f, ok := rpn.Variables[name]
	if !ok {
		return fmt.Errorf("unknown variable: @%s", name)
	}
	if f.Type == COMPLEX_FRAME {
		// Just push the frame
		return rpn.PushFrame(f)
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
