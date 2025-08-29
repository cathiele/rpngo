package rpn

import (
	"fmt"
	"sort"
)

func (rpn *RPN) initHelp() {
	rpn.conceptHelp = map[string]string{
		"basics": "- Enter numbers to push them to the stack\n" +
			"- Numbers can be separated by spaces or newlines\n" +
			"- Enter an operator to replace numbers on the stack with a result\n" +
			"- For example: 2 3 +",

		"complex": "Enter a complex value as i, -i, 3+i or 3-i\n" +
			"Do not use spaces.",

		"macros": "Execute a variable as @name\n" +
			"Example:\n" +
			"'. 3.14159 * *' cirarea=\n" +
			"5 @cirarea\n" +
			"See Also: variables",

		"strings": "Enter a string value as 'example 1' or \"example 2\"",

		"variables": "Set a variable as name=\n" +
			"Use a variable with $name\n" +
			"Example: 5 x= $x $x *\n" +
			"\n" +
			"Push a variable frame with vpush, pop with vpop\n" +
			"Variables added after vpush will be reverted after a vpop,\n" +
			"allowing for 'local variables' to be temporarily defined." +
			"\n" +
			"See Also: macros",
	}
	rpn.commandHelp = make(map[string]string)
}

func (r *RPN) RegisterConceptHelp(concept, help string) {
	r.conceptHelp[concept] = help
}

func (r *RPN) PushHelp(topic string, windoww int) error {
	if len(topic) == 0 {
		r.listCommands(windoww)
		return nil
	}
	help, ok := r.conceptHelp[topic]
	if !ok {
		help, ok = r.commandHelp[topic]
	}
	if !ok {
		return fmt.Errorf("no help found for %s. Use ? to list all", topic)
	}
	r.PushMessage("")
	r.PushMessage(help)
	r.PushMessage("")
	return nil
}

func (r *RPN) listCommands(windoww int) {
	r.dumpMap("Concepts", windoww, r.conceptHelp)
	r.dumpMap("Commands", windoww, r.commandHelp)
}

const colWidth = 40

func (r *RPN) dumpMap(title string, windoww int, m map[string]string) {
	r.PushMessage(title)
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
			r.PushMessage(string(line))
			line = line[:2]
			nextCol = len(line) + colWidth
		}
	}
	if len(line) > 0 {
		r.PushMessage(string(line))
	}
}
