package shell

import (
	"io"
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"os/exec"
	"strings"
)

func Register(r *rpn.RPN) {
	r.Register("sh", Shell, rpn.CatIO, ShellHelp)
	r.Register("source", Source, rpn.CatIO, SourceHelp)
	r.Register(".", Source, rpn.CatIO, SourceHelp)
}

const ShellHelp = `Executes a string as a shell command.

There are many ways to execute a shell command and the following special
variables control the behavior:

.stdin  - If set, the contents will be sent to stdin of the process
.stdout - If empty or false, stdout/stderr is simply printed.  
          If set to true, stdout/stderr is pushed to the stack
.env    - If set, environment variabls will be set using KEY=VALUE with
          one variable per line

The exit code of the shell command is set to the variable $.rc.
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

	stdin, err := checkStdinVar(r)
	if err != nil {
		return err
	}
	if stdin != nil {
		cmd.Stdin = stdin
	}

	output, err := cmd.CombinedOutput()
	rc := 0
	if err != nil {
		// TODO handle RC better.
		rc = 1
		r.Print("Error: " + err.Error() + " " + string(output) + "\n")
	}
	r.PushInt(int64(rc), rpn.INTEGER_FRAME)
	r.SetVariable(".rc")
	if err == nil {
		if err := setCmdOutput(r, string(output)); err != nil {
			return err
		}
	}
	return nil
}

func checkStdinVar(r *rpn.RPN) (io.Reader, error) {
	val, err := r.GetStringVariable(".stdin")
	if err != nil {
		return nil, nil
	}
	return strings.NewReader(val + "\n"), nil
}

func setCmdOutput(r *rpn.RPN, output string) error {
	stack := false
	stdout, err := r.GetVariable(".stdout")
	if err == nil {
		if stdout.Type != rpn.BOOL_FRAME {
			return rpn.ErrExpectedABoolean
		}
		stack = stdout.Bool()
	}
	if stack {
		return r.PushString(strings.TrimSpace(output))
	}
	r.Print(output)
	return nil
}
