package functions

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
	"os"
)

const LoadHelp = "Loads the given path and executes commands within it.\n" +
	"Example: 'myfile.txt' load"

func Load(r *rpn.RPN) error {
	path, err := r.PopString()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fields, err := parse.Fields(string(data))
	if err != nil {
		return err
	}
	if err := r.Exec(fields); err != nil {
		return err
	}
	return nil
}
