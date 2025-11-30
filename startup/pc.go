//go:build !pico && !pico2

package startup

const defaultConfig = commonStartup + `
# Create and layout windows
{
  'g' w.new.group
  'g' w.columns
  'i' 'g' w.move.end
  'i' 25 w.weight
  'g2' w.new.group
  'g2' 25 w.weight
  'g2' w.columns
  'g2' .wtarget=
  's' w.new.stack
  'v' w.new.var
  'g' .wtarget=
} .init=

@.init

'w.reset @.init' .f1=
'w.reset "root" w.columns "i" 30 w.weight' .f2=
'w.reset "root" w.columns "v" w.new.var "v" "showdot" true w.setp' .f3=

'time t1= 0 x= "$x 1 + x= $x 3000000 <" for time $t1 - 3000000 1> /' benchmark=

# Plot defaults
'p' .plotwin=
'$.plotwin w.new.plot' .plotinit=

'/dev/ttyACM0' .serial=
`
