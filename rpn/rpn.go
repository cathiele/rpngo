package rpn

import "errors"

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
	messages    []string
	Variables   map[string]Frame
	functions   map[string]func(*RPN) error
	commandHelp map[string]string
	conceptHelp map[string]string
}

// Init initializes an RPNCalc object
func (rpn *RPN) Init() {
	rpn.Clear()
	rpn.functions = make(map[string]func(*RPN) error)
	rpn.Variables = make(map[string]Frame)
	rpn.initHelp()
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *RPN) error, help string) {
	rpn.functions[name] = fn
	rpn.commandHelp[name] = help
}
