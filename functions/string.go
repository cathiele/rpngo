package functions

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const FieldsHelp = "Splits a string into fields and places all fields on the stack"

func Fields(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	fn := func(s string) error {
		return r.PushFrame(rpn.StringFrame(s, f.Type()))
	}
	return parse.Fields(f.String(false), fn)
}
