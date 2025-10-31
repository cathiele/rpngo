package startup

const commonStartup = `
# set some useful vars
3.141592653589793 pi=
2.718281828459045 e=

# some useful equations
{vpush
 c= $0 neg bn= b= $0 2 * a2= a= $b sq 4 $a $c * * - sqrt root=
 $bn $root + $a2 /
 $bn $root - $a2 /
 vpop} quad=

{0 {+} 1 filtern} sum=
{$0 {$1 $1 < {0/} {1/} ifelse} 1 filtern} min=
{$0 {$1 $1 > {0/} {1/} ifelse} 1 filtern} max=
`
