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

func TestSubtract(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"-"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "-"},
			WantErr: rpn.ErrNotEnoughStackFrames,
			Want:    []string{"1"},
		},
		{
			Args:    []string{"1", "true", "-"},
			WantErr: rpn.ErrExpectedANumber,
			Want:    []string{"1", "true"},
		},
		{
			Args:    []string{"true", "1", "-"},
			WantErr: rpn.ErrExpectedANumber,
			Want:    []string{"true", "1"},
		},
		{
			Args: []string{"1", "2", "-"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"1", "2d", "-"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"1d", "2", "-"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"1d", "2d", "-"},
			Want: []string{"-1d"},
		},
		{
			Args:    []string{"'foo'", "7", "-"},
			Want:    []string{"\"foo\"", "7"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "-"},
			Want:    []string{"7", "\"foo\""},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "-"},
			Want:    []string{"\"foo\"", "\"bar\""},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestMultiply(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"*"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "*"},
			WantErr: rpn.ErrNotEnoughStackFrames,
			Want:    []string{"1"},
		},
		{
			Args:    []string{"1", "true", "*"},
			WantErr: rpn.ErrExpectedANumber,
			Want:    []string{"1", "true"},
		},
		{
			Args:    []string{"true", "1", "*"},
			WantErr: rpn.ErrExpectedANumber,
			Want:    []string{"true", "1"},
		},
		{
			Args: []string{"2", "3", "*"},
			Want: []string{"6"},
		},
		{
			Args: []string{"2", "3d", "*"},
			Want: []string{"6"},
		},
		{
			Args: []string{"2d", "3", "*"},
			Want: []string{"6"},
		},
		{
			Args: []string{"2d", "3d", "*"},
			Want: []string{"6d"},
		},
		{
			Args:    []string{"'foo'", "7", "*"},
			Want:    []string{"\"foo\"", "7"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "*"},
			Want:    []string{"7", "\"foo\""},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "*"},
			Want:    []string{"\"foo\"", "\"bar\""},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestDivide(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"/"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "/"},
			WantErr: rpn.ErrNotEnoughStackFrames,
			Want:    []string{"1"},
		},
		{
			Args:    []string{"1", "true", "/"},
			WantErr: rpn.ErrExpectedANumber,
			Want:    []string{"1", "true"},
		},
		{
			Args:    []string{"true", "1", "/"},
			WantErr: rpn.ErrExpectedANumber,
			Want:    []string{"true", "1"},
		},
		{
			Args: []string{"5", "2", "/"},
			Want: []string{"2.5"},
		},
		{
			Args: []string{"5", "2d", "/"},
			Want: []string{"2.5"},
		},
		{
			Args: []string{"5d", "2", "/"},
			Want: []string{"2.5"},
		},
		{
			Args: []string{"5d", "2d", "/"},
			Want: []string{"2d"},
		},
		{
			Args:    []string{"'foo'", "7", "/"},
			Want:    []string{"\"foo\"", "7"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "/"},
			Want:    []string{"7", "\"foo\""},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "/"},
			Want:    []string{"\"foo\"", "\"bar\""},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestNegate(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"neg"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"'foo'", "neg"},
			Want:    []string{"\"foo\""},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"1", "neg"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"-1", "neg"},
			Want: []string{"1"},
		},
		{
			Args: []string{"0", "neg"},
			Want: []string{"-0"},
		},
		{
			Args: []string{"1d", "neg"},
			Want: []string{"-1d"},
		},
		{
			Args: []string{"0d", "neg"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"true", "neg"},
			Want: []string{"false"},
		},
		{
			Args: []string{"false", "neg"},
			Want: []string{"true"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestExec(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"@"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"''", "@"},
		},
		{
			Args: []string{"'3'", "@"},
			Want: []string{"3"},
		},
		{
			Args: []string{"'2 3 +'", "@"},
			Want: []string{"5"},
		},
		{
			Args:    []string{"foo", "@"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestRand(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"rand", "$0", "0", ">=", "1>", "1", "<="},
			Want: []string{"true", "true"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestReal(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"real"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1d", "real"},
			Want:    []string{"1d"},
			WantErr: rpn.ErrExpectedAComplexNumber,
		},
		{
			Args:    []string{"'foo'", "real"},
			Want:    []string{"\"foo\""},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "real"},
			Want: []string{"1"},
		},
		{
			Args: []string{"i", "real"},
			Want: []string{"0"},
		},
		{
			Args: []string{"2+3i", "real"},
			Want: []string{"2"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestImag(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"imag"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1d", "imag"},
			Want:    []string{"1d"},
			WantErr: rpn.ErrExpectedAComplexNumber,
		},
		{
			Args:    []string{"'foo'", "imag"},
			Want:    []string{"\"foo\""},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "imag"},
			Want: []string{"0"},
		},
		{
			Args: []string{"i", "imag"},
			Want: []string{"i"},
		},
		{
			Args: []string{"2+3i", "imag"},
			Want: []string{"3i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestTrueFalse(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"true"},
			Want: []string{"true"},
		},
		{
			Args: []string{"false"},
			Want: []string{"false"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
