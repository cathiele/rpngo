package rpn

import (
	"reflect"
	"testing"
)

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
	var r RPN
	r.Init()
	f := StringFrame("foo")
	err := r.PushFrame(f)
	if err != nil {
		t.Fatalf("err=%v, want nil", err)
	}
	if len(r.Frames) != 1 {
		t.Fatalf("len(frames)=%v, want 1", len(r.Frames))
	}
	if !reflect.DeepEqual(r.Frames[0], f) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", r.Frames[0], f)
	}
}

func TestPopFrame(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopFrame()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	f := StringFrame("foo")
	r.PushFrame(f)
	got, err := r.PopFrame()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if !reflect.DeepEqual(got, f) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, f)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
}

func TestPopStackIndex(t *testing.T) {
	var r RPN
	r.Init()
	_, err := r.PopStackIndex()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	r.PushFrame(StringFrame("foo"))
	_, err = r.PopStackIndex()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	r.PushFrame(IntFrame(-1, INTEGER_FRAME))
	_, err = r.PopStackIndex()
	if err != ErrIllegalValue {
		t.Errorf("err: %v, want: ErrIllegalValue", err)
	}
	r.PushFrame(IntFrame(4, INTEGER_FRAME))
	_, err = r.PopStackIndex()
	if err != ErrIllegalValue {
		t.Errorf("err: %v, want: ErrIllegalValue", err)
	}
	r.PushFrame(StringFrame("foo")) // 2
	r.PushFrame(StringFrame("bar")) // 1
	r.PushFrame(StringFrame("baz")) // 0
	r.PushFrame(ComplexFrame(complex(2, 0)))
	got, err := r.PopStackIndex()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	want := 2
	if want != got {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, want)
	}
	if len(r.Frames) != 3 {
		t.Errorf("stack size is %v, want 3", len(r.Frames))
	}
}

func TestPop2Frames(t *testing.T) {
	var r RPN
	r.Init()
	_, _, err := r.Pop2Frames()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	wanta := StringFrame("foo")
	r.PushFrame(wanta)
	wantb := StringFrame("bar")
	r.PushFrame(wantb)
	gota, gotb, err := r.Pop2Frames()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if !reflect.DeepEqual(gota, wanta) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", gota, wanta)
	}
	if !reflect.DeepEqual(gotb, wantb) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", gotb, wantb)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
}

func TestPop2Numbers(t *testing.T) {
	var r RPN
	r.Init()
	_, _, err := r.Pop2Numbers()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}

	r.PushFrame(StringFrame("foo"))
	r.PushFrame(StringFrame("bar"))
	_, _, err = r.Pop2Numbers()
	if err != ErrExpectedANumber {
		t.Errorf("err: %v, want: ErrExpectedANumber", err)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}

	wanta := IntFrame(123, INTEGER_FRAME)
	r.PushFrame(wanta)
	wantb := IntFrame(456, INTEGER_FRAME)
	r.PushFrame(wantb)
	gota, gotb, err := r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	wanta = ComplexFrame(123)
	wantb = ComplexFrame(complex(1, 2))
	r.PushFrame(IntFrame(123, INTEGER_FRAME))
	r.PushFrame(wantb)
	gota, gotb, err = r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	wantb = ComplexFrame(123)
	wanta = ComplexFrame(complex(1, 2))
	r.PushFrame(wanta)
	r.PushFrame(IntFrame(123, INTEGER_FRAME))
	gota, gotb, err = r.Pop2Numbers()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}

	wantb = ComplexFrame(123)
	wanta = ComplexFrame(complex(1, 2))
	r.PushFrame(wanta)
	r.PushFrame(wantb)
	_, _, err = r.Pop2Numbers()
	if !reflect.DeepEqual(wanta, gota) {
		t.Errorf("want: %v, go %v", wanta, gota)
	}
	if !reflect.DeepEqual(wantb, gotb) {
		t.Errorf("want: %v, go %v", wantb, gotb)
	}
}

func TestPeekFrame(t *testing.T) {
	var r RPN
	r.Init()
	r.PushFrame(StringFrame("foo"))
	r.PushFrame(StringFrame("bar"))
	_, err := r.PeekFrame(-1)
	if err != ErrIllegalValue {
		t.Errorf("want ErrIllegalValue, got %v", err)
	}
	_, err = r.PeekFrame(2)
	if err != ErrNotEnoughStackFrames {
		t.Errorf("want ErrNotEnoughStackFrames, got %v", err)
	}

	want := StringFrame("bar")
	got, err := r.PeekFrame(0)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, for %v", want, got)
	}

	want = StringFrame("foo")
	got, err = r.PeekFrame(1)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, for %v", want, got)
	}
}

func TestDeleteFrame(t *testing.T) {
	var r RPN
	r.Init()
	r.PushFrame(StringFrame("foo"))
	r.PushFrame(StringFrame("bar"))
	r.PushFrame(StringFrame("baz"))
	_, err := r.DeleteFrame(-1)
	if err != ErrIllegalValue {
		t.Errorf("want ErrIllegalValue, got %v", err)
	}
	_, err = r.DeleteFrame(3)
	if err != ErrNotEnoughStackFrames {
		t.Errorf("want ErrNotEnoughStackFrames, got %v", err)
	}

	got, err := r.DeleteFrame(0)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	want := StringFrame("baz")
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	wants := []Frame{
		StringFrame("foo"),
		StringFrame("bar"),
	}
	if !reflect.DeepEqual(wants, r.Frames) {
		t.Errorf("want %v, get %v", wants, r.Frames)
	}

	got, err = r.DeleteFrame(1)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	want = StringFrame("foo")
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	wants = []Frame{
		StringFrame("bar"),
	}
	if !reflect.DeepEqual(wants, r.Frames) {
		t.Errorf("want %v, get %v", wants, r.Frames)
	}
}

func TestInsertFrame(t *testing.T) {
	var r RPN
	r.Init()
	err := r.InsertFrame(Frame{}, -1)
	wantErr := ErrIllegalValue
	if err != wantErr {
		t.Errorf("err=%v, want %v", err, wantErr)
	}
	err = r.InsertFrame(Frame{}, 1)
	wantErr = ErrNotEnoughStackFrames
	if err != wantErr {
		t.Errorf("err=%v, want %v", err, wantErr)
	}
	r.PushFrame(StringFrame("1"))
	err = r.InsertFrame(StringFrame("foo"), 1)
	wantErr = nil
	if err != wantErr {
		t.Errorf("err=%v, want %v", err, wantErr)
	}

	wants := []Frame{StringFrame("foo"), StringFrame("1")}
	if !reflect.DeepEqual(wants, r.Frames) {
		t.Errorf("want %v, get %v", wants, r.Frames)
	}
}

func TestPushAndPopStack(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args:    []string{"spop"},
			WantErr: ErrStackEmpty,
		},
		{
			Args: []string{"1", "2", "spush"},
			Want: []string{"1", "2"},
		},
		{
			Args: []string{"1", "2", "spush", "5", "spop"},
			Want: []string{"1", "2"},
		},
	}
	UnitTestExecAll(t, data, nil)
}

func TestStackSize(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args: []string{"ssize"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"1", "ssize"},
			Want: []string{"1", "1d"},
		},
		{
			Args: []string{"1", "2", "ssize"},
			Want: []string{"1", "2", "2d"},
		},
	}
	UnitTestExecAll(t, data, nil)
}
