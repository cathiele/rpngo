#RPNGO

## TODO Bugs

- Border rendering on LCD is sometimes incomplete

## TODO Cleanup

## TODO Features

- L: Implement PixelWindow
- L: Implement large output viewer
- M: Add stack output formatting options
- M: Add free memory command
- S: Blink on panic
- L: Rewrite readme doc
- L: Add I2C keyboard support
- M: Ad F or ctrl key macro shortcuts

## Wishlist

- Integrated editor?
- SD Card support?
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

## RPN Introduction

RPN is an old and proven way to do calculations that is popular in engineering fields.
Much has been written on the subject. You can start here: https://en.wikipedia.org/wiki/Reverse_Polish_notation

As a quick example, here is one way to calculate `(1 + 5) / (4 - 2)`:

    1 5 + 4 2 - /

## Variables

You can define and use variables.  This can make it easier to work through
calculations. Most implementation will define some varialbes, such as `.pi`,
on startup.

    5 a=
    2 .a +  -> 7

## Macros

You can define a string as a macro, here is one for the area of a circle (`pi * r * r / 2`),
given the radius (the `d` command means "duplicate"):

    "d * .pi * 2 /" carea=
    5 `carea -> 39.269908169872416

## Conditionals and simnple programs

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

