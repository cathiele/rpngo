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

# My TinyGO Impressions (10/2025)

Traditional Go on Rasperry Pi, PCs and servers is a mature product used
in production enviroments.  For example, Kubernetes and Docker are both
implemented in Go.

TinyGo, by contrast, is not currently ready for "serious" work (e.g. used in a product for-sale)
for the following reasons (with the cavet that, by the time you read this, the
situation may have been improved).

## Memory Management

Memory management issues form the biggest problem areas for TinyGo.  

### Implicit Allocation

Memory is usually a precious resource on a microcontroller,
with a couple hundred kilobytes considered a generous offering.
As a quick refresher, there are usually 4 types of directly accessible memory:

- *Read only*: Data and code stored in flash memory
- *Static*: Globally-defined structures
- *Stack*: Memory that holds local variables for functions (assuming they are not too big)
- *Heap*: General purpose memory that can be used when the above choices don't work well
  for your project.  For example, say you are implementing an programmable RPN calculator
  and do not know what kind of programs the user will send your way `:)`

In the classic C/C++, all 4 of these are explicitly user controlled.  In many embedded applications
heap memory can be avoided entirely.  In cases where it is needed, the `malloc()` and `new` calls
can return a null pointer if they fail, giving the program more options than just crashing.

In Go, whether to use stack or heap memory is implicit and guided by heuristics (e.g. if object is
less than N bytes, use the stack otherwise the heap) and escape analysis. Both of these
can change or regress over time as they are not guaranteed to work a certain way.

In TinyGO, you can compile with `-show-allocs` to list the allocs that will occur but,
again, the report is not guaranteed to remain stable as new versions of Go are released.

### Heap Fragmentation

The current heap allocator in TinyGo will not move allocated memory.  This means that, over time
it's possible to have enough total free memory for an allocation but fail anyway.  As an
abstract example, say you have memory filled with `AAAA`, `BBB`, `CC` and some free memory
`...` like this:

    AAAABBBCC...

Now B is freed, giving us

    AAAA...CC...

Now we need a new memory region `DDDD`.  Since the allocator can not move `C` to make room,
the request fails.

Fixing this would involve reworking how pointers are implemented (behind the scenes).
It will also require a new framework be created for interrupts handlers (to avoid
the interrupt trying to access heap memory that can move).
It should be solvable, but will take work to get there.

### Panic behavior

Running out of memory thows a Panic.  TinyGO's current reaction to a panic depends on
what code throws it.  In the case of the memory allocator, it seems to be some other
task that throws the panic and the result is a message it written to the serial port
followed by a processor hang even if a panic handler was registered.  This is ok during
early development but not good behavior for a shipping product.

In my opinion, a failed memory allocation in TinyGO should not panic but return a `nil`
pointer instead.  Just like `C`, this means more checks are needed but they actually
are needed since panics are difficult to handle well.  For example, if I load
a file that's too large into the RPN calculator, I'd rather see an error "not enough memory"
over the calculator resetting or hard locking...

### Refinement and Polish

While I was working on TinyGO, the version I was using is using a new `cores` scheduler
that happens to be unstable and cause the chip to hang when it runs garbage colletion.
Switching to the `tasks` scheduler fixes this.  The real point here is that the
release version of TinyGO uses `cores` by default *and* has an error like this.
It's a sign that more refinement and polish is needed to graduate from the
"experimental" status.


