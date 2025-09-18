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
			name: "frame",
			fn:   func(r *RPN) error { return r.PushFrame(Frame{Type: INTEGER_FRAME, Int: 1234}) },
			want: Frame{Type: INTEGER_FRAME, Int: 1234},
		},
		{
			name: "complex",
			fn:   func(r *RPN) error { return r.PushComplex(complex(1, 2)) },
			want: Frame{Type: COMPLEX_FRAME, Complex: complex(1, 2)},
		},
		{
			name: "string",
			fn:   func(r *RPN) error { return r.PushString("foo") },
			want: Frame{Type: STRING_FRAME, Str: "foo"},
		},
		{
			name: "integer",
			fn:   func(r *RPN) error { return r.PushInt(1234, INTEGER_FRAME) },
			want: Frame{Type: INTEGER_FRAME, Int: 1234},
		},
		{
			name: "bool true",
			fn:   func(r *RPN) error { return r.PushBool(true) },
			want: Frame{Type: BOOL_FRAME, Int: 1},
		},
		{
			name: "bool false",
			fn:   func(r *RPN) error { return r.PushBool(false) },
			want: Frame{Type: BOOL_FRAME, Int: 0},
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

func TestPopFrame(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopFrame()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushString("foo")
	got, err := r.PopFrame()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want := Frame{Type: STRING_FRAME, Str: "foo"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.frames))
	}
}

func TestPopString(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopString()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushInt(1234, INTEGER_FRAME)
	_, err = r.PopString()
	if err != ErrExpectedAString {
		t.Errorf("err: %v, want: ErrExpectedAString", err)
	}
	r.PushString("foo")
	got, err := r.PopString()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want := "foo"
	if got != want {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
}

func TestPopBool(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopBool()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushInt(1234, INTEGER_FRAME)
	_, err = r.PopBool()
	if err != ErrExpectedABoolean {
		t.Errorf("err: %v, want: ErrExpectedABoolean", err)
	}
	r.PushBool(true)
	got, err := r.PopBool()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want := true
	if got != want {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
}

func TestPopNumber(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopNumber()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushString("foo")
	_, err = r.PopNumber()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	r.PushInt(1234, INTEGER_FRAME)
	got, err := r.PopNumber()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want := Frame{Type: INTEGER_FRAME, Int: 1234}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
	r.PushComplex(complex(1, 2))
	got, err = r.PopNumber()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want = Frame{Type: COMPLEX_FRAME, Complex: complex(1, 2)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
}

func TestPopComplex(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopComplex()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushString("foo")
	_, err = r.PopComplex()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	r.PushInt(1234, INTEGER_FRAME)
	got, err := r.PopComplex()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want := complex(1234, 0)
	if got != want {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
	r.PushComplex(complex(1, 2))
	got, err = r.PopComplex()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want = complex(1, 2)
	if got != want {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
}

func TestPopReal(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopReal()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushString("foo")
	_, err = r.PopReal()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	r.PushComplex(complex(1, 2))
	_, err = r.PopReal()
	if err != ErrComplexNumberNotSupported {
		t.Errorf("err: %v, want: ErrComplexNumberNotSupported", err)
	}
	r.PushComplex(complex(1234, 0))
	got, err := r.PopReal()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	var want float64 = 1234
	if got != want {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	r.PushInt(12, INTEGER_FRAME)
	got, err = r.PopReal()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want = 12
	if got != want {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
}

func TestPop2Frames(t *testing.T) {
	var r RPN
	r.Init()
	_, _, err := r.Pop2Frames()
	if err != ErrNotEnoughStackFrames {
		t.Errorf("err: %v, want: ErrNotEnoughStackFrames", err)
	}
	r.PushString("foo")
	_, _, err = r.Pop2Frames()
	if err != ErrNotEnoughStackFrames {
		t.Errorf("err: %v, want: ErrNotEnoughStackFrames", err)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
	r.PushString("bar")
	gota, gotb, err := r.Pop2Frames()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	wanta := Frame{Type: STRING_FRAME, Str: "foo"}
	wantb := Frame{Type: STRING_FRAME, Str: "bar"}
	if !reflect.DeepEqual(gota, wanta) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", gota, wanta)
	}
	if !reflect.DeepEqual(gotb, wantb) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", gotb, wantb)
	}
	if len(r.frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.frames))
	}
}

func TestPop2Numbers(t *testing.T) {
	var r RPN
	r.Init()
	_, _, err := r.Pop2Numbers()
	if err != ErrNotEnoughStackFrames {
		t.Errorf("err: %v, want: ErrNotEnoughStackFrames", err)
	}
	r.PushString("foo")
	_, _, err = r.Pop2Numbers()
	if err != ErrNotEnoughStackFrames {
		t.Errorf("err: %v, want: ErrNotEnoughStackFrames", err)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
	r.PushString("bar")
	_, _, err = r.Pop2Numbers()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}

	r.PushInt(123, INTEGER_FRAME)
	_, _, err = r.Pop2Numbers()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushInt(456, INTEGER_FRAME)
	gota, gotb, err := r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta := Frame{Type: INTEGER_FRAME, Int: 123}
	wantb := Frame{Type: INTEGER_FRAME, Int: 456}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	r.PushInt(123, INTEGER_FRAME)
	_, _, err = r.Pop2Numbers()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushComplex(complex(1, 2))
	gota, gotb, err = r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta = Frame{Type: COMPLEX_FRAME, Complex: complex(123, 0), Int: 123} // Int: 123 is ignored due to type
	wantb = Frame{Type: COMPLEX_FRAME, Complex: complex(1, 2)}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	r.PushComplex(complex(1, 2))
	_, _, err = r.Pop2Numbers()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushInt(123, INTEGER_FRAME)
	gota, gotb, err = r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta = Frame{Type: COMPLEX_FRAME, Complex: complex(1, 2)}
	wantb = Frame{Type: COMPLEX_FRAME, Complex: complex(123, 0), Int: 123} // Int: 123 is ignored due to type
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	r.PushComplex(complex(1, 2))
	_, _, err = r.Pop2Numbers()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushComplex(complex(3, 4))
	gota, gotb, err = r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta = Frame{Type: COMPLEX_FRAME, Complex: complex(1, 2)}
	wantb = Frame{Type: COMPLEX_FRAME, Complex: complex(3, 4)}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}
}

func TestPop2Ints(t *testing.T) {
	var r RPN
	r.Init()
	_, _, err := r.Pop2Ints()
	if err != ErrNotEnoughStackFrames {
		t.Errorf("err: %v, want: ErrNotEnoughStackFrames", err)
	}
	r.PushString("foo")
	_, _, err = r.Pop2Ints()
	if err != ErrNotEnoughStackFrames {
		t.Errorf("err: %v, want: ErrNotEnoughStackFrames", err)
	}
	if len(r.frames) != 1 {
		t.Errorf("stack size is %v, want 1", len(r.frames))
	}
	r.PushString("bar")
	_, _, err = r.Pop2Ints()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}

	r.PushInt(123, INTEGER_FRAME)
	_, _, err = r.Pop2Ints()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushInt(456, INTEGER_FRAME)
	gota, gotb, err := r.Pop2Ints()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta := Frame{Type: INTEGER_FRAME, Int: 123}
	wantb := Frame{Type: INTEGER_FRAME, Int: 456}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	r.PushInt(123, INTEGER_FRAME)
	_, _, err = r.Pop2Ints()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushComplex(complex(1, 2))
	gota, gotb, err = r.Pop2Ints()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta = Frame{Type: INTEGER_FRAME, Int: 123}
	wantb = Frame{Type: INTEGER_FRAME, Complex: complex(1, 2), Int: 1} // complex part is ignored
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	r.PushComplex(complex(1, 2))
	_, _, err = r.Pop2Ints()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushInt(123, INTEGER_FRAME)
	gota, gotb, err = r.Pop2Ints()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta = Frame{Type: INTEGER_FRAME, Complex: complex(1, 2), Int: 1} // complex part is ignored
	wantb = Frame{Type: INTEGER_FRAME, Int: 123}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	r.PushComplex(complex(1, 2))
	_, _, err = r.Pop2Ints()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.frames))
	}
	r.PushComplex(complex(3, 4))
	gota, gotb, err = r.Pop2Ints()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.frames) != 2 {
		t.Errorf("stack size is %v, want 2", len(r.frames))
	}
	wanta = Frame{Type: INTEGER_FRAME, Complex: complex(1, 2), Int: 1}
	wantb = Frame{Type: INTEGER_FRAME, Complex: complex(3, 4), Int: 3}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}
}
