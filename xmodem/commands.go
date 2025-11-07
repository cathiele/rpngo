package xmodem

import (
	"errors"
	"mattwach/rpngo/rpn"
	"time"
)

var (
	errCRCMismatch          = errors.New("crc mismatch")
	errEndOfTransmission    = errors.New("end of transmission")
	errIncorrectPacketId    = errors.New("incorrect packet id")
	errIncorrectPacketIdSum = errors.New("incorrect packet id sum")
	errInvalidInitialPacket = errors.New("invalid initial packet type")
	errReadTimeout          = errors.New("read timeout")
	errNakSent              = errors.New("nak sent")
	errUnexpectedHeaderByte = errors.New("unexpected header byte")
)

type Serial interface {
	Open(path string) error
	Close() error
	ReadByte() (byte, error) // Non-blocking
	WriteByte(c byte) error
}

type XmodemCommands struct {
	serial       Serial
	attemptsLeft int
	state        readState
	deadline     time.Time
	nextPacketId uint8
	idx          int
	packet       [133]uint8
	buff         []byte
}

/*
func (sc *XmodemCommands) debugDump(r *rpn.RPN, err error) {
	r.Print("err: ")
	r.Println(err.Error())
	r.Println(fmt.Sprintf(
		"atmp-left: %v idx: %v state: %v nextId: %v len(buff): %v",
		sc.attemptsLeft,
		sc.idx,
		sc.state,
		sc.nextPacketId,
		len(sc.buff),
	))

	r.Println(fmt.Sprintf("phead: %02x %02x %02x", sc.packet[0], sc.packet[1], sc.packet[2]))
	r.Println(fmt.Sprintf("packet[3:16]: %v", string(sc.packet[3:16])))
	r.Println(fmt.Sprintf("packet-crc: %02x %02x", sc.packet[131], sc.packet[132]))
}
*/

const handshakeAttempt = 10

func (sc *XmodemCommands) InitAndRegister(r *rpn.RPN, serial Serial) {
	sc.serial = serial
	r.Register("xmodemr", sc.xmodemRead, rpn.CatIO, xmodemReadHelp)
	r.Register("xmodemw", sc.xmodemWrite, rpn.CatIO, xmodemWriteHelp)
}

const xmodemReadHelp = "Attempts to read data using the xmodem protocal to the top of the stack."

type readState uint8

const (
	InitialHandshake readState = iota
	ReadingPackets
	Finished
)

const (
	SOH uint8 = 0x01
	EOT       = 0x04
	ACK       = 0x06
	NAK       = 0x15
	ETB       = 0x17
	CAN       = 0x18
	C         = 0x43
)

func (sc *XmodemCommands) xmodemRead(r *rpn.RPN) error {
	sc.state = InitialHandshake
	sc.attemptsLeft = handshakeAttempt
	sc.nextPacketId = 0x01
	sc.buff = make([]byte, 0, 128)
	var err error
	for err == nil {
		switch sc.state {
		case InitialHandshake:
			err = sc.serial.WriteByte(C)
			if err == nil {
				err = sc.readPacket(r)
			}
		case ReadingPackets:
			err = sc.readPacket(r)
		case Finished:
			return r.PushFrame(rpn.StringFrame(sc.trimZeros(), rpn.STRING_BRACES))
		}
	}
	return err
}

func (sc *XmodemCommands) trimZeros() string {
	end := len(sc.buff) - 1
	for ; end > 0; end-- {
		if sc.buff[end] != 0 {
			break
		}
	}
	sc.buff = sc.buff[:end]
	return string(sc.buff)
}

func (sc *XmodemCommands) readPacket(r *rpn.RPN) error {
	for {
		if sc.state == InitialHandshake {
			sc.deadline = time.Now().Add(3 * time.Second)
		} else {
			sc.deadline = time.Now().Add(1 * time.Second)
		}

		sc.idx = 0
		for sc.idx < len(sc.packet) {
			if time.Now().After(sc.deadline) {
				return sc.nakWith(r, errReadTimeout)
			}
			// read in bursts so we are not checking the timer on every byte
			for n := 0; n < 1024; n++ {
				b, err := sc.serial.ReadByte()
				if err != nil {
					// We get errors when the buffer is empty, which is almost
					// guaranteed to happen.  Just rely on the timeout.
					continue
				}
				if (sc.idx == 0) && (sc.state == ReadingPackets) {
					if b == EOT {
						if err := sc.serial.WriteByte(ACK); err != nil {
							return err
						}
						sc.state = Finished
						return nil
					}
				}
				sc.packet[sc.idx] = b
				sc.idx++
				if sc.idx >= len(sc.packet) {
					break
				}
			}
		}

		if err := sc.validatePacket(); err != nil {
			return sc.nakWith(r, err)
		}

		if sc.packet[1] == sc.nextPacketId {
			sc.nextPacketId++
			sc.attemptsLeft = handshakeAttempt
			sc.buff = append(sc.buff, sc.packet[3:131]...)
			if sc.state == InitialHandshake {
				if sc.packet[0] != SOH {
					return errInvalidInitialPacket
				}
				sc.state = ReadingPackets
			}
			return sc.serial.WriteByte(ACK)
		} else if sc.packet[1] == (sc.nextPacketId - 1) {
			// sender might have not received our previous ack
			if err := sc.serial.WriteByte(ACK); err != nil {
				return err
			}
		} else {
			return errIncorrectPacketId
		}
	}

}

func (sc *XmodemCommands) nakWith(r *rpn.RPN, err error) error {
	//sc.debugDump(r, err)
	sc.serial.WriteByte(NAK)
	sc.attemptsLeft--
	if sc.attemptsLeft <= 0 {
		return err
	}
	return nil
}

func (sc *XmodemCommands) validatePacket() error {
	switch sc.packet[0] {
	case SOH, EOT, ETB, CAN:
		// ok
	default:
		return errUnexpectedHeaderByte
	}

	if uint16(sc.packet[1])+uint16(sc.packet[2]) != 0xFF {
		return errIncorrectPacketIdSum
	}

	var crc uint16
	for _, b := range sc.packet[3:131] {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}

	if (sc.packet[131] != uint8(crc>>8)) || (sc.packet[132] != uint8(crc&0xFF)) {
		return errCRCMismatch
	}

	return nil
}

const xmodemWriteHelp = "Attempts to send the data at the top of the stack using the xmodem protocol."

func (sc *XmodemCommands) xmodemWrite(r *rpn.RPN) error {
	return rpn.ErrNotSupported
}
