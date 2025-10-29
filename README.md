#RPNGO


This project is implements a programmable RPN graphing calculator using golang
and TinyGO.

It can be run on regular PCs (including Raspberry Pi):

![regular pc](img/running_on_pc.png)
![editor on pc](img/editor_on_pc.png)

Or embedded microcontrollers.  Here it is running on a breadboard:

![breadboard](img/running_on_breadboard.jpg)

And a PicoCalc:

![picocalc](img/running_on_picocalc.jpg)

## Why

I've been writing code professionally and for fun for many decades in 20+
languages. Of those languages, my favorites are golang for PC utilities
and tools and C for microcontroller projects.  What about golang on microcontrollers?

I decided to give TinyGO a try and see. Later in this doc, I talk about my views
on TinyGO so far (TLDR: I think it's fun for experimental work but is not yet
suitable for more serious efforts).

Why an RPN calculator? We'll I spent my early days on an HP-48G user and
wanted a modern calculator that worked in a similar-but-modernized way. I then
discovered the PicoCalc, an $80 platform for makig your own retro computer /
calculator:

![picocalc and hp48](img/picocalc_and_hp48.png)

The calculator ships with Micro Basic, Python, LISP and an NES emulator - plenty
to do right there! But I wanted to do something lower-level and work directly
with the hardware.

## Base Features (as of 10/2025):

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
- Disk and serial IO (still experimental on TinyGO)
- Build in editor

## Non-implemented (as of 10/2025)

- Matrix algebra
- Symbolic equation support / features
- Wifi Support

Want to learn more?  Here is a [User's Guide]](USER_GUIDE.md) that is filled
with examples.

# Upcoming Work

## Bugs

## Cleanup

## Features

- L: Basic text editor
- M: Add stack output formatting options
- L: Rewrite readme doc
- M: Experiement with SD card or on-chip storage
- M: Improve font
- M: Create snapshot
- M: Foreach, del