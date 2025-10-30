package functions

import (
	"mattwach/rpngo/parse"
	"mattwach/rpngo/rpn"
)

const FilterHelp = "Copies each value from the bottom to the top of the stack " +
	"to the top of the stack.  Calling filter function after each one.  Discards the " +
	"original values, leaving only the new ones.\n" +
	"\n" +
	"Examples:\n" +
	"{2 *} filter  # multiply every stack value by 2\n" +
	"{$0 100 >= {0/} if} filter  # keep only values < 100\n" +
	"0 {+} filter  # sum all values\n" +
	"$0 min= {$0 $min < {min=} {0/} ifelse} filter $min  # find minimum\n"

func Filter(r *rpn.RPN) error {
	fn, err := r.PopFrame()
	if err != nil {
		return err
	}
	ssize := len(r.Frames)
	if ssize <= 0 {
		return nil
	}
	fields := make([]string, 0, 16)
	addField := func(t string) error {
		fields = append(fields, t)
		return nil
	}
	if err := parse.Fields(fn.String(false), addField); err != nil {
		return err
	}
	for i := 0; i < ssize; i++ {
		if i >= len(r.Frames) {
			return rpn.ErrNotEnoughStackFrames
		}
		if err := r.PushFrame(r.Frames[i]); err != nil {
			return err
		}
		if err := r.ExecSlice(fields); err != nil {
			return err
		}
	}
	newSize := len(r.Frames) - ssize
	if newSize <= 0 {
		r.Frames = r.Frames[:0]
	} else {
		copy(r.Frames, r.Frames[ssize:])
		r.Frames = r.Frames[:newSize]
	}
	return nil
}
