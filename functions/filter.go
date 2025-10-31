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
	"\n" +
	"See Also: filterm, filtern, filtermn"

func Filter(r *rpn.RPN) error {
	return filterMN(r, len(r.Frames)-2, 0)
}

const FilterMHelp = "This function works like filter, but selects the top m elements " +
	"on the stack instead of working on every element (m does not include filterm arguments)" +
	"\n" +
	"Examples:\n" +
	"{2 *} 5 filterm  # multiply the top 5 stack values by 2\n" +
	"{$0 100 >= {0/} if} 10 filterm  # filter the top 10 stack values\n" +
	"\n" +
	"See Also: filter, filtern, filtermn"

func FilterM(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	m, err := f.Int()
	if err != nil {
		return err
	}
	return filterMN(r, int(m-1), 0)
}

const FilterNHelp = "This function works like filter, but ends at element head-n (not including " +
	"filtern arguments)\n" +
	"This is useful for preserving \"working values\"" +
	"\n" +
	"Examples:\n" +
	"0 {+} 1 filtern # sum all values\n" +
	"$0 {$1 $1 < {0/} {1/} ifelse} 1 filtern  # minimum\n" +
	"\n" +
	"See Also: filter, filterm, filtermn"

func FilterN(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	n, err := f.Int()
	if err != nil {
		return err
	}
	return filterMN(r, len(r.Frames)-2, int(n))
}

const FilterMNHelp = "This function works like filterm and filtern combined, filtering\n" +
	"from m elements back to head-n (not including filtermn arguments)\n" +
	"\n" +
	"Examples:\n" +
	"0 {+} 5 1 filtermn # sum top 5 values\n" +
	"$0 {$1 $1 < {0/} {1/} ifelse} 5 1 filtermn  # minimum of top 5 values\n" +
	"\n" +
	"See Also: filter, filterm, filtern"

func FilterMN(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	n, err := f.Int()
	if err != nil {
		return err
	}
	f, err = r.PopFrame()
	if err != nil {
		return err
	}
	m, err := f.Int()
	if err != nil {
		return err
	}
	return filterMN(r, int(m), int(n))
}

// Common execution function
func filterMN(r *rpn.RPN, m, n int) error {
	fn, err := r.PopFrame()
	if err != nil {
		return err
	}
	endIdx := len(r.Frames) - n
	if endIdx < 0 {
		return rpn.ErrNotEnoughStackFrames
	}
	if m == 0 {
		return nil
	}
	startIdx := len(r.Frames) - m - 1
	if startIdx < 0 {
		return rpn.ErrNotEnoughStackFrames
	}
	if endIdx < startIdx {
		return rpn.ErrIllegalValue
	}
	if endIdx == startIdx {
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
	for i := startIdx; i < endIdx; i++ {
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
	newSize := len(r.Frames) - (endIdx - startIdx)
	if newSize <= 0 {
		r.Frames = r.Frames[:0]
	} else {
		copy(r.Frames[startIdx:], r.Frames[endIdx:])
		r.Frames = r.Frames[:newSize]
	}
	return nil
}
