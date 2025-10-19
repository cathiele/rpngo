package rpn

import (
	"errors"
	"testing"
)

func TestIsIntDefault(t *testing.T) {
	var f Frame
	if f.IsInt() {
		t.Error("IsInt() = true, want false")
	}
}

func TestBool(t *testing.T) {
	data := []struct {
		frame   Frame
		wantErr error
		want    bool
	}{
		{
			frame: BoolFrame(true),
			want:  true,
		},
		{
			frame: BoolFrame(false),
			want:  false,
		},
		{
			frame:   StringFrame("true"),
			wantErr: ErrExpectedABoolean,
		},
		{
			frame:   IntFrame(1, INTEGER_FRAME),
			wantErr: ErrExpectedABoolean,
		},
		{
			frame:   ComplexFrame(1),
			wantErr: ErrExpectedABoolean,
		},
	}

	for _, d := range data {
		t.Run(d.frame.String(false), func(t *testing.T) {
			got, err := d.frame.Bool()
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want=%v", err, d.wantErr)
			}
		})
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
			frame: StringFrame("foo"),
			want:  "foo",
		},
		{
			name:  "quoted string",
			frame: StringFrame("foo"),
			quote: true,
			want:  "\"foo\"",
		},
		{
			name:  "complex 1",
			frame: ComplexFrame(-1),
			want:  "-1",
		},
		{
			name:  "complex 2",
			frame: ComplexFrame(0),
			want:  "0",
		},
		{
			name:  "complex 3",
			frame: ComplexFrame(123),
			want:  "123",
		},
		{
			name:  "complex 4",
			frame: ComplexFrame(complex(0, 1)),
			want:  "i",
		},
		{
			name:  "complex 5",
			frame: ComplexFrame(complex(0, -1)),
			want:  "-i",
		},
		{
			name:  "complex 6",
			frame: ComplexFrame(complex(12, 4)),
			want:  "12+4i",
		},
		{
			name:  "complex 7",
			frame: ComplexFrame(complex(-12, -4)),
			want:  "-12-4i",
		},
		{
			name:  "complex 8",
			frame: ComplexFrame(complex(-1.2, -0.4)),
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
		{
			name:  "integer",
			frame: IntFrame(1234, INTEGER_FRAME),
			want:  "1234d",
		},
		{
			name:  "hex",
			frame: IntFrame(0x1234, HEXIDECIMAL_FRAME),
			want:  "1234x",
		},
		{
			name:  "octal",
			frame: IntFrame(01234, OCTAL_FRAME),
			want:  "1234o",
		},
		{
			name:  "binary",
			frame: IntFrame(9, BINARY_FRAME),
			want:  "1001b",
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
