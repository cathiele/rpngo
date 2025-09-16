package rpn

import "testing"

func TestInit(t *testing.T) {
  var r RPN
  r.Init()
  err := r.Exec([]string{"11", "22", "33", "ssize"})
  if err != nil {
    t.Fatalf("r.Exec() err=%v", err)
  }
  if len(r.frames) != 4 {
    t.Fatalf("want len(r.frames) = 4, got %v", len(r.frames))
  }
  f := r.frames[3]
  if f.Type != INTEGER_FRAME {
    t.Errorf("want frame.Type = INTEGER_FRAMFE, got %v", f.Type)
  }
  if f.Int != 3 {
    t.Errorf("want frame.Int = 3, got %v", f.Int)
  }
}

