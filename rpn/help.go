package rpn

import (
	"fmt"
	"sort"
)

func (rpn *RPN) initHelp() {
	conceptHelp := map[string]string{
		"base": "Different number bases are supported:\n" +
			"Examples:\n" +
			"  1.1 -1 0 # complex floats\n" +
			"  2+3i # also a complex float\n" +
			"  32d # a decimal integer\n" +
			"  ffx # hexidecimal ff (255 in decimal)\n" +
			"  55o # octal 55 (45 in decimal)\n" +
			"  1001b # binary (9 in decimal)\n" +
			"Numbers can be converted to different formats:\n" +
			"  bin # convert to binary\n" +
			"  float # convert to complex float\n" +
			"  hex # convert to hecidecimal\n" +
			"  int # cnvert to integer\n" +
			"  oct # cnvert to octal\n" +
			"  str # cnvert to string",

		"basics": "- Enter numbers to push them to the stack\n" +
			"- Numbers can be separated by spaces or newlines\n" +
			"- Enter an operator to replace numbers on the stack with a result\n" +
			"- For example: 2 3 +",

		"conditionals": "The operators, >, >=, <, <=, =, != can be used to\n" +
			"compare numbers.  Note that > and friends ignore the complex\n" +
			"part of numbers.\n" +
			"Example: 3 5 > # this will put false on the stack",

		"control": "The if, ifelse, and for operators can be used to provide\n" +
			"support for simple programmng\n" +
			"Example (print 1-50): 1 'println 1 + c 50 <' for",

		"complex": "Enter a complex value as i, -i, 3+i or 3-i\n" +
			"Do not use spaces.",

		"macros": "Execute a variable as @name.  Execute a string with just @\n" +
			"Convert any value to a string with str.\n" +
			"Example:\n" +
			"'. 3.14159 * *' cirarea=\n" +
			"5 @cirarea\n" +
			"See Also: variables",

		"plot": "Plot functions using plot. Plot will push an 'x' value to the stack,\n" +
			"run the provided string, and pop the value as y value.\n" +
			"Examples:\n" +
			"    '2 *' plot # plots y = x * 2\n" +
			"    'c *' plot # plots y = x * x\n" +
			"    'sin' plot # plots y = sin(x)\n" +
			"Various properties can be set on the plot window to change the number\n" +
			"of points and the boundaries of the plot.\n" +
			"There are some special variables that plot uses:\n" +
			"    plot.win  : Name of the window to send plots to (there can be more than one)\n" +
			"                at a time.\n" +
			"    plot.init : If no plot window exists, this string is executed and is expected\n" +
			"                create one. Making this a variable allows for user customization.\n" +
			"See Also: window.props, plot.parametric",

		"plot.parametric": "Plot parametric functions using pplot. pplot will push a 't' value to\n" +
			"the stack, run the provided string then pop y, then x to determine the plot point x, y\n" +
			"Examples:\n" +
			"    'c cos sw sin' pplot # draws an arc or full circle, depending on t range\n" +
			"    't= $t sin $t * $t cos $t *' pplot # draw a spiral\n" +
			"    '1 sw' draw a vertical line\n",

		"printing": "There are various printing functions that print values\n" +
			"at the head of the stack. These include:\n" +
			"  - print : print the value at the head of the stack\n" +
			"  - printx : pop and print the value at the head of the stack\n" +
			"  - prints : print with a space\n" +
			"  - printsx : printx with a space\n" +
			"  - println : print with a newline\n" +
			"  - printlnx : printx with a newline",

		"stack": "Operators are provided to manipute the stack to set up calculations\n" +
			"    If things are getting complex, consider using variables.\n" +
			"Examples:\n" +
			"  $x c + 1 /  # Uses c to copy the stack head and execute a/(a+1)'\n" +
			"  'c cos sw sin' pplot # parametric plot of a circle using sw to swap elements\n" +
			"   x # drop the element at the head of the stack\n" +
			"   1 xi # drop the second element from the stack\n" +
			"   X # clear the stack\n" +
			"   2 mi # Moves the third element in the stack to the head\n" +
			"   2 ci # Copies the third element in the stack to the head\n" +
			"\n",

		"strings": "Enter a string value as 'example 1' or \"example 2\"",

		"variables": "Set a variable as name=\n" +
			"Use a variable with $name\n" +
			"Example: 5 x= $x $x *\n" +
			"\n" +
			"Clear a variable with a trailing /. e.g. x/\n" +
			"Push a variable frame with vpush, pop with vpop\n" +
			"Variables added after vpush will be reverted after a vpop,\n" +
			"allowing for 'local variables' to be temporarily defined." +
			"\n" +
			"See Also: macros",

		"window.layout": "Windows are arranged with window groups.  There\n" +
			"is always a window group named 'root' which is the parent of all \n" +
			"windows and groups.\n" +
			"- Add a new window group to the root window with w.new.group.\n" +
			"- Move a window or group to a different window group with w.move.beg and w.move.end\n" +
			"- Change the weight of a window or group with w.weight (default weight is 100).\n" +
			"- Change the layout mode of a window group to columns with w.columns.\n" +
			"- Print info on all existing windows and groups with w.dump.\n" +
			"See Also: windows, window.props",

		"window.props": "Each window supports properties that changes how the window operates\n" +
			"- Print all properties and values for a window with w.listp\n" +
			"- Get a single property with w.getp\n" +
			"- Set a single property with w.setp\n" +
			"See Also: windows, window.layout, plotting",

		"windows": "The display can be customized with different windows\n" +
			"- Add a window with a w.new.<type> command. Example: w.new.stack\n" +
			"- Reset to a single window with w.reset.\n" +
			"See Also: window.layout, window.props",
	}
	rpn.help = map[string]map[string]string{CatConcepts: conceptHelp}
}

func (r *RPN) RegisterConceptHelp(concept string, help string) {
	r.help[CatConcepts][concept] = help
}

func (r *RPN) PrintHelp(topic string, windoww int) error {
	if len(topic) == 0 {
		r.listCommands(windoww)
		return nil
	}
	var help string
	for cat := range r.help {
		help = r.help[cat][topic]
		if help != "" {
			break
		}
	}
	if help == "" {
		return fmt.Errorf("no help found for %s. Use ? to list all", topic)
	}
	r.Print("\n")
	r.Println(help)
	return nil
}

func (r *RPN) listCommands(windoww int) {
	var cats []string
	for cat := range r.help {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	for _, cat := range cats {
		r.dumpMap(cat, windoww, r.help[cat])
	}
}

const colWidth = 32

func (r *RPN) dumpMap(title string, windoww int, m map[string]string) {
	r.Println(title)
	var topics []string
	for k := range m {
		topics = append(topics, k)
	}
	sort.Strings(topics)
	line := []byte("  ")
	nextCol := len(line) + colWidth
	for _, t := range topics {
		line = append(line, []byte(t)...)
		for len(line) < nextCol {
			line = append(line, ' ')
		}
		nextCol += colWidth
		if nextCol > windoww {
			r.Println(string(line))
			line = line[:2]
			nextCol = len(line) + colWidth
		}
	}
	if len(line) > 0 {
		r.Println(string(line))
	}
}
