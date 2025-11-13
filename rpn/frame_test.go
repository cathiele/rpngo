package rpn

import (
	"errors"
	"testing"
)

func TestIs(t *testing.T) {
	data := []struct {
		frame Frame
		fn    func(Frame) bool
		want  bool
	}{
		{
			frame: IntFrame(1, INTEGER_FRAME),
			fn:    func(f Frame) bool { return f.IsInt() },
			want:  true,
		},
		{
			frame: ComplexFrame(1, COMPLEX_FRAME),
			fn:    func(f Frame) bool { return f.IsInt() },
			want:  false,
		},
		{
			frame: ComplexFrame(1, COMPLEX_FRAME),
			fn:    func(f Frame) bool { return f.IsComplex() },
			want:  true,
		},
		{
			frame: ComplexFrame(1, POLAR_FRAME),
			fn:    func(f Frame) bool { return f.IsComplex() },
			want:  true,
		},
		{
			frame: IntFrame(1, INTEGER_FRAME),
			fn:    func(f Frame) bool { return f.IsComplex() },
			want:  false,
		},
		{
			frame: ComplexFrame(1, COMPLEX_FRAME),
			fn:    func(f Frame) bool { return f.IsNumber() },
			want:  true,
		},
		{
			frame: ComplexFrame(1, POLAR_FRAME),
			fn:    func(f Frame) bool { return f.IsNumber() },
			want:  true,
		},
		{
			frame: IntFrame(1, INTEGER_FRAME),
			fn:    func(f Frame) bool { return f.IsNumber() },
			want:  true,
		},
		{
			frame: BoolFrame(true),
			fn:    func(f Frame) bool { return f.IsNumber() },
			want:  false,
		},
		{
			frame: BoolFrame(true),
			fn:    func(f Frame) bool { return f.IsBool() },
			want:  true,
		},
		{
			frame: BoolFrame(false),
			fn:    func(f Frame) bool { return f.IsBool() },
			want:  true,
		},
		{
			frame: IntFrame(1, INTEGER_FRAME),
			fn:    func(f Frame) bool { return f.IsBool() },
			want:  false,
		},
		{
			frame: StringFrame("foo", STRING_DOUBLE_QUOTE),
			fn:    func(f Frame) bool { return f.IsString() },
			want:  true,
		},
		{
			frame: IntFrame(1, INTEGER_FRAME),
			fn:    func(f Frame) bool { return f.IsString() },
			want:  false,
		},
	}

	for _, d := range data {
		t.Run(d.frame.String(false), func(t *testing.T) {
			got := d.fn(d.frame)
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
		})
	}
}

func TestComplex(t *testing.T) {
	data := []struct {
		frame   Frame
		wantErr error
		want    complex128
	}{
		{
			frame: ComplexFrame(complex(1, 1), COMPLEX_FRAME),
			want:  complex(1, 1),
		},
		{
			frame: ComplexFrame(complex(1, 1), POLAR_FRAME),
			want:  complex(1, 1),
		},
		{
			frame: IntFrame(1, INTEGER_FRAME),
			want:  complex(1, 0),
		},
		{
			frame:   BoolFrame(true),
			wantErr: ErrExpectedANumber,
		},
		{
			frame:   StringFrame("foo", STRING_DOUBLE_QUOTE),
			wantErr: ErrExpectedANumber,
		},
	}

	for _, d := range data {
		t.Run(d.frame.String(false), func(t *testing.T) {
			got, err := d.frame.Complex()
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want=%v", err, d.wantErr)
			}
		})
	}
}

func TestReal(t *testing.T) {
	data := []struct {
		frame   Frame
		wantErr error
		want    float64
	}{
		{
			frame: RealFrame(1),
			want:  1,
		},
		{
			frame: IntFrame(1, INTEGER_FRAME),
			want:  1,
		},
		{
			frame:   ComplexFrame(complex(1, 1), COMPLEX_FRAME),
			wantErr: ErrComplexNumberNotSupported,
		},
		{
			frame:   ComplexFrame(complex(1, 1), POLAR_FRAME),
			wantErr: ErrComplexNumberNotSupported,
		},
		{
			frame:   BoolFrame(true),
			wantErr: ErrExpectedANumber,
		},
		{
			frame:   StringFrame("foo", STRING_DOUBLE_QUOTE),
			wantErr: ErrExpectedANumber,
		},
	}

	for _, d := range data {
		t.Run(d.frame.String(false), func(t *testing.T) {
			got, err := d.frame.Real()
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want=%v", err, d.wantErr)
			}
		})
	}
}

func TestInt(t *testing.T) {
	data := []struct {
		frame   Frame
		wantErr error
		want    int64
	}{
		{
			frame: IntFrame(1, INTEGER_FRAME),
			want:  1,
		},
		{
			frame: RealFrame(1),
			want:  1,
		},
		{
			frame:   ComplexFrame(complex(1, 1), COMPLEX_FRAME),
			wantErr: ErrComplexNumberNotSupported,
		},
		{
			frame:   ComplexFrame(complex(1, 1), POLAR_FRAME),
			wantErr: ErrComplexNumberNotSupported,
		},
		{
			frame:   BoolFrame(true),
			wantErr: ErrExpectedANumber,
		},
		{
			frame:   StringFrame("foo", STRING_DOUBLE_QUOTE),
			wantErr: ErrExpectedANumber,
		},
	}

	for _, d := range data {
		t.Run(d.frame.String(false), func(t *testing.T) {
			got, err := d.frame.Int()
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want=%v", err, d.wantErr)
			}
		})
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
			frame:   StringFrame("true", STRING_DOUBLE_QUOTE),
			wantErr: ErrExpectedABoolean,
		},
		{
			frame:   IntFrame(1, INTEGER_FRAME),
			wantErr: ErrExpectedABoolean,
		},
		{
			frame:   ComplexFrame(1, COMPLEX_FRAME),
			wantErr: ErrExpectedABoolean,
		},
		{
			frame:   ComplexFrame(1, POLAR_FRAME),
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
			frame: StringFrame("foo", STRING_DOUBLE_QUOTE),
			want:  "foo",
		},
		{
			name:  "quoted string",
			frame: StringFrame("foo", STRING_DOUBLE_QUOTE),
			quote: true,
			want:  "\"foo\"",
		},
		{
			name:  "complex 1",
			frame: ComplexFrame(-1, COMPLEX_FRAME),
			want:  "-1",
		},
		{
			name:  "complex 2",
			frame: ComplexFrame(0, COMPLEX_FRAME),
			want:  "0",
		},
		{
			name:  "complex 3",
			frame: ComplexFrame(123, COMPLEX_FRAME),
			want:  "123",
		},
		{
			name:  "complex 4",
			frame: ComplexFrame(complex(0, 1), COMPLEX_FRAME),
			want:  "i",
		},
		{
			name:  "complex 5",
			frame: ComplexFrame(complex(0, -1), COMPLEX_FRAME),
			want:  "-i",
		},
		{
			name:  "complex 6",
			frame: ComplexFrame(complex(12, 4), COMPLEX_FRAME),
			want:  "12+4i",
		},
		{
			name:  "complex 7",
			frame: ComplexFrame(complex(-12, -4), COMPLEX_FRAME),
			want:  "-12-4i",
		},
		{
			name:  "complex 8",
			frame: ComplexFrame(complex(-1.2, -0.4), COMPLEX_FRAME),
			want:  "-1.2-0.4i",
		},
		{
			name:  "polar 1",
			frame: ComplexFrame(-1, POLAR_FRAME),
			want:  "1<3.141592653589793",
		},
		{
			name:  "polar 2",
			frame: ComplexFrame(0, POLAR_FRAME),
			want:  "0<0",
		},
		{
			name:  "polar 3",
			frame: ComplexFrame(complex(0, 1), POLAR_FRAME),
			want:  "1<1.570796326794897",
		},
		{
			name:  "polar 4",
			frame: ComplexFrame(complex(0, -1), POLAR_FRAME),
			want:  "1<-1.570796326794897",
		},
		{
			name:  "polar 5",
			frame: PolarFrame2(1, 1),
			want:  "1<1",
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
