package xmodem

import (
	"errors"
	"mattwach/rpngo/rpn"
)

var (
	errReadTimeout = errors.New("read timeout")
)

type Serial interface {
	Open(path string) error
	Close() error
	ReadByte() (byte, error) // Non-blocking
	WriteByte(c byte) error
}

type XmodemCommands struct {
	serial Serial
}

func (sc *XmodemCommands) InitAndRegister(r *rpn.RPN, serial Serial) {
	sc.serial = serial
	r.Register("xmodemr", sc.xmodemRead, rpn.CatIO, xmodemReadHelp)
	r.Register("xmodemw", sc.xmodemWrite, rpn.CatIO, xmodemWriteHelp)
}

const xmodemReadHelp = "Attempts to read data using the xmodem protocal to the top of the stack."

func (sc *XmodemCommands) xmodemRead(r *rpn.RPN) error {
	return rpn.ErrNotSupported
}

const xmodemWriteHelp = "Attempts to send th data at the top of the stack using the xmodem protocol."

func (sc *XmodemCommands) xmodemWrite(r *rpn.RPN) error {
	return rpn.ErrNotSupported
}
