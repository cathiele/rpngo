package rpn

import (
	"errors"
	"reflect"
	"testing"
)

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

func TestSetGetAndClear(t *testing.T) {
	data := []UnitTestExecData{
		{
			Name:    "no such var",
			Args:    []string{"$x"},
			WantErr: ErrNotFound,
		},
		{
			Name: "basic",
			Args: []string{"5", "x_1=", "$x_1", "$x_1", "6", "x_1=", "$x_1"},
			Want: []string{"5", "5", "6"},
		},
		{
			Name:    "illegal set 1",
			Args:    []string{"5", "=="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Name:    "illegal set 2",
			Args:    []string{"5", "$="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Name:    "illegal set 3",
			Args:    []string{"5", "1="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Name:    "illegal set 4",
			Args:    []string{"5", "-="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Name:    "clear unknown",
			Args:    []string{"x/"},
			WantErr: ErrNotFound,
		},
		{
			Name: "clear basic",
			Args: []string{"5", "x=", "x/"},
		},
		{
			Name:    "clear basic with check",
			Args:    []string{"5", "x=", "x/", "$x"},
			WantErr: ErrNotFound,
		},
	}
	UnitTestExecAll(t, data, nil)
}

func TestGetStringVariable(t *testing.T) {
	data := []struct {
		name    string
		args    []string
		vname   string
		want    string
		wantErr error
	}{
		{
			name:  "basic",
			args:  []string{"'foo'", "x="},
			vname: "x",
			want:  "foo",
		},
		{
			name:    "not found",
			vname:   "x",
			wantErr: ErrNotFound,
		},
		{
			name:    "not a string",
			args:    []string{"5", "x="},
			vname:   "x",
			wantErr: ErrExpectedAString,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init()
			err := r.Exec(d.args)
			if err != nil {
				t.Fatalf("err=%v, want nil", err)
			}
			got, err := r.GetStringVariable(d.vname)
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want %v", err, d.wantErr)
			}
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
		})
	}
}

func TestGetComplexVariable(t *testing.T) {
	data := []struct {
		name    string
		args    []string
		vname   string
		want    complex128
		wantErr error
	}{
		{
			name:  "basic",
			args:  []string{"5+i", "x="},
			vname: "x",
			want:  complex(5, 1),
		},
		{
			name:    "not found",
			vname:   "x",
			wantErr: ErrNotFound,
		},
		{
			name:  "integer",
			args:  []string{"5d", "x="},
			vname: "x",
			want:  complex(5, 0),
		},
		{
			name:    "string",
			args:    []string{"'foo'", "x="},
			vname:   "x",
			wantErr: ErrExpectedANumber,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init()
			err := r.Exec(d.args)
			if err != nil {
				t.Fatalf("err=%v, want nil", err)
			}
			got, err := r.GetComplexVariable(d.vname)
			if !errors.Is(err, d.wantErr) {
				t.Fatalf("err=%v, want %v", err, d.wantErr)
			}
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
		})
	}
}

func TestAllVariableNamesAndValues(t *testing.T) {
	var r RPN
	r.Init()
	r.Exec([]string{"1d", "a=", "2d", "b=", "vpush", "'foo'", "a=", "3d", "c="})
	want := []NameAndValues{
		{
			Name:   "a",
			Values: []Frame{{Type: INTEGER_FRAME, Int: 1}, {Type: STRING_FRAME, Str: "foo"}},
		},
		{
			Name:   "b",
			Values: []Frame{{Type: INTEGER_FRAME, Int: 2}},
		},
		{
			Name:   "c",
			Values: []Frame{{}, {Type: INTEGER_FRAME, Int: 3}},
		},
	}
	got := r.AllVariableNamesAndValues()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\n got=%+v\nwant=%+v", got, want)
	}
}

func TestExecVariableAsMacro(t *testing.T) {
	data := []UnitTestExecData{
		{
			Name:    "not found",
			Args:    []string{"@x"},
			WantErr: ErrNotFound,
		},
		{
			Name: "simple",
			Args: []string{"'1 2'", "x=", "@x"},
			Want: []string{"1", "2"},
		},
		{
			Name: "nested",
			Args: []string{"'1 2'", "x=", "'@x @x'", "y=", "@y"},
			Want: []string{"1", "2", "1", "2"},
		},
	}
	UnitTestExecAll(t, data, nil)
}
