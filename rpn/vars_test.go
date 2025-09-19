package rpn

import "testing"

func TestPushPopVarFrame(t *testing.T) {
	data := []UnitTestExecData{
		{
			Name:    "empty pop",
			Args:    []string{"vpop"},
			WantErr: ErrStackEmpty,
		},
		{
			Name: "push",
			Args: []string{"1234", "x=", "vpush", "$x"},
			Want: []string{"1234"},
		},
		{
			Name: "push, then pop",
			Args: []string{"1234", "x=", "vpush", "2345", "$x", "x=", "vpop", "$x"},
			Want: []string{"2345", "1234"},
		},
		{
			Name: "push twice, then pop twice",
			Args: []string{"1234", "x=", "vpush", "$x", "vpush", "$x", "2345", "x=", "$x", "vpop", "vpop", "$x"},
			Want: []string{"1234", "1234", "2345", "1234"},
		},
	}
	UnitTestExecAll(t, data, nil)
}
