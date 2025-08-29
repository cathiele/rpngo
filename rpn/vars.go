package rpn

import (
	"errors"
	"fmt"
	"mattwach/rpngo/parse"
	"sort"
	"strings"
)

const pushVariableFrameHelp = "Pushes a variable frame to the variable stack"

func pushVariableFrame(r *RPN) error {
	r.variables = append(r.variables, make(map[string]Frame))
	return nil
}

const popVariableFrameHelp = "Pops a variable frame from the variable stack"

func popVariableFrame(r *RPN) error {
	if len(r.variables) <= 1 {
		return errors.New("no variable stack to pop")
	}
	r.variables = r.variables[:len(r.variables)-1]
	return nil
}

// Sets a variable
func (r *RPN) setVariable(name string) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.variables[len(r.variables)-1][name] = f
	return nil
}

// Gets a variable
func (r *RPN) getVariable(name string) (Frame, bool) {
	for i := len(r.variables) - 1; i >= 0; i-- {
		f, ok := r.variables[i][name]
		if ok {
			return f, true
		}
	}
	return Frame{}, false
}

// Gets all variable values as a string
func (r *RPN) getAllValuesForVariable(name string) string {
	var values []string
	lastVal := 0
	for i := 0; i < len(r.variables); i++ {
		f, ok := r.variables[i][name]
		if ok {
			values = append(values, f.String())
			lastVal = i
		} else {
			values = append(values, "nil")
		}
	}
	if len(values) == 0 {
		return "nil"
	}
	return strings.Join(values[:lastVal+1], " -> ")
}

// Gets all variable names
func (r *RPN) AllVariableNamesAndValues() []string {
	var names []string
	for i := 0; i < len(r.variables); i++ {
		for k := range r.variables[i] {
			names = append(names, k)
		}
	}
	if len(names) == 0 {
		return nil
	}
	sort.Strings(names)
	// names may contain duplicates
	var lastName string
	var results []string
	for _, name := range names {
		if name == lastName {
			continue
		}
		lastName = name
		results = append(
			results,
			fmt.Sprintf("%s: %s", name, r.getAllValuesForVariable(name)))
	}
	return results
}

// Executes a Variables as a macro
func (r *RPN) execVariableAsMacro(name string) error {
	f, ok := r.getVariable(name)
	if !ok {
		return fmt.Errorf("unknown variable: @%s", name)
	}
	if f.Type == COMPLEX_FRAME {
		// Just push the frame
		return r.PushFrame(f)
	}
	fields, err := parse.Fields(f.Str)
	if err != nil {
		return err
	}
	for _, f := range fields {
		if err := r.Exec(f); err != nil {
			return fmt.Errorf("@%s(%s): %v", name, f, err)
		}
	}
	return nil
}
