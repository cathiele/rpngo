package startup

import (
	"fmt"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const lcd320ConfigFile = `
"
's' w.new.stack
's' 30 w.weight
" .init=

@.init

'w.reset @.init' .f1=

'time t1= 0 x= "$x 1 + x= $x 50000 <" for time $t1 - 50000 1> /' benchmark=

# Plot defaults
'p' .plotwin=
'w.reset $.plotwin w.new.plot $.plotwin "root" w.move.beg $.plotwin $.plotwin 200 w.weight' .plotinit=
` + commonStartup

// LCD320Startup is startup logic when using a 320x240 display
// filesystem available)
func LCD320Startup(r *rpn.RPN) error {
	fields := make([]string, 256)  // object allocated on the heap: object size 2048 exceeds maximum stack allocation size 256
	fields, err := parse.Fields(lcd320ConfigFile, fields)
	if err != nil {
		return fmt.Errorf("while parsing lcd320ConfigFile var: %w", err)
	}
	if err := r.Exec(fields); err != nil {
		return fmt.Errorf("while executing commands in lcd320ConfigFile: %w", err)
	}
	return nil
}