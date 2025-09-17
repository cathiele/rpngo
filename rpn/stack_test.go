package rpn

import (
	"reflect"
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
		{
			name:  "integer",
			frame: Frame{Type: INTEGER_FRAME, Int: 1234},
			want:  "1234d",
		},
		{
			name:  "hex",
			frame: Frame{Type: HEXIDECIMAL_FRAME, Int: 0x1234},
			want:  "1234x",
		},
		{
			name:  "octal",
			frame: Frame{Type: OCTAL_FRAME, Int: 01234},
			want:  "1234o",
		},
		{
			name:  "binary",
			frame: Frame{Type: BINARY_FRAME, Int: 9},
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

func TestLengthAndClear(t *testing.T) {
	var r RPN
	r.Init()
	r.Exec([]string{"1", "2", "3"})
	if r.StackLen() != 3 {
		t.Errorf("StackLen()=%v, want 3", r.StackLen())
	}
	r.Clear()
	if r.StackLen() != 0 {
		t.Errorf("StackLen()=%v, want 0", r.StackLen())
	}
}

func TestPush(t *testing.T) {
	data := []struct {
		name string
		fn   func(r *RPN) error
		want Frame
	}{
		{
			name: "integer",
			fn:   func(r *RPN) error { return r.PushFrame(Frame{Type: INTEGER_FRAME, Int: 1234}) },
			want: Frame{Type: INTEGER_FRAME, Int: 1234},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init()
			err := d.fn(&r)
			if err != nil {
				t.Fatalf("err=%v, want nil", err)
			}
			if len(r.frames) != 1 {
				t.Fatalf("len(frames)=%v, want 1", len(r.frames))
			}
			if !reflect.DeepEqual(r.frames[0], d.want) {
				t.Errorf("frame mismatch. got=%+v, want=%+v", r.frames[0], d.want)
			}
		})
	}
}
