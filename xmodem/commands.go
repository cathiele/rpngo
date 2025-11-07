package xmodem

import (
	"errors"
	"log"
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
	packet       [133]uint8
	buff         []byte
}

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
	for {
		switch sc.state {
		case InitialHandshake:
			if err := sc.serial.WriteByte(C); err != nil {
				return err
			}
			if err := sc.readPacket(); err != nil {
				return err
			}
		case ReadingPackets:
			if err := sc.readPacket(); err != nil {
				return err
			}
		case Finished:
			return r.PushFrame(rpn.StringFrame(string(sc.buff), rpn.STRING_BRACES))
		}
	}
}

func (sc *XmodemCommands) readPacket() error {
	log.Printf("readPacket. State=%v", sc.state)
	for {
		if sc.state == InitialHandshake {
			sc.deadline = time.Now().Add(3 * time.Second)
		} else {
			sc.deadline = time.Now().Add(1 * time.Second)
		}

		var idx int
		for idx < len(sc.packet) {
			if time.Now().After(sc.deadline) {
				log.Print("timeout")
				return sc.nakWith(errReadTimeout)
			}
			// read in bursts so we are not checking the timer on every byte
			for n := 0; n < 1024; n++ {
				b, err := sc.serial.ReadByte()
				if (idx == 0) && (sc.state == ReadingPackets) {
					if b == EOT {
						log.Print("EOT")
						return sc.serial.WriteByte(ACK)
					} else if b == ETB {
						if err := sc.serial.WriteByte(ACK); err != nil {
							return err
						}
						log.Print("ETB")
						sc.state = Finished
						return nil
					}
				}
				if err != nil {
					return err
				}
				if b != 0 {
					sc.packet[idx] = b
					idx++
					if idx >= len(sc.packet) {
						break
					}
				}
			}
		}

		if err := sc.validatePacket(); err != nil {
			return sc.nakWith(err)
		}

		if sc.packet[1] == sc.nextPacketId {
			log.Printf("Packet ok: id=%v", sc.packet[1])
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
			log.Printf("Packet repeated: id=%v", sc.packet[1])
			// sender might have not received our previous ack
			if err := sc.serial.WriteByte(ACK); err != nil {
				return err
			}
		} else {
			log.Printf("Unexpected packet: id=%v", sc.packet[1])
			return errIncorrectPacketId
		}
	}

}

func (sc *XmodemCommands) nakWith(err error) error {
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
		log.Printf("unexpted header byte %v", sc.packet[0])
		return errUnexpectedHeaderByte
	}

	if uint16(sc.packet[1])+uint16(sc.packet[2]) != 0xFF {
		log.Printf("bad sum 1=%v 2=%v", sc.packet[1], sc.packet[2])
		return errIncorrectPacketIdSum
	}

	var crc uint16
	for _, b := range sc.packet[4:131] {
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
		log.Printf("bad crc 0: want %x, got %v, 1: want %x, got %x",
			sc.packet[131], uint8(crc>>8), sc.packet[132], uint8(crc&0xFF))
		return errCRCMismatch
	}

	return nil
}

const xmodemWriteHelp = "Attempts to send the data at the top of the stack using the xmodem protocol."

func (sc *XmodemCommands) xmodemWrite(r *rpn.RPN) error {
	return rpn.ErrNotSupported
}
