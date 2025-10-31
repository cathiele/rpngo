package startup

const commonStartup = `
# set some useful vars
3.141592653589793 pi=
2.718281828459045 e=

# some useful equations
# (-b +/- sqrt(b*b - 4*a*c)) / (2 * a)
{$2 * 4 * $1 sq - neg sqrt 1> neg $0 $2 - $3 2 * / 3< + 1> 2 * /} quad=

{0 {+} 1 filtern} sum=
{$0 {$1 $1 < {0/} {1/} ifelse} 1 filtern} min=
{$0 {$1 $1 > {0/} {1/} ifelse} 1 filtern} max=
`
