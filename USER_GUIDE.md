## Build

You'll first need to install golang.

There is a `bin/` directory thast contains various different builds of rpn
for the computer and for microtrollers.  Try the ones you like and ignore the rest/  Here is an overview:

- bin/minimal/rpn - The most basic version.  Parses args and exsts.
- bin/ncurses/rpn - Uses ncurses to support multiple view windows and even text-based plotting
- bin/tinygo/serialonly - The most basic version that runs on vaious microtrollers and uses serial console to talk to the host pc
- bin/tinygo/ili9341 - A tinygo mcrocontroller build that still uses serial fo input, but uses an ili9341 color LCD for output.  Supports fancy pixel-based graping.  The pins are configured for a raspberry pi pico or pico2.  By adding a pin mappiong file, it should be possible to support other microcontrollers too.

To build the normal build versions, you cd into the directory you want and type

### Desktop / Raspberry Pi

```
go build
```

Note that th ncurses build will likely give you an error unless you have `libncurses-dev` already installed.  In Ubntu/Debian, the fix is:

```
sudo apt install libncurses-dev
```

## Microcontrollers using TinyGo (Raspberry Pi PICO and PICO2 tested as working)

You'll need to install tinygo. After that, check the `Makefile` for the correct command

e.g. to look

```
$ cd bin/tinygo/picocalc
$ make -n flash
tinygo flash -target=pico2 -scheduler=tasks -serial=uart
```

to actually flash the chip

```
$ cd bin/tinygo/picocalc
$ make flash
```

or

```
$ cd bin/tinygo/picocalc
$ tinygo flash -target=pico2 -scheduler=tasks -serial=uart
```

### Build Everything and Run Unit Tests

In the top-level directory, you can type

```
make
```

or 

```
make all
```

to build the more obsure targets too.  Tinygo is needed to build everything


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

One difference, however, is the "up arrow" for command history.  Sometimes it's nice to batch
a set of commands in a line with spaces to make the command history more useful.

The nice thing about RPN is that you don't have to enter paranthesis, which
many people find faster and less error prone.  For example, do calculate `sqrt((10 - 2)/(5 - 3))`,
you could say:

    10 2 - 5 3 - / sqrt

## Editing features

- Left, right, insert, backspace, home, and end all work like you would expect
- Press up and down arrows to scroll through command history
- Ctrl-C to cancel a running program
- Ctrl-D to exit the program
- Press "esc" or "page up" to enter scrolling mode where you can
  use page Up, page Down, up arrow and down arrow to view text
  that has scrolled off the top of the window. "esc" or scrolling
  down far enough exits this mode.

Type `edit` to enter a full-panel multiline editor.  The window
will contain the value at the top of the stack.  For example:

```
'animate_sin.rpn' load edit
```

The editor is intended for basic tasks and supports only the following:

- arrow key, page up, page down navigation
- insert, replace, backspace, delete
- syntax highlighting

While editing, press "esc" to keep your changes (which will be
at the top of the stack) or "ctrl-c" to exit without changing the
value.  If you happen to want to save your work permanently, exit with
"esc", then type something like:

```
'my_file.txt' save
```

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

## Stack Deletions

    10 20 30   # put 10 20 30 on the stack
    1/         # now the stack is 10 30
    0/         # now the stack is 10
    X          # emptys all values from the stack

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

    {sq $pi * 2 /} carea=
    5 @carea -> 39.26990817

## Conditonals

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
