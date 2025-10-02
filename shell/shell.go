package shell

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"os/exec"
)

func Register(r *rpn.RPN) {
	r.Register("sh", Shell, rpn.CatIO, ShellHelp)
}

const ShellHelp = `Executes a string as a shell command.

There are many ways to execute a shell command and the following special
variables control the behavior:

.stdin -  This can be set to '$v' to use the value of $v as stdin
          If can be set to 'stack' to pull stdin from the stack
		  If can be set to any other value to try and load/use a file
		  If empty, no stdin will be provided
.stdout - This can be set to v= to store stdout to $v
          It can be set to stack to push the stdout on the stack
		  It can be set to any other string to try and save stdout to a file
		  If empty, stdout will be printed and not captured
.env    - If set, environment variabls will be set using KEY=VALUE with
          one variable per line

The exit code of the shell command is pushed to the stack.
`

func Shell(r *rpn.RPN) error {
	s, err := r.PopString()
	if err != nil {
		return err
	}
	fields, err := parse.Fields(s)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return rpn.ErrIllegalValue
	}

	cmd := exec.Command(fields[0], fields[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		r.Print("Error: " + string(output))
		return nil
	}
	r.Print(string(output))
	return nil
}
