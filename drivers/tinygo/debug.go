package tinygo

import "strconv"

// a temporary check to narrow down panics
func Check(ctx string, idx int, maxidx int) {
	if (idx < 0) || (idx >= maxidx) {
		panic("out of bounds ctx " + ctx + " idx " + strconv.Itoa(idx) + " maxidx " + strconv.Itoa(maxidx))
	}
}
