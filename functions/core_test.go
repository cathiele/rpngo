package functions

import (
	"mattwach/rpngo/rpn"
	"testing"
)

func TestAdd(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"+"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "+"},
			WantErr: rpn.ErrNotEnoughStackFrames,
			Want:    []string{"1"},
		},
		{
			Args:    []string{"1", "true", "+"},
			WantErr: rpn.ErrIllegalValue,
			Want:    []string{"1", "true"},
		},
		{
			Args:    []string{"true", "1", "+"},
			WantErr: rpn.ErrIllegalValue,
			Want:    []string{"true", "1"},
		},
		{
			Args: []string{"1", "2", "+"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1", "2d", "+"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1d", "2", "+"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1d", "2d", "+"},
			Want: []string{"3d"},
		},
		{
			Args: []string{"'foo'", "7", "+"},
			Want: []string{"\"foo7\""},
		},
		{
			Args: []string{"7", "'foo'", "+"},
			Want: []string{"\"7foo\""},
		},
		{
			Args: []string{"'foo'", "7d", "+"},
			Want: []string{"\"foo7d\""},
		},
		{
			Args: []string{"7d", "'foo'", "+"},
			Want: []string{"\"7dfoo\""},
		},
		{
			Args: []string{"'foo'", "true", "+"},
			Want: []string{"\"footrue\""},
		},
		{
			Args: []string{"true", "'foo'", "+"},
			Want: []string{"\"truefoo\""},
		},
		{
			Args: []string{"\"foo\"", "'bar'", "+"},
			Want: []string{"\"foobar\""},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
