#RPNGO

## TODO Bugs

## TODO Cleanup

## TODO Features

- L: Basic text editor
- M: Add stack output formatting options
- L: Rewrite readme doc
- M: Experiement with SD card or on-chip storage
- S: Develop options for transferring data over serial 
- M: Improve font
- M: Create snapshot
- M: Foreach, del

## Wishlist

- Look into using picocalc memory layout for bootloader compatibility

This library is intended to serve as a backend engine and implements a
programmable RPN calculator. Included is a simple commandline frontend for
getting a taste of what the program can do.

Base features include:

- An RPN compute engine
- Macro support
- Units support (m/s, A, kg, etc)
- Complex Numbers, Float, Integers, Hexidecimal, Octal, Binary
- Conditionals support (for writing simple programs)
- Support for adding your own functionality

Core functions that are available (can be all be added, or selectively added):

- Arithimetic: +, -, \*, /, **
- Scientific: sqrt, sin, cos, tan, asin, acos, atan

## Build

You'll first need to install golang.

There is a `bin/` directory thast contains various different builds of rpn
for the computer and for microtrollers.  Try the ones you like and ignore the rest/  Here is an overview:

- bin/minimal/rpn - The most basic version.  Parses args and exsts.
- bin/ncurses/rpn - Uses ncurses to support multiple view windows and even text-based plotting
- bin/tinygo/serialonly - The most basic version that runs on vaious microtrollers and uses serial console to talk to the host pc
- bin/tinygo/ili9341 - A tinygo mcrocontroller build that still uses serial fo input, but uses an ili9341 color LCD for output.  Supports fancy pixel-based graping.  The pins are configured for a raspberry pi pico or pico2.  By adding a pin mappiong file, it should be possible to support other microcontrollers too.

To build the normal build versions, you cd into the directory you want and type

```
go build
```

For tinygo builds, you'll need to install tinygo. After that, the command is along the lines of.

```
tinygo build -target=pico
```

Adjusting to the miccrocontroller you are using and perhaps using `flash` instead of `build`.

Note that th ncurses build will likely give you an error unless you have `libncurses-dev` already installed.  In Ubntu/Debian, the fix is:

```
sudo apt install libncurses-dev
```

In the top-level directory, you can type

```
make
```

To run all unit tests and build all targets/  This will only work if you have tinygo installed.

## Features

A quick list of things the calculator can do:

- All regular and scientific calculator operations (e.g. `+`, `-`, `sqrt`, 'sin`, ...)
- Working with the following number formats: complex, integer, binary, octal, hexidecimal
- Bitwise and logical operations
- Working with string data
- Unit conversion (e.g. miles/hour -> meters/sec)
- 2D plotting (regular and parametric)
- Simple programming
- Variables
- Customisable window layouts
- Custmizable keyboard shortcuts
- Disk and serial IO
- Build in editor

## RPN Introduction

This "users guide" will take the format of being mostly working examples.
Let's start with the basics.

RPN is an old and proven way to do calculations that is popular in engineering fields.
Much has been written on the subject. You can start here: https://en.wikipedia.org/wiki/Reverse_Polish_notation

Let's start with '2 + 3'


    2
    3
    +

or

    2 3 +

Both formns (spaces and new lines) do the following:

1. Push a `2` to the stack
2. Push a `3` to the stack
3. Call the `+` operator which pulls 2 stack values (`2`, `3`) and pushes the result `5`

The nice thing about RPN is that you don't have to enter paranthesis, which
many people find faster and less error prone.  For example, do calculate `sqrt((10 - 2)/(5 - 3))`,
you could say:

    10 2 - 5 3 - / sqrt

## Variables

You can define and use variables.  This can make it easier to work through
calculations. Most implementation will define some varialbes, such as `$pi`,
on startup.

    5 a=
    2 $a +  -> 7

Variables that start with a `.` are hidden by default in the variable window
(which we'll get to later).

There are also special variables, `$0`, `$1`, etc, which represent values on
stack. `$0` represents the value at the top of the stack. If the stack
is empty, then using `$0` results in an error.  e.g.

    5 sq    # square the number
    5 $0 *  # same result

## Stack Shifts

You can move values around the stack.

    10 20 30   # put 10 20 30 on th stack
    1>         # now the stack is 10 30 20
    1<         # back to 10 20 30
    2>         # now the stack is 20 30 10
    2<         # back to 10 20 30
    $2 2<      # now it's 10 10 20 30

## Strings

There are three ways to specify a string

    "hello world"
    'hello world'
    {hello world}

All three are just strings.  The third form can be useful when you want to nest
terms in a program.

    {0 x= {$x 1 + x= $x println 1000 <} for} count=

It could also be done with other string types:

    '0 x= "$x 1 + x= $x println 1000 <" for' count=

but this is arguably not as easy to read and further nesting would
make it even less readable than using `{}`.

## Macros

You can define a string as a macro, here is one for the area of a circle (`$pi * r * r / 2`),
given the radius:

    {$0 * $pi * 2 /} carea=
    5 @carea -> 39.269908169872416

## Conditionals and simple programs

These can be combined with macros and varialbes for simple programming.  Here
is a program that counts to 100.

    "print d 100 == '`loop' ifjmp 1 +" loop= 0 `loop

Breakdown:

- `print` Print the head of the stack
- `d` Duplicate the head of the stack
- `100` push 100 to the stack
- `==` pop 2 elements from the stack and push 1 if they are equal, 0 otherwise
- `'\`loop'` A string that contains a macro, pushed to the stack
- `ifjmp` pop two elements.  If the conditional one is 1, execute the string one instead of the rest of the string
- `1 +` Adds one to the loop
- `loop=` define the macro
- `0` initial condition
- `\`loop` execute

## Type Conversions

A `;` character can be used to add a type to a number.  Adding types can help prevent errors
during calculations.  For example, the area of a circle:

    5;mm d * .pi *
    -> 78.53981633974483 mm*mm

    mm*mm->m*m
    -> 7.853981633974483e-05 m*m

    5;m 10;s /
    -> 0.5 m/s

    m/s->miles/h
    -> 1.118468146027201 miles/h

    5;m 5;mm +
    -> 5.005 m

    5;m/s 10m +
    -> error mismatched units

## Integers, Hex, Binary

The C convention is used for inputting the values.  How the values are displayed depends on
the code that uses the library:

    0xff   # hex
    0b1101 # binary
    0i123  # integer

## Windows

The RPN calculator has differerent optoional display windows, including:

- Input
- Stack Display (in various formats)
- Graphics

Every display window (other than Input) can have multiple instances.  This allows
for multiple graphs, seeing hex along decimal values, etc.

The layout of these windows follows a window-group -> window tree with the window
group `root` always at the top.

Window groups can contain window or other wiundow groups.  They have either a
horizontal or vertical tile orientation. Each child in a window group gets
a default "weight" of 100. Relative weights determine how much space
a window group child is given relative to it's siblings.

For example, say we want the following:

```
+--------------------------------------------------------------+---------------+
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                   GRAPH (g1)                                 |  STACK (s1)   |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
|                                                              |               |
+--------------------------------------------------------------+               |
|                                                              |               |
|                                                              |               |
|                   INPUT (i)                                  |               |
|                                                              |               |
|                                                              |               |
+--------------------------------------------------------------+---------------+
```

We would get this using the following commands:

```
resetwg              # reset any existing wnidow groups, leading to 'root' with 'i'
'root' vertwg        # set root window group orientation to vertical 
'gr1' 'root' newwg   # create a left window group and add it to root
'g1' 'gr1' newgraph  # create a graph window and add it to 'gr1'
'i' 'gr1' movewg     # move window 'i' from 'root' to 'gr1'
'i' 30 weight        # change the weight of the 'i' window to 30 so it takes less space
's1' 'root' newstack # create a new stack window
's1' 25 weight       # change the weight of the 's1' window so it takes less space
```

Normally, you would not use these commands on the fly but instead have precooked
macros that can change the window layout to your preferences.

