package rpn

import (
	"errors"
	"fmt"
)

const MaxStackDepth = 4096

var (
	ErrExpectedABoolean     = errors.New("expected a boolean")
	ErrExpectedANumber      = errors.New("expected a number")
	ErrExpectedAString      = errors.New("expected a string")
	ErrStackEmpty           = errors.New("stack empty")
	ErrStackFull            = errors.New("stack is full")
	errNotEnoughStackFrames = errors.New("not enough stack frames")
)

type FrameType uint8

const (
	STRING_FRAME FrameType = iota
	COMPLEX_FRAME
	BOOL_FRAME
)

// Frame Defines a single stack frame
type Frame struct {
	Type    FrameType
	Str     string
	Complex complex128
	Bool    bool
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
	r.addDefaultPlotVars()
}

func (r *RPN) addDefaultPlotVars() {
	// Set the default plot window to p1
	r.PushString("p1")
	r.setVariable("plot.win")
	// set the default plot init function
	r.PushString("$plot.win w.new.plot $plot.win 'root' w.move.beg $plot.win 200 w.weight")
	r.setVariable("plot.init")
	r.pushComplex("-1")
	r.setVariable("plot.min")
	r.pushComplex("1")
	r.setVariable("plot.max")
	r.pushComplex("400")
	r.setVariable("plot.steps")
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
