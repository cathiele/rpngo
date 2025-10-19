package rpn

import (
	"errors"
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
			wantFrame:  ComplexFrame(10),
		},
		{
			name:       "complex 2",
			args:       []string{"-10"},
			frameCount: 1,
			wantFrame:  ComplexFrame(-10),
		},
		{
			name:       "complex 3",
			args:       []string{"i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(0, 1)),
		},
		{
			name:       "complex 4",
			args:       []string{"-i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(0, -1)),
		},
		{
			name:       "complex 5",
			args:       []string{"10i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(0, 10)),
		},
		{
			name:       "complex 6",
			args:       []string{"-10i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(0, -10)),
		},
		{
			name:       "complex 7",
			args:       []string{".2+3i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(0.2, 3)),
		},
		{
			name:       "complex 8",
			args:       []string{"2-30.2i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(2, -30.2)),
		},
		{
			name:       "complex 9",
			args:       []string{"-20+0.3i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(-20, 0.3)),
		},
		{
			name:       "complex 10",
			args:       []string{"-0.2-.3i"},
			frameCount: 1,
			wantFrame:  ComplexFrame(complex(-0.2, -0.3)),
		},
		{
			name:       "integer 1",
			args:       []string{"10d"},
			frameCount: 1,
			wantFrame:  IntFrame(10, INTEGER_FRAME),
		},
		{
			name:       "integer 2",
			args:       []string{"-6d"},
			frameCount: 1,
			wantFrame:  IntFrame(-6, INTEGER_FRAME),
		},
		{
			name:       "hex 1",
			args:       []string{"10x"},
			frameCount: 1,
			wantFrame:  IntFrame(0x10, HEXIDECIMAL_FRAME),
		},
		{
			name:       "hex 2",
			args:       []string{"-10x"},
			frameCount: 1,
			wantFrame:  IntFrame(-0x10, HEXIDECIMAL_FRAME),
		},
		{
			name:       "octal 1",
			args:       []string{"10o"},
			frameCount: 1,
			wantFrame:  IntFrame(010, OCTAL_FRAME),
		},
		{
			name:       "octal 2",
			args:       []string{"-10o"},
			frameCount: 1,
			wantFrame:  IntFrame(-010, OCTAL_FRAME),
		},
		{
			name:       "binary 1",
			args:       []string{"10b"},
			frameCount: 1,
			wantFrame:  IntFrame(0b10, BINARY_FRAME),
		},
		{
			name:       "binary 2",
			args:       []string{"-10b"},
			frameCount: 1,
			wantFrame:  IntFrame(-0b10, BINARY_FRAME),
		},
		{
			name:       "empty string double quote",
			args:       []string{"\"\""},
			frameCount: 1,
			wantFrame:  StringFrame(""),
		},
		{
			name:       "empty string single quote",
			args:       []string{"''"},
			frameCount: 1,
			wantFrame:  StringFrame(""),
		},
		{
			name:       "string double quote",
			args:       []string{"\"foo\""},
			frameCount: 1,
			wantFrame:  StringFrame("foo"),
		},
		{
			name:       "string single quote",
			args:       []string{"'foo'"},
			frameCount: 1,
			wantFrame:  StringFrame("foo"),
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
			wantFrame:  IntFrame(0, INTEGER_FRAME),
		},
		{
			name:       "set and get variable",
			args:       []string{"55d", "foo=", "$foo"},
			frameCount: 1,
			wantFrame:  IntFrame(55, INTEGER_FRAME),
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
			wantFrame:  IntFrame(55, INTEGER_FRAME),
		},
		{
			name:       "conversion (parse check only)",
			args:       []string{"0", "mi>km"},
			frameCount: 1,
			wantFrame:  ComplexFrame(0),
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
			if len(r.Frames) != d.frameCount {
				t.Errorf("frame count want %v, got %v", d.frameCount, len(r.Frames))
			}
			if len(r.Frames) > 0 {
				gotf := r.Frames[len(r.Frames)-1]
				if !reflect.DeepEqual(gotf, d.wantFrame) {
					t.Errorf("frame got %+v, want %+v", gotf, d.wantFrame)
				}
			}
		})
	}
}

const float64EqualityThreshold = 1e-4
