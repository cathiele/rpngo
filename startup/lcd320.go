package startup

import (
	"fmt"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const lcd320ConfigFile = `
'root' w.columns
'v' w.new.var
's' w.new.stack
'g1' w.new.group
'g1' 50 w.weight
'v' 'g1' w.move.end
's' 'g1' w.move.end
'i' 'g1' w.move.end

# Plot defaults
'p1' plot.win=
'$plot.win w.new.plot' plot.init=
` + commonStartup

// LCD320Startup is startup logic when using a 320x240 display
// filesystem available)
func LCD320Startup(r *rpn.RPN) error {
	fields, err := parse.Fields(lcd320ConfigFile)
	if err != nil {
		return fmt.Errorf("while parsing lcd320ConfigFile var: %w", err)
	}
	if err := r.Exec(fields); err != nil {
		return fmt.Errorf("while executing commands in lcd320ConfigFile: %w", err)
	}
	return nil
}
