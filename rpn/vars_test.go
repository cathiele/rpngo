package rpn

import (
	"reflect"
	"testing"
)

func testPushPopVariableFrame(t *testing.T) {
	var r RPN
	r.Init()
	r.PushString("bar")
	r.setVariable("foo")

	err := pushVariableFrame(&r)
	if err != nil {
		t.Fatalf("want err=nil, got %v", err)
	}
	r.PushString("x")
	r.setVariable("y")
	want := []map[string]Frame{
		{
			"foo": Frame{Type: STRING_FRAME, Str: "bar"},
		},
		{
			"y": Frame{Type: STRING_FRAME, Str: "x"},
		},
	}
	if !reflect.DeepEqual(want, r.variables) {
		t.Errorf("want variables = %+v, got %+v", want, r.variables)
	}

	err = popVariableFrame(&r)
	if err != nil {
		t.Fatalf("want err=nil, got %v", err)
	}
	want = []map[string]Frame{
		{
			"foo": Frame{Type: STRING_FRAME, Str: "bar"},
		},
	}
	if !reflect.DeepEqual(want, r.variables) {
		t.Errorf("want variables = %+v, got %+v", want, r.variables)
	}

	err = popVariableFrame(&r)
	if err != nil {
		t.Fatalf("want err=nil, got %v", err)
	}
	want = []map[string]Frame{}
	if !reflect.DeepEqual(want, r.variables) {
		t.Errorf("want variables = %+v, got %+v", want, r.variables)
	}

	err = popVariableFrame(&r)
	if err != ErrStackEmpty {
		t.Fatalf("want err=ErrStackEmpty, got %v", err)
	}
}
