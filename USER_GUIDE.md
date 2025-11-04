## Build

You'll first need to [install golang](https://go.dev/doc/install).

There is a `bin/` directory that contains various different builds of `rpngo` for
the computer and for microcontrollers. Here is an overview:

- `bin/minimal/rpn` - The most basic version.  Parses `args` and exits.
- `bin/ncurses/rpn` - Uses [`ncurses`](https://en.wikipedia.org/wiki/Ncurses) to
  support multiple view windows and even text-based plotting
- `bin/tinygo/serialonly` - A stripped down build for  microcontrollers that uses
  serial communication only.
- `bin/tinygo/ili9341` - A TinyGo microcontroller build that still uses serial for input,
  and a ili9341 color LCD for output.  Full-featured.
- `bin/tinygo/picocalc` - A TinyGo microcontroller build that targets a
  [PicoCalc](https://www.clockworkpi.com/picocalc). Full-featured.


### Desktop / Raspberry Pi

Minimal version:

```
cd bin/minimal/rpn
go build
./rpn 2 3 +
```

ncurses version 

You need have `libncurses-dev` already installed.  In Ubuntu/Debian, the fix is:

```
sudo apt install libncurses-dev
```

then

```
cd bin/ncurses/rpn
go build
./rpn
```

### Microcontrollers using TinyGo (Raspberry Pi Pico and Pico2 tested as working)

You'll need to [install TinyGo](https://tinygo.org/getting-started/install/).

Minimal. Set your microcontroller target.  Tested on Pico and Pico2.  Other
chips may need configuration change.

```
cd bin/tinygo/serialonly
tinygo build -target=pico
```

ili9341 LCD, using serial for communication

```
cd bin/tinygo/ili9341
# check the command
make -n flash
# if looks ok, execute
make flash
```

PicoCalc

```
cd bin/tinygo/picocalc
# check the command
make -n flash
# if looks ok, execute
make flash
```


### Build Everything and Run Unit Tests

In the top-level directory, you can type

```
make all
```


## RPN Introduction

This "users guide" will take the format of being mostly working examples.
Let's start with the basics.

RPN is an old and proven way to do calculations that is popular in engineering
fields.  Much has been written on the subject. You can [read more here](
https://en.wikipedia.org/wiki/Reverse_Polish_notation)

Let's start with '2 + 3'


    2
    3
    +

or

    2 3 +

Both forms (spaces and new lines) do the following:

1. Push a `2` to a "stack"
2. Push a `3` to the same stack
3. Call the `+` operator which pulls 2 stack values (`2`, `3`) and pushes the result `5`

You can use the up arrow to browse command history. Command history might make
the `2 3 +` style more convenient than using separate lines.

There are some technical reasons I could use to argue for RPN (no parenthesis
needed, etc) but my main personal reason is that I saw my test scores in
college jump after switching because I made fewer mistakes than with my
previous calculator.  I carry that positive feeling with me to this day and
simply enjoy using them now.

## Editing Features

The experience is similar to most terminals:

- Left, right, insert, backspace, home, and end all work like you would expect
- Press up and down arrows to scroll through command history
- Ctrl-C to cancel a running program
- Ctrl-D to exit the program
- Press "ESC" or "page up" to enter scrolling mode where you can
  use page up, page down, up arrow and down arrow to view text
  that has scrolled off the top of the window. "ESC" or scrolling
  down far enough exits this mode.

Type `edit` to enter a full-panel multi-line editor.  The window
will contain the value at the top of the stack.  For example:

```
'animate_sin.rpn' load edit
```

The editor is currently only supports basic features:

- arrow key, page up, page down, home, end navigation
- insert, replace, backspace, delete
- syntax highlighting

While editing, press "ESC" to keep your changes (which will be
at the top of the stack) or "Ctrl-c" to exit without changing the
value.  If you happen to want to save your work permanently, exit with
"ESC", then type something like:

```
'my_file.txt' save
```

If you need more editing power, you can use a text editor on
your computer and copy over the file via serial or
an SD Card.

## Beyond the Stack: Variables

You can define and use variables.  Most build variants will define some
variables, such as `$pi`, on startup.

    5 a=
    2 $a +  # 7

Variables that start with a `.` are hidden by default in the variable window
(which we'll get to later). These often are used for special configuration
settings.

There are also special variables, `$0`, `$1`, etc, which represent values on
stack. `$0` represents the value at the top of the stack. If the stack
is empty, then using `$0` results in an error.  e.g.

    5 sq    # squares the number 5
    5 $0 *  # also squares the number 5

## Stack Shifts

You can move values around the stack.

    10 20 30   # put 10 20 30 on th stack
    1>         # now the stack is 10 30 20
    1<         # back to 10 20 30
    2>         # now the stack is 20 30 10
    2<         # back to 10 20 30
    $2 2<      # now it's 10 10 20 30

## Stack Deletions

    10 20 30   # put 10 20 30 on the stack
    1/         # now the stack is 10 30
    0/         # now the stack is 10
    d          # emptys all values from the stack

## Number Bases

Many type of numbers are supported

    50       # floating point (actually a 5+0i complex)
    50+i     # complex number
    50d      # Integer
    32x      # Hexidecimal
    62o      # octal
    110010b  # binary

Most operations can us a mix of these types, using the following rules:

    # Any number type mixed with float results in a float
    12.4 5d +  ->  17.4

    # Two integer types combined takes the base of the left term
    32x 50d +  ->  64x

You can also convert between types using `hex`, `bin`, `oct`, `float`, `real`,
`imag`, and `str`. You can convert from a string to a type by executing
it with `@`

    "54x" @  ->  54x

## Booleans and Conditionals

Boolean values include `true` and `false`.  Conditionals return a boolean:

    true       #  true
    false      #  false
    1 2 >      #  false
    1 2 <      #  true
    1 2 =      #  false
    1 2 !=     #  true
    1 2 >=     #  false
    1 2 <=     #  false
    false neg  #  true
    true neg   #  false

Conditionals are an essential part of programming, which we will cover
with examples later.

## Strings

There are three ways to specify a string

    "hello world"
    'hello world'
    {hello world}

All three examples define the same string.  The third form can be useful
when you want to nest terms in a program.

    {0 x= {$x 1 + x= $x println 1000 <} for} count=

It could also be done with other string types, rpngo will execute them
the same way:

    '0 x= "$x 1 + x= $x println 1000 <" for' count=

Which do you find easier to read?

## Macros

You can define a string as a macro, here is one for the area of a circle (`$pi
* r * r / 2`), given the radius:

    {sq $pi * 2 /} carea=
    5 @carea -> 39.26990817

A macro can also contain loops, conditionals and calls to other macros (needed
for programming). 

## Labels

Labels can be added to non string values.  A label shows up in the stack window
and can be used to annotate what the number represents.  This can be
useful when attempting to do complex calculations purely using the stack.
It cal also be used to communicate what values are (The `heapstats` function
uses this)

For example, say you want to make a formula that converts velocity, time, 
and acceleration into distance. The formula is `(v * t) + (0.5 * a * t * t)`.
You could use variables or the pure stack. If you use the pure stack, labels
can help keep track of what is what:

     1 `v
     2 `a
     3 `t

Now the stack will look like this:

     2: 1 `v
     1: 2 `a
     0: 3 `t

You can see that `$1` is `a` and see how `a` changes stack slots as you
work the equation.  This makes it easier to create a macro that does
not define any variables (which can be fun little puzzles to solve).

    {$0 sq 2> * 2 / 2< * +} dist=
    1 2 3 @dist  ->  12
   
## Conditionals

`if` and `ifelse` can be used to conditionally execute a bit of code
based on the result of a `true`/`false` condition.

```
> true {'yes' printlnx} if
yes

> false {'yes' printlnx} if

> 5 1 > {'is greater' printlnx} if
is greater

> true {'yes'} {'no'} ifelse printlnx
yes

> false {'yes'} {'no'} ifelse printlnx
no
```

## Unit Conversions

You can convert between several unit types, some examples:

    5 km>mi
    3.106855961

    60 mi/h>m/s
    26.8224

    10 liter>m*m*m
    0.01

    1 megabyte>bits
    8388608

See all possible conversions with

    conversions?


## Window Layout

The `rpngo` program supports several window types:

- Input
- Stack
- Variables
- Plot

Of the above, only the input is required.  For the others,
you can have zero or more of them.  For example, if you
want an input window and two separate plot windows, you
can do it. You might also have two stack windows with
different configuration options set.  Whatever you want.

These options need a means to create windows and put
them in the position and size you want.

For example, say we want the following:

```
+------------------+---------------+
|                  |               |
|                  |               |
|                  |               |
|   Input (i)      |  STACK (s)    |
|                  |               |
|                  |               |
+------------------+               |
|                  |               |
|                  |               |
|   Vars (v)       |               |
|                  |               |
|                  |               |
+------------------+---------------+
```

How do we do it?

The `rpngo` window system works a bit like html tables. You have a 'root'
window (like a table) that can either contain `1xn` or `nx1` children.
Each child can either be a window (input, stack, etc) or a window
group with the same rules as the root group.

Let get on with how you would create the layout above.  First, lets reset
everything:

    w.reset

You'll have this default starting point:

```
+----------------------------------+
|                                  |
|                                  |
|                                  |
|                                  |
|                                  |
|                                  |
|         Input (i)                |
|                                  |
|                                  |
|                                  |
|                                  |
|                                  |
+----------------------------------+
```

Now lets add a window group names g

    'g' w.new.group 

and we'll have

```
+----------------------------------+
|                                  |
|                                  |
|         Input (i)                |
|                                  |
|                                  |
|                                  |
+----------------------------------+
|                                  |
|                                  |
|         Group (g)                |
|                                  |
|                                  |
+----------------------------------+
```

Let's switch the root window to column layout

    'root' w.columns

```
+----------------+------------------+
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|   Input (i)    |    Group (g)     |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
+----------------+------------------+
```

and move input into the group

    'i' 'g' w.move.beg

```
+----------------------------------+
|                                  |
|                                  |
|                                  |
|                                  |
|                                  |
|                                  |
|         Input (i)                |
|         Group (g)                |
|                                  |
|                                  |
|                                  |
|                                  |
+----------------------------------+
```

Let's add the stack window:

    's' w.new.stack

```
+----------------+------------------+
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|   Input (i)    |    Stack (s)     |
|   Group (g)    |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
|                |                  |
+----------------+------------------+
```

and finally the vars window. We have two options here. We can either add the vars
window to the root window, then move it to `g` or create it in `g` to begin with.
If we just say:

    $.warget

It will say

    'root'

We can change this "special" variable to alter the behavior of `w.new.*` commands.
Note there is also `w.beg` and `w.weight` for controlling placement and size.
Let's not worry about that yet:

    'g' .warget=
    'v' w.new.var

and we are done

```
+------------------+---------------+
|                  |               |
|                  |               |
|                  |               |
|   Input (i)      |  STACK (s)    |
|                  |               |
|                  |               |
+------------------+               |
|                  |               |
|                  |               |
|   Vars (v)       |               |
|                  |               |
|                  |               |
+------------------+---------------+
```

OK, one more things.  Say we want to change some sizes above.
Window and group "weights" are used to do this. Each window and group
above is assigned a default weight of 100.  Weights are shared
amongst siblings in a group to decide how much space each one gets.

For example, ley's change the weight of window group 'g' from 100
to 200.

    'g' 200 w.weight

Things will shift as so:

```
+----------------------+-----------+
|                      |           |
|                      |           |
|                      |           |
|   Input (i)          | STACK (s) |
|                      |           |
|                      |           |
+----------------------+           |
|                      |           |
|                      |           |
|   Vars (v)           |           |
|                      |           |
|                      |           |
+----------------------+-----------+
```

A "real example":

![window layout](img/window_layout.png)

There is a command, `w.dump` that will output the current layout for
inspection. Let's try it:

    w.dump

    root(x=0, y=0, w=80, h=30, cols=true, weight=100):
      g(x=0, y=0, w=53, h=30, cols=false, weight=200):
        i(type=input, x=1, y=1, w=51, h=13, weight=100)
        v(type=var, x=1, y=15, w=51, h=14, weight=100)
      s(type=stack, x=53, y=1, w=26, h=28, weight=100)

    .wtarget=g .wend=true .wweight=100d
 
We might want to set `$.wtarget` back (or push/pop it with `.wtarget<` and `.wtarget>`,
but note that `w.reset` also sets `$.wtarget` back to `root` and most window
management macros start with `w.reset` anyway.

## Programmable Function Keys

Configuring windows like above is flexible but cumbersome.  For everyday use,
there are variables like `.f1`, `.f2` that invoke when you
press the corresponding F key.  I like to set these to change window
layouts.  For example:

```
# Input and stack window
{
  w.reset
  's' w.new.stack
  's' 30 w.weight
} .f1=

# Just one big input window
{w.reset} .f2=
```

and so on.

## Window Properties

Each of the four window types (input, stack, variable, plot) have customizable
properties.  We'll start with common commands, then explore details (expect
for plot, which we cover later).

### Input Window Properties

You can list properties for a window with `w.listp`

    'i' w.listp

       autofn: {}
       showframes: 1d

These are what the input properties do:

- `autofn` Executes the given code before showing a prompt.  This can be used to
  track memory usage, print a custom message, update a graph, or anything you want.
- `showframes` How much of the stack to copy to the input window after each
  entered command. More or less stack in the history has it's trade-offs and is
  personal preference.

Here is an example of setting properties:

    0 cmd=
    'i' 'autofn' {$cmd 1 + print cmd=} w.setp

Now w have an `autocmd` property set in the input window that will print an
incrementing command number every time you press enter.  Only one more
property command to cover:

    'i' 'autofn' w.getp

This copies the property value to the stack. It can be useful if you
want to use a property value in a program.

### Stack Window Properties

    's' w.listp

      round: -1d

For now, the only property is `round` which allow you to round
floating point numbers to a given number of decimal places (-1
represents no rounding).

### Variable Window Properties

    'v' w.listp

      showdot: false
      multiline: false

- `showdot`: If true, then variable names that start with a `.` (such as
  `.wtarget`) will also be shown.
- `multiline`: If true, then string that expand multiple lines will
  show all of these lines in the window (assuming there is space).


