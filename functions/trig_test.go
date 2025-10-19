package functions

import (
	"mattwach/rpngo/rpn"
	"testing"
)

func TestSin(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"sin"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "sin"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "sin", "3", "round"},
			Want: []string{"0.841"},
		},
		{
			Args: []string{"1+i", "sin", "3", "round"},
			Want: []string{"1.298+0.635i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestCos(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"cos"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "cos"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "cos", "3", "round"},
			Want: []string{"0.54"},
		},
		{
			Args: []string{"21+i", "cos", "3", "round"},
			Want: []string{"-0.845-0.983i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestTan(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"tan"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "tan"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "tan", "3", "round"},
			Want: []string{"1.557"},
		},
		{
			Args: []string{"1.5+2.1i", "tan", "3", "round"},
			Want: []string{"0.004+1.03i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
