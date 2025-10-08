package rpn

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

func TestExecInterrupted(t *testing.T) {
	var r RPN
	r.Init()
	r.Interrupt = func() bool { return true }
	err := r.exec("5")
	if !errors.Is(err, ErrInterrupted) {
		t.Errorf("err got %v, want %v", err, ErrInterrupted)
	}
}

func TestExec(t *testing.T) {
	data := []struct {
		name       string
		args       []string
		wantErr    error
		frameCount int
		wantFrame  Frame
	}{
		{
			name: "empty",
		},
		{
			name:    "unknown fn",
			args:    []string{"foo"},
			wantErr: ErrSyntax,
		},
		{
			name:       "complex 1",
			args:       []string{"10"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(10, 0)},
		},
		{
			name:       "complex 2",
			args:       []string{"-10"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(-10, 0)},
		},
		{
			name:       "complex 3",
			args:       []string{"i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(0, 1)},
		},
		{
			name:       "complex 4",
			args:       []string{"-i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(0, -1)},
		},
		{
			name:       "complex 5",
			args:       []string{"10i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(0, 10)},
		},
		{
			name:       "complex 6",
			args:       []string{"-10i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(0, -10)},
		},
		{
			name:       "complex 7",
			args:       []string{".2+3i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(0.2, 3)},
		},
		{
			name:       "complex 8",
			args:       []string{"2-30.2i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(2, -30.2)},
		},
		{
			name:       "complex 9",
			args:       []string{"-20+0.3i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(-20, 0.3)},
		},
		{
			name:       "complex 10",
			args:       []string{"-0.2-.3i"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(-0.2, -0.3)},
		},
		{
			name:       "integer 1",
			args:       []string{"10d"},
			frameCount: 1,
			wantFrame:  Frame{Type: INTEGER_FRAME, Int: 10},
		},
		{
			name:       "integer 2",
			args:       []string{"-6d"},
			frameCount: 1,
			wantFrame:  Frame{Type: INTEGER_FRAME, Int: -6},
		},
		{
			name:       "hex 1",
			args:       []string{"10x"},
			frameCount: 1,
			wantFrame:  Frame{Type: HEXIDECIMAL_FRAME, Int: 0x10},
		},
		{
			name:       "hex 2",
			args:       []string{"-10x"},
			frameCount: 1,
			wantFrame:  Frame{Type: HEXIDECIMAL_FRAME, Int: -0x10},
		},
		{
			name:       "octal 1",
			args:       []string{"10o"},
			frameCount: 1,
			wantFrame:  Frame{Type: OCTAL_FRAME, Int: 010},
		},
		{
			name:       "octal 2",
			args:       []string{"-10o"},
			frameCount: 1,
			wantFrame:  Frame{Type: OCTAL_FRAME, Int: -010},
		},
		{
			name:       "binary 1",
			args:       []string{"10b"},
			frameCount: 1,
			wantFrame:  Frame{Type: BINARY_FRAME, Int: 2},
		},
		{
			name:       "binary 2",
			args:       []string{"-10b"},
			frameCount: 1,
			wantFrame:  Frame{Type: BINARY_FRAME, Int: -2},
		},
		{
			name:       "empty string double quote",
			args:       []string{"\"\""},
			frameCount: 1,
			wantFrame:  Frame{Type: STRING_FRAME, Str: ""},
		},
		{
			name:       "empty string single quote",
			args:       []string{"''"},
			frameCount: 1,
			wantFrame:  Frame{Type: STRING_FRAME, Str: ""},
		},
		{
			name:       "string double quote",
			args:       []string{"\"foo\""},
			frameCount: 1,
			wantFrame:  Frame{Type: STRING_FRAME, Str: "foo"},
		},
		{
			name:       "string single quote",
			args:       []string{"'foo'"},
			frameCount: 1,
			wantFrame:  Frame{Type: STRING_FRAME, Str: "foo"},
		},
		{
			name:    "string mismatched quotes 1",
			args:    []string{"\"foo'"},
			wantErr: ErrSyntax,
		},
		{
			name:    "string mismatched quotes 2",
			args:    []string{"'foo\""},
			wantErr: ErrSyntax,
		},
		{
			name:    "string mismatched quotes 3",
			args:    []string{"\""},
			wantErr: ErrSyntax,
		},
		{
			name:    "string mismatched quotes 4",
			args:    []string{"'"},
			wantErr: ErrSyntax,
		},
		{
			name:       "some fn",
			args:       []string{"ssize"},
			frameCount: 1,
			wantFrame:  Frame{Type: INTEGER_FRAME},
		},
		{
			name:       "set and get variable",
			args:       []string{"55d", "foo=", "$foo"},
			frameCount: 1,
			wantFrame:  Frame{Type: INTEGER_FRAME, Int: 55},
		},
		{
			name:    "set variable no value",
			args:    []string{"foo="},
			wantErr: ErrStackEmpty,
		},
		{
			name:    "set empty var",
			args:    []string{"="},
			wantErr: ErrSyntax,
		},
		{
			name:    "set and clear variable",
			args:    []string{"55d", "foo=", "foo/", "$foo"},
			wantErr: ErrNotFound,
		},
		{
			name:    "clear unknown variable",
			args:    []string{"foo/"},
			wantErr: ErrNotFound,
		},
		{
			name:    "unknown variable",
			args:    []string{"$foo"},
			wantErr: ErrNotFound,
		},
		{
			name:       "set and execute macro 1",
			args:       []string{"'1 2 55d'", "foo=", "@foo"},
			frameCount: 3,
			wantFrame:  Frame{Type: INTEGER_FRAME, Int: 55},
		},
		{
			name:       "conversion (parse check only)",
			args:       []string{"0", "mi>km"},
			frameCount: 1,
			wantFrame:  Frame{Type: COMPLEX_FRAME, Complex: complex(0, 0)},
		},
		{
			name: "help all",
			args: []string{"?"},
		},
		{
			name: "help one",
			args: []string{"spush?"},
		},
		{
			name:    "help unknown",
			args:    []string{"foo?"},
			wantErr: ErrNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init()
			r.Print = func(string) {}
			err := r.Exec(d.args)
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err got %v, want %v", err, d.wantErr)
			}
			if len(r.frames) != d.frameCount {
				t.Errorf("frame count want %v, got %v", d.frameCount, len(r.frames))
			}
			if len(r.frames) > 0 {
				gotf := r.frames[len(r.frames)-1]
				if !reflect.DeepEqual(gotf, d.wantFrame) {
					t.Errorf("frame got %+v, want %+v", gotf, d.wantFrame)
				}
			}
		})
	}
}

const float64EqualityThreshold = 1e-4

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
