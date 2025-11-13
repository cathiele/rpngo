package rpn

import (
	"mattwach/rpngo/convert"
	"mattwach/rpngo/elog"
	"sort"
)

type AngleUnit int

const (
	ANGLE_RADIANS AngleUnit = iota
	ANGLE_DEGREES
	ANGLE_GRADS
)

// RPN is the main structure
type RPN struct {
	Frames    []Frame
	variables map[string][]Frame
	functions map[string]func(*RPN) error
	// maps are category -> command -> help
	help          map[string]map[string]string
	Interrupt     func() bool
	Print         func(string)
	Input         func(*RPN) (string, error)
	TextWidth     int
	maxStackDepth int
	AngleUnits    AngleUnit
	conv          *convert.Conversion
}

// Init initializes an RPNCalc object
func (r *RPN) Init(maxStackDepth int) {
	r.Clear()
	elog.Heap("alloc: /rpn/rpn.go:26: r.Frames = make([]Frame, 0, maxStackDepth)")
	r.Frames = make([]Frame, 0, 16) // object allocated on the heap: object size 10240 exceeds maximum stack allocation size 256
	r.maxStackDepth = maxStackDepth
	r.functions = make(map[string]func(*RPN) error)
	elog.Heap("alloc: /rpn/rpn.go:28: r.variables = []map[string]Frame{make(map[string]Frame)}")
	r.variables = make(map[string][]Frame) // object allocated on the heap: escapes at line 28
	r.conv = convert.Init()                // must come before initHelp()
	r.initHelp()
	r.Register("ssize", stackSize, CatStack, stackSizeHelp)
	r.Register("vlist", listVariables, CatVariables, listVariablesHelp)
	r.Print = DefaultPrint
	r.Interrupt = DefaultInterrupt
	r.AngleUnits = ANGLE_RADIANS
	r.TextWidth = 80
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *RPN) error, helpcat, helptxt string) {
	rpn.functions[name] = fn
	cat := rpn.help[helpcat]
	if cat == nil {
		rpn.help[helpcat] = map[string]string{name: helptxt}
	} else {
		cat[name] = helptxt
	}
}

func (rpn *RPN) AllFunctionNames() []string {
	var names []string
	for name := range rpn.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func DefaultPrint(msg string) {
	print(msg)
}

func DefaultInterrupt() bool {
	return false
}

func (r *RPN) Println(msg string) {
	r.Print(msg)
	r.Print("\n")
}


func (r *RPN) FromRadians(rad complex128, t Frame) Frame {
	switch r.AngleUnits {
	case ANGLE_DEGREES:
		return Frame{ftype: t.ftype, cmplx: rad * 57.29577951308232, str: "`deg"}
	case ANGLE_GRADS:
		return Frame{ftype: t.ftype, cmplx: rad * 63.66197723675813, str: "`grad"}
	default:
		return Frame{ftype: t.ftype, cmplx: rad, str: "`rad"}
	} 
}

func fromRadiansBase(rad complex128, u AngleUnit) complex128 {
	switch u {
	case ANGLE_DEGREES:
		return rad * 57.29577951308232
	case ANGLE_GRADS:
		return rad * 63.66197723675813
	default:
		return rad
	} 
}

func (r *RPN) ToRadians(angle complex128) complex128 {
	switch r.AngleUnits {
	case ANGLE_DEGREES:
		return angle * 0.0174532925199433
	case ANGLE_GRADS:
		return angle * 0.01570796326794897
	default:
		return angle
	} 
}

func toRadiansBase(angle complex128, u AngleUnit) complex128 {
	switch u {
	case ANGLE_DEGREES:
		return angle * 0.0174532925199433
	case ANGLE_GRADS:
		return angle * 0.01570796326794897
	default:
		return angle
	} 
}

