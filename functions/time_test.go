package functions

import (
	"mattwach/rpngo/rpn"
	"testing"
)

func TestDelay(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"delay"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"0.1", "delay"},
		},
		{
			Args: []string{"0d", "delay"},
		},
		{
			Args: []string{"-0.1", "delay"},
		},
		{
			Args:    []string{"true", "delay"},
			Want:    []string{"true"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
