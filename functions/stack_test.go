package functions

import (
	"mattwach/rpngo/rpn"
	"testing"
)

func TestDropAll(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"d"},
		},
		{
			Args: []string{"1", "2", "3", "d"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
