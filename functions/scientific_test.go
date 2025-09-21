package functions

import (
	"mattwach/rpngo/rpn"
	"testing"
)

func TestPower(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"**"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"2", "**"},
			Want:    []string{"2"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"2", "3", "**"},
			Want: []string{"8"},
		},
		{
			Args: []string{"2d", "3", "**"},
			Want: []string{"8"},
		},
		{
			Args: []string{"2", "3d", "**"},
			Want: []string{"8"},
		},
		{
			Args: []string{"2d", "3d", "**"},
			Want: []string{"8d"},
		},
		{
			Args:    []string{"2", "true", "**"},
			Want:    []string{"2", "true"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "2", "**"},
			Want:    []string{"true", "2"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestSquareRoot(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"sqrt"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "sqrt"},
			Want: []string{"2"},
		},
		{
			Args: []string{"4d", "sqrt"},
			Want: []string{"2"},
		},
		{
			Args:    []string{"true", "sqrt"},
			Want:    []string{"true"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestAbs(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"abs"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "abs"},
			Want: []string{"4"},
		},
		{
			Args: []string{"-4", "abs"},
			Want: []string{"4"},
		},
		{
			Args: []string{"4d", "abs"},
			Want: []string{"4d"},
		},
		{
			Args: []string{"-4d", "abs"},
			Want: []string{"4d"},
		},
		{
			Args:    []string{"true", "abs"},
			Want:    []string{"true"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestSquare(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"sq"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "sq"},
			Want: []string{"16"},
		},
		{
			Args: []string{"4d", "sq"},
			Want: []string{"16d"},
		},
		{
			Args:    []string{"true", "sq"},
			Want:    []string{"true"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
