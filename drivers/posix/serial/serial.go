//go:build !pico && !pico2

package serial

import (
	"errors"
	"fmt"
	"os"
)

var errSerialPortNotOpen = errors.New("serial port not open")

type Serial struct {
	f *os.File
}

func (sc *Serial) Open(path string) error {
	if sc.f != nil {
		return fmt.Errorf("serial file is already open: v", sc.f.Name())
	}
	var err error
	sc.f, err = os.OpenFile(path, os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (sc *Serial) Close() error {
	if sc.f == nil {
		return errSerialPortNotOpen
	}
	err := sc.f.Close()
	sc.f = nil
	return err
}

var readbuff = make([]byte 1)
func (sc *Serial) ReadByte() (byte, error) {
	if sc.f == nil {
		return 0, errSerialPortNotOpen
	}
	// TODO: This probably needs to be non-blocking so that
	// a timeout can be implemented.
	n, err := sc.f.Read(readbuff)
	return readbuff[0], err
}

func (sc *Serial) WriteByte(c byte) error {
	if sc.f == nil {
		return errSerialPortNotOpen
	}
	_, err := sc.f.Write([]byte{c})
	return err
}
