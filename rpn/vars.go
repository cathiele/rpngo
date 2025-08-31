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

// gets a variable as a string
func (r *RPN) GetStringVariable(name string) (string, error) {
	v, ok := r.getVariable(name)
	if !ok {
		return "", fmt.Errorf("$%s is not defined", name)
	}
	if v.Type != STRING_FRAME {
		return "", fmt.Errorf("$%s is not a string", name)
	}
	return v.Str, nil
}

// gets a variable as a complex
func (r *RPN) GetComplexVariable(name string) (complex128, error) {
	v, ok := r.getVariable(name)
	if !ok {
		return 0, fmt.Errorf("$%s is not defined", name)
	}
	if v.Type != COMPLEX_FRAME {
		return 0, fmt.Errorf("$%s is not a number", name)
	}
	return v.Complex, nil
}

// Gets all variable values as a string
func (r *RPN) getAllValuesForVariable(name string) string {
	var values []string
	lastVal := 0
	for i := 0; i < len(r.variables); i++ {
		f, ok := r.variables[i][name]
		if ok {
			values = append(values, f.String(true))
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
	if err := r.Exec(fields); err != nil {
		return err
	}
	return nil
}
