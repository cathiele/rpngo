//go:build pico || pico2

package startup

const defaultConfig = commonStartup + `
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
{ hists 'command history saved to' printsx 'i' 'histpath' w.getp printlnx } .f5=

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

true .echo=
`
