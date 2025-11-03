package startup

import (
	"fmt"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const lcd320ConfigFile = `
{
  30 .wweight<
  's' w.new.stack
  .wweight> 0/
} .init=

@.init

{w.reset @.init} .f1=
{w.reset} .f2=
{
  w.reset
  'root' w.columns
  'v' w.new.var
} .f3=
{@.plotinit} .f4=
{
	w.reset
	'g' w.new.group
	'g' w.columns
	'i' 'g' w.move.beg
	'g' .wtarget<
	'v' w.new.var
	.wtarget> 0/
	30 .wweight=
	's' w.new.stack
	.wweight> 0/
} .f5=

'i' 'autofn'
{
  heapstats
  2> 1024d / float 'a:' printx printx 'k f:'
  printx float printx
  0/ 0/
} w.setp 

{
  time t1=
  0 x=
  {$x 1 + x= $x 50000 <} for
  time $t1 - 50000 1> /
  ' loops/second' + printlnx
  heapstats
} benchmark=

# Plot defaults
'p' .plotwin=
{
  w.reset
  false .wend<
  250 .wweight<
  $.plotwin w.new.plot
  .wend> 0/ .wweight> 0/
} .plotinit=
` + commonStartup

// LCD320Startup is startup logic when using a 320x240 display
// filesystem available)
func LCD320Startup(r *rpn.RPN) error {
	err := parse.Fields(lcd320ConfigFile, r.Exec)
	if err != nil {
		return fmt.Errorf("while parsing lcd320ConfigFile var: %w", err)
	}
	return nil
}
