package serial

import (
	"errors"
	"fmt"
	"mattwach/rpngo/rpn"
	"os"
)

var (
	errSerialPortNotOpen      = errors.New("serial port not open")
	errSerialPortNeedsToBeSet = errors.New("$.serial needs to be set (e.g. /dev/ttyAMC0)")
)

type readData struct {
	read chan byte
	err  chan error
	done chan bool
}

type Serial struct {
	f     *os.File
	rdata *readData
}

func (sc *Serial) Open(r *rpn.RPN) error {
	if sc.f != nil {
		return fmt.Errorf("serial file is already open: v", sc.f.Name())
	}
	f, err := r.GetVariable(".serial")
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	sc.f, err = os.OpenFile(f.UnsafeString(), os.O_RDWR|os.O_SYNC, 0666)
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
	if sc.rdata != nil {
		<-sc.rdata.done
		sc.rdata = nil
	}
	sc.f = nil
	return err
}

var readbuff = make([]byte, 1)

func (sc *Serial) ReadByte() (byte, error) {
	if sc.f == nil {
		return 0, errSerialPortNotOpen
	}

	if sc.rdata == nil {
		sc.rdata = &readData{
			read: make(chan byte, 16),
			err:  make(chan error, 1),
			done: make(chan bool),
		}
		go func() {
			for {
				_, err := sc.f.Read(readbuff)
				if err != nil {
					sc.rdata.err <- err
					break
				}
				sc.rdata.read <- readbuff[0]
			}
			sc.rdata.done <- true
		}()
	}

	select {
	case rd := <-sc.rdata.read:
		return rd, nil
	case err := <-sc.rdata.err:
		return 0, err
	default:
		return 0, nil
	}
}

func (sc *Serial) WriteByte(c byte) error {
	if sc.f == nil {
		return errSerialPortNotOpen
	}
	_, err := sc.f.Write([]byte{c})
	return err
}
