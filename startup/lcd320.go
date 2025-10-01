package startup

import (
	"fmt"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const lcd320ConfigFile = `
's' w.new.stack
's' 30 w.weight

# Plot defaults
'p' plot.win=
'w.reset $plot.win w.new.plot $plot.win "root" w.move.beg $plot.win $plot.win 200 w.weight' plot.init=
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
