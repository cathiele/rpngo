package startup

const commonStartup = `
# set some useful vars
3.141592653589793 pi=
2.718281828459045 e=

# some useful equations
# (-b +/- sqrt(b*b - 4*a*c)) / (2 * a)
{$2 * 4 * $1 sq - neg sqrt 1> neg $0 $2 - $3 2 * / 3< + 1> 2 * /} quad=

{0 {+ ssize 1 >} for} sum=
{ssize n< 0 {+ ssize 1 >} for n> /} mean=
{$0 {min ssize 1 >} for} min=
{$0 {max ssize 1 >} for} max=

histl
hists
`
