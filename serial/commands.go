package serial

import (
	"errors"
	"mattwach/rpngo/rpn"
	"time"
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

type SerialCommands struct {
	serial Serial
}

func (sc *SerialCommands) InitAndRegister(r *rpn.RPN, serial Serial) {
	sc.serial = serial
	r.Register("serialr", sc.serialRead, rpn.CatIO, serialReadHelp)
	r.Register("serialw", sc.serialWrite, rpn.CatIO, serialWriteHelp)
}

const serialReadHelp = "Opens the $.serial device (if applicable).  Waits up to 30 seconds " +
	"for data.  Once data starts, the command will exit after 250ms passes with no new data"

func (sc *SerialCommands) serialRead(r *rpn.RPN) error {
	// Some platforms (such as pico) do not need this variable set.
	// So we'll try the open with whatever it's set to and respond if the
	// open fails.
	port, _ := r.GetStringVariable(".serial")
	if err := sc.serial.Open(port); err != nil {
		return err
	}
	defer sc.serial.Close()
	c, err := sc.readInitialByte(r)
	if err != nil {
		return err
	}
	data, err := sc.readRemainingBytes(r, c)
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(data, rpn.STRING_BRACES))
}

func (sc *SerialCommands) readInitialByte(r *rpn.RPN) (byte, error) {
	deadline := time.Now().Add(time.Second * 30)
	for time.Now().Before(deadline) {
		if r.Interrupt() {
			return 0, rpn.ErrInterrupted
		}
		c, err := sc.serial.ReadByte()
		if err != nil {
			return 0, err
		}
		if c != 0 {
			return c, nil
		} else {
			time.Sleep(20 * time.Millisecond)
		}
	}
	return 0, errReadTimeout
}

func (sc *SerialCommands) readRemainingBytes(r *rpn.RPN, firstChar byte) (string, error) {
	data := make([]byte, 0, 16)
	data[0] = firstChar
	deadline := time.Now().Add(time.Millisecond * 250)
	for time.Now().Before(deadline) {
		if r.Interrupt() {
			return "", rpn.ErrInterrupted
		}
		c, err := sc.serial.ReadByte()
		if err != nil {
			return "", err
		}
		if c != 0 {
			data = append(data, c)
			deadline = time.Now().Add(time.Millisecond * 250)
		}
	}
	return string(data), nil
}

const serialWriteHelp = "Opens the $.serial device (if applicable). Writes the top stack element as a string"

func (sc *SerialCommands) serialWrite(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	data := f.String(false)
	// Some platforms (such as pico) do not need this variable set.
	// So we'll try the open with whatever it's set to and respond if the
	// open fails.
	port, _ := r.GetStringVariable(".serial")
	if err := sc.serial.Open(port); err != nil {
		return err
	}
	defer sc.serial.Close()
	for _, c := range data {
		if r.Interrupt() {
			return rpn.ErrInterrupted
		}
		if err := sc.serial.WriteByte(byte(c)); err != nil {
			return err
		}
	}
	return nil
}
