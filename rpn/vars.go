package rpn

import (
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
		return ErrStackEmpty
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
	if strings.Contains(name, "=") {
		return ErrIllegalName
	}
	if strings.Contains(name, "$") {
		return ErrIllegalName
	}
	r.variables[len(r.variables)-1][name] = f
	return nil
}

// Clears a variable
func (r *RPN) clearVariable(name string) error {
	vframe := r.variables[len(r.variables)-1]
	_, ok := vframe[name]
	if !ok {
		return ErrNotFound
	}
	delete(vframe, name)
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
		return "", ErrNotFound
	}
	if v.Type != STRING_FRAME {
		return "", ErrExpectedAString
	}
	return v.Str, nil
}

// gets a variable as a complex
func (r *RPN) GetComplexVariable(name string) (complex128, error) {
	v, ok := r.getVariable(name)
	if !ok {
		return 0, ErrNotFound
	}
	if v.Type != COMPLEX_FRAME {
		return 0, ErrExpectedANumber
	}
	return v.Complex, nil
}

// Gets all variable values as a string
func (r *RPN) getAllValuesForVariable(name string) []Frame {
	var values []Frame
	lastVal := 0
	for i := 0; i < len(r.variables); i++ {
		f, ok := r.variables[i][name]
		if ok {
			values = append(values, f)
			lastVal = i
		} else {
			values = append(values, Frame{Type: EMPTY_FRAME})
		}
	}
	if len(values) == 0 {
		return []Frame{{Type: EMPTY_FRAME}}
	}
	return values[:lastVal+1]
}

type NameAndValues struct {
	Name   string
	Values []Frame
}

// Gets all variable names
func (r *RPN) AllVariableNamesAndValues() []NameAndValues {
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
	var results []NameAndValues
	for _, name := range names {
		if name == lastName {
			continue
		}
		lastName = name
		results = append(results, NameAndValues{Name: name, Values: r.getAllValuesForVariable(name)})
	}
	return results
}

// Executes a Variables as a macro
func (r *RPN) execVariableAsMacro(name string) error {
	f, ok := r.getVariable(name)
	if !ok {
		return ErrNotFound
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
