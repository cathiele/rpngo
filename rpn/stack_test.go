package rpn

import (
	"testing"
)

func TestIsIntDefault(t *testing.T) {
	var f Frame
	if f.IsInt() {
		t.Error("IsInt() = true, want false")
	}
}

func TestBool(t *testing.T) {
	f := BoolFrame(false)
	if f.Bool() {
		t.Error("Bool() = true, want false")
	}
	f = BoolFrame(true)
	if !f.Bool() {
		t.Error("Bool() = false, want true")
	}
}

func TestString(t *testing.T) {
	data := []struct {
		name  string
		frame Frame
		quote bool
		want  string
	}{
		{
			name:  "default",
			frame: Frame{},
			want:  "nil",
		},
		{
			name:  "string",
			frame: Frame{Type: STRING_FRAME, Str: "foo"},
			want:  "foo",
		},
		{
			name:  "quoted string",
			frame: Frame{Type: STRING_FRAME, Str: "foo"},
			quote: true,
			want:  "\"foo\"",
		},
		{
			name:  "complex 1",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(-1, 0)},
			want:  "-1",
		},
		{
			name:  "complex 2",
			frame: Frame{Type: COMPLEX_FRAME},
			want:  "0",
		},
		{
			name:  "complex 3",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(123, 0)},
			want:  "123",
		},
		{
			name:  "complex 4",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(0, 1)},
			want:  "i",
		},
		{
			name:  "complex 5",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(0, -1)},
			want:  "-i",
		},
		{
			name:  "complex 6",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(12, 4)},
			want:  "12+4i",
		},
		{
			name:  "complex 7",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(-12, -4)},
			want:  "-12-4i",
		},
		{
			name:  "complex 8",
			frame: Frame{Type: COMPLEX_FRAME, Complex: complex(-1.2, -0.4)},
			want:  "-1.2-0.4i",
		},
		{
			name:  "bool true",
			frame: BoolFrame(true),
			want:  "true",
		},
		{
			name:  "bool false",
			frame: BoolFrame(false),
			want:  "false",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			got := d.frame.String(d.quote)
			if got != d.want {
				t.Errorf("want: %v, got: %v", d.want, got)
			}
		})
	}
}
