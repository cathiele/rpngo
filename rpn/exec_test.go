package rpn

import (
	"errors"
	"reflect"
	"testing"
)

func TestExecInterrupted(t *testing.T) {
	var r RPN
	r.Init()
	r.Interrupt = make(chan bool, 1)
	r.Interrupt <- true
	err := r.exec("5")
	if !errors.Is(err, ErrInterrupted) {
		t.Errorf("err got %v, want %v", err, ErrInterrupted)
	}
}

func TestExec(t *testing.T) {
	data := []struct {
		name      string
		args      []string
		wantErr   error
		noFrame   bool
		wantFrame Frame
	}{
		{
			name:    "empty",
			noFrame: true,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init()
			err := r.Exec(d.args)
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err got %v, want %v", err, d.wantErr)
			}
			if err != nil {
				return
			}
			if d.noFrame {
				if len(r.frames) != 0 {
					t.Errorf("want no frames, len(r.frames)=%v", len(r.frames))
				}
			} else {
				if len(r.frames) == 0 {
					t.Fatalf("no stack frames with noFrame=false")
				}
				gotf := r.frames[len(r.frames)-1]
				if !reflect.DeepEqual(gotf, d.wantFrame) {
					t.Errorf("frame got %+v, want %+v", gotf, d.wantFrame)
				}
			}
		})
	}
}
