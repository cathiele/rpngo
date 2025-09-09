package functions

import (
	"mattwach/rpngo/rpn"
)

const AndHelp = "Performs a logical AND operation"

func And(r *rpn.RPN) error {
  a, b, err := pop2Bins(r)
  if err != nil {
    return err
  }
  return r.PushInt(a.Int & b.Int, a.Type)
}

const OrHelp = "Performs a logical OR operation"

func Or(r *rpn.RPN) error {
  a, b, err := pop2Bins(r)
  if err != nil {
    return err
  }
  return r.PushInt(a.Int | b.Int, a.Type)
}

const XOrHelp = "Performs a logical XOR operation"

func XOr(r *rpn.RPN) error {
  a, b, err := pop2Bins(r)
  if err != nil {
    return err
  }
  return r.PushInt(a.Int ^ b.Int, a.Type)
}

const ShiftLeftHelp = "Performs a logical shift left operation"

func ShiftLeft(r *rpn.RPN) error {
  a, b, err := pop2Bins(r)
  if err != nil {
    return err
  }
  return r.PushInt(a.Int << b.Int, a.Type)
}

const ShiftRightHelp = "Performs a logical shift right operation"

func ShiftRight(r *rpn.RPN) error {
  a, b, err := pop2Bins(r)
  if err != nil {
    return err
  }
  return r.PushInt(a.Int >> b.Int, a.Type)
}

func pop2Bins(r *rpn.RPN) (a rpn.Frame, b rpn.Frame, err error) {
  a, b, err = r.Pop2Ints()
  if err != nil {
    return
  }
  if (a.Int < 0) || (b.Int < 0) {
    r.PushFrame(a)
    r.PushFrame(b)
    err = rpn.ErrExpectedAPositiveNumber
  }
  return
}
