package rpn

import (
	"errors"
	"fmt"
)

var (
	errExpectedANumber      = errors.New("expected a number")
	errExpectedAString      = errors.New("expected a string")
	ErrStackEmpty           = errors.New("stack empty")
	errNotEnoughStackFrames = errors.New("not enough stack frames")
)

type FrameType uint8

const (
	STRING_FRAME FrameType = iota
	COMPLEX_FRAME
)

// Frame Defines a single stack frame
type Frame struct {
	Type    FrameType
	Str     string
	Complex complex128
}

// RPN is the main structure
type RPN struct {
	frames      []Frame
	variables   []map[string]Frame
	functions   map[string]func(*RPN) error
	commandHelp map[string]string
	conceptHelp map[string]string
	Print       func(string)
}

// Init initializes an RPNCalc object
func (r *RPN) Init() {
	r.Clear()
	r.functions = make(map[string]func(*RPN) error)
	r.variables = []map[string]Frame{make(map[string]Frame)}
	r.initHelp()
	r.Register("vpush", pushVariableFrame, pushVariableFrameHelp)
	r.Register("vpop", popVariableFrame, popVariableFrameHelp)
	r.Print = DefaultPrint
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *RPN) error, help string) {
	rpn.functions[name] = fn
	rpn.commandHelp[name] = help
}

func DefaultPrint(msg string) {
	fmt.Print(msg)
}

func (r *RPN) Println(msg string) {
	r.Print(msg)
	r.Print("\n")
}
