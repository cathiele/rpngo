package rpn

import (
	"mattwach/rpngo/parse"
	"sort"
	"strconv"
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

const listVariablesHelp = "Prints all variable names"

func listVariables(r *RPN) error {
	names := r.AllVariableNames()
	for _, n := range names {
		r.Print(n)
		r.Print("\n")
	}
	return nil
}

// Sets a variable
func (r *RPN) SetVariable(name string) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if err := checkVariableName(name); err != nil {
		r.PushFrame(f)
		return err
	}
	r.variables[len(r.variables)-1][name] = f
	return nil
}

func checkVariableName(name string) error {
	if len(name) == 0 {
		return ErrIllegalName
	}
	if !isAlpha(rune(name[0])) {
		return ErrIllegalName
	}
	for _, r := range name {
		if !isAlphaNum(r) {
			return ErrIllegalName
		}
	}
	return nil
}

func isAlpha(r rune) bool {
	return (r == '.') || (r == '_') || ((r >= 'A') && (r <= 'Z')) || ((r >= 'a') && (r <= 'z'))
}

func isNum(r rune) bool {
	return (r >= '0') && (r <= '9')
}

func isAlphaNum(r rune) bool {
	return (r == '.') || (r == '_') || ((r >= '0') && (r <= '9')) || ((r >= 'A') && (r <= 'Z')) || ((r >= 'a') && (r <= 'z'))
}

// Clears a variable
func (r *RPN) ClearVariable(name string) error {
	if len(name) == 0 {
		return ErrIllegalName
	}
	if isNum(rune(name[0])) {
		return r.clearStackVariable(name)
	}
	vframe := r.variables[len(r.variables)-1]
	_, ok := vframe[name]
	if !ok {
		return ErrNotFound
	}
	delete(vframe, name)
	return nil
}

// Gets a variable
func (r *RPN) GetVariable(name string) (Frame, error) {
	if len(name) == 0 {
		return Frame{}, ErrIllegalName
	}
	if isNum(rune(name[0])) {
		return r.getStackVariable(name)
	}
	for i := len(r.variables) - 1; i >= 0; i-- {
		f, ok := r.variables[i][name]
		if ok {
			return f, nil
		}
	}
	return Frame{}, ErrNotFound
}

// Gets a variable from the stack
func (r *RPN) getStackVariable(name string) (Frame, error) {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return Frame{}, ErrIllegalName
	}
	return r.PeekFrame(idx)
}

// Removes a variable from the stack
func (r *RPN) clearStackVariable(name string) error {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return ErrIllegalName
	}
	_, err = r.DeleteFrame(idx)
	return err
}

func (r *RPN) moveStackVariableToHead(name string) error {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return ErrIllegalName
	}
	f, err := r.DeleteFrame(idx)
	if err != nil {
		return err
	}
	return r.PushFrame(f)
}

func (r *RPN) moveHeadStackVariable(name string) error {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return ErrIllegalName
	}
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	err = r.InsertFrame(f, idx)
	if err != nil {
		r.PushFrame(f)
	}
	return err
}

// gets a variable as a string
func (r *RPN) GetStringVariable(name string) (string, error) {
	v, err := r.GetVariable(name)
	if err != nil {
		return "", err
	}
	if v.Type != STRING_FRAME {
		return "", ErrExpectedAString
	}
	return v.Str, nil
}

// gets a variable as a complex
func (r *RPN) GetComplexVariable(name string) (complex128, error) {
	v, err := r.GetVariable(name)
	if err != nil {
		return 0, err
	}
	if v.Type == COMPLEX_FRAME {
		return v.Complex, nil
	}
	if v.IsInt() {
		return complex(float64(v.Int), 0), nil
	}
	return 0, ErrExpectedANumber
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
func (r *RPN) AllVariableNames() []string {
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
	return names
}

// Gets all variable names
func (r *RPN) AllVariableNamesAndValues() []NameAndValues {
	names := r.AllVariableNames()
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
	f, err := r.GetVariable(name)
	if err != nil {
		return err
	}
	if f.Type == COMPLEX_FRAME {
		// Just push the frame
		return r.PushFrame(f)
	}
	// this call can be recursive so we need to allocate here
	fields := make([]string, 16)
	fields, err = parse.Fields(f.Str, fields)
	if err != nil {
		return err
	}
	if err := r.Exec(fields); err != nil {
		return err
	}
	return nil
}
