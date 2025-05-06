#RPNGO

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

