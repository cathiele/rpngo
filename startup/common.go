package startup

const commonStartup = `

# common index-out-of-bound code
"w.reset
'g1' w.new.group
'g1' w.columns
'i' 'g1' w.move.beg
'g2' w.new.group
'g2' 'g1' w.move.end
's1' w.new.stack
's1' 'g2' w.move.beg
'v1' w.new.var
'v1' 'g2' w.move.beg
" crash=

# set some useful vars
3.141592653589793 pi=
2.718281828459045 e=

# some useful equations
'vpush
 c= $0 neg bn= b= $0 2 * a2= a= $b sq 4 $a $c * * - sqrt root=
 $bn $root + $a2 /
 $bn $root - $a2 /
 vpop' quad=
`
