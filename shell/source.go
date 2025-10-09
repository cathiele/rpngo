package shell

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"os"
)

const SourceHelp = "Loads the given path and executes commands within it.\n" +
	"Example: 'myfile.txt' source"

func Source(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fields := make([]string, 256)
	fields, err = parse.Fields(string(data), fields)
	if err != nil {
		return err
	}
	if err := r.Exec(fields); err != nil {
		return err
	}
	return nil
}
