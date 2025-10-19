package fileops

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const SourceHelp = "Loads the given path and executes commands within it.\n" +
	"Example: 'myfile.txt' source"

func (fo *FileOps) Source(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	path := f.UnsafeString()
	data, err := fo.driver.ReadFile(path)
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
