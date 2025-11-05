package functions

import (
	"mattwach/rpngo/rpn"
	"sort"
)


const SortHelp = "Sorts values on the stack.  Uses the " +
	"< conditional as the baseline for sorting."

func Sort(r *rpn.RPN) error {
	return sortMN(r, len(r.Frames), 0)
}

const SortNHelp = "Sorts values on the stack from the bottom to n frames back."

func SortN(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	n, err := f.Int()
	if err != nil {
		return err
	}
	return sortMN(r, len(r.Frames), int(n))
}

type sortInterface struct {
	data []rpn.Frame
}

func (si *sortInterface) Len() int {
	return len(si.data)
}

func (si *sortInterface) Less(i, j int) bool {
	return si.data[i].IsLessThan(si.data[j])
}

func (si *sortInterface) Swap(i, j int) {
	si.data[i], si.data[j] = si.data[j], si.data[i]
}

func sortMN(r *rpn.RPN, m, n int) error {
	endIdx := len(r.Frames) - n
	if endIdx < 0 {
		return rpn.ErrNotEnoughStackFrames
	}
	if m == 0 {
		return nil
	}
	startIdx := len(r.Frames) - m
	if startIdx < 0 {
		return rpn.ErrNotEnoughStackFrames
	}
	if endIdx < startIdx {
		return rpn.ErrIllegalValue
	}
	if endIdx == startIdx {
		return nil
	}
	si := sortInterface{data: r.Frames[startIdx:endIdx]}
	sort.Sort(&si)
	return nil
}

