package functions

import (
	"mattwach/rpngo/rpn"
	"testing"
)

func TestFilter(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"filter"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"{sq}", "filter"},
		},
		{
			Args: []string{"1", "2", "3", "{sq}", "filter"},
			Want: []string{"1", "4", "9"},
		},
		{
			Args: []string{"10", "100", "30", "200", "{$0 100 >= {0/} if}", "filter"},
			Want: []string{"10", "30"},
		},
		{
			Args: []string{"1", "2", "3", "0", "sum=", "{$sum + sum=}", "filter", "$sum"},
			Want: []string{"6"},
		},
		{
			Args: []string{"10", "5", "7", "12", "$0", "min=", "{$0 $min < {min=} {0/} ifelse}", "filter", "$min"},
			Want: []string{"5"},
		},
		{
			Args: []string{"1", "2", "3", "{0/}", "filter"},
		},
		{
			Args:    []string{"1", "2", "3", "{0/ 0/}", "filter"},
			Want:    []string{"1"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"1", "2", "3", "4", "{sq}", "2", "filterm"},
			Want: []string{"1", "2", "9", "16"},
		},
		{
			Args: []string{"10", "100", "30", "200", "{$0 100 >= {0/} if}", "2", "filterm"},
			Want: []string{"10", "100", "30"},
		},
		{
			Args:    []string{"filterm"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "filterm"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "filterm"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"{sq}", "0", "filterm"},
		},
		{
			Args: []string{"{sq}", "-1", "filterm"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"1", "2", "3", "0", "{+}", "1", "filtern"},
			Want: []string{"6"},
		},
		{
			Args: []string{"10", "2", "30", "5", "$0", "{$1 $1 < {0/} {1/} ifelse}", "1", "filtern"},
			Want: []string{"2"},
		},
		{
			Args:    []string{"filtern"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "filtern"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "filtern"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"{sq}", "0", "filtern"},
		},
		{
			Args: []string{"{sq}", "-1", "filtern"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"{sq}", "1", "filtern"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"4", "5", "6", "7", "{sq}", "2", "1", "filtermn"},
			Want: []string{"4", "7", "25", "36"},
		},
		{
			Args: []string{"10", "2", "30", "5", "0", "{+}", "3", "1", "filtermn"},
			Want: []string{"10", "37"},
		},
		{
			Args: []string{"10", "2", "30", "5", "$0", "{$1 $1 < {0/} {1/} ifelse}", "3", "1", "filtermn"},
			Want: []string{"10", "2"},
		},
		{
			Args:    []string{"filtermn"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "filtermn"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "1", "filtermn"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"{sq}", "0", "0", "filtermn"},
		},
		{
			Args:    []string{"-1", "0", "filtermn"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "-1", "filtermn"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"{sq}", "1", "0", "filtermn"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"{sq}", "0", "1", "filtermn"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "{sq}", "1", "1", "filtermn"},
			Want:    []string{"1"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "2", "3", "{sq}", "1", "3", "filtermn"},
			Want:    []string{"1", "2", "3"},
			WantErr: rpn.ErrIllegalValue,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
