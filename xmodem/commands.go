package xmodem

import (
	"errors"
	"mattwach/rpngo/rpn"
	"time"
)

var (
	errCRCMismatch          = errors.New("crc mismatch")
	errEarlyReceiverResponse = errors.New("receiver responded mid-packet")
	errEndOfTransmission    = errors.New("end of transmission")
	errIncorrectPacketId    = errors.New("incorrect packet id")
	errIncorrectPacketIdSum = errors.New("incorrect packet id sum")
	errInvalidInitialPacket = errors.New("invalid initial packet type")
	errReadTimeout          = errors.New("read timeout")
	errNakReceived              = errors.New("nak received")
	errNakSent              = errors.New("nak sent")
	errUnexpectedHeaderByte = errors.New("unexpected header byte")
	errTimeoutWaitingForStart = errors.New("timeout waiting for start message")
	errUnexpectedByteReceived = errors.New("unexpected byte received (not ack/nak)")
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
	r.Register("rx", sc.xmodemRead, rpn.CatIO, xmodemReadHelp)
	r.Register("sx", sc.xmodemWrite, rpn.CatIO, xmodemWriteHelp)
}

const xmodemReadHelp = "Attempts to read data using the xmodem protocal to the top of the stack."

type readState uint8

const (
	InitialHandshake readState = iota
	XferPackets
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
		case XferPackets:
			err = sc.readPacket(r)
		case Finished:
			return r.PushFrame(rpn.StringFrame(sc.trimExtra(), rpn.STRING_BRACES))
		}
	}
	return err
}

func (sc *XmodemCommands) xmodemWrite(r *rpn.RPN) error {
	s, err := r.PopFrame()
	if err != nil {
		return err
	}
	sc.buff = []byte(s)
	sc.state = InitialHandshake
	sc.attemptsLeft = handshakeAttempt
	sc.nextPacketId = 0x01
	sc.deadline = time.Now().Add(60 * time.Second)
	var err error
	for err == nil {
		switch sc.state {
		case InitialHandshake:
			b, err := sc.serial.ReadByte()
			if err == nil {
				if b == 'C' {
					sc.state = XferPackets
				}
			}
			if time.Now().After(sc.deadline) {
				return errTimeoutWaitingForStart
			}
		case XferPackets:
			err = sc.writePacket(r)
		case Finished:
			return nil
		}
	}
	return err
}

func (sc *XmodemCommands) trimExtra() string {
	// I see a string of 0x1A (SUB) characters in the last
	// block.  A bit surprising.
	end := len(sc.buff) - 1
	for ; end > 0; end-- {
		if (sc.buff[end] == '\n') || (sc.buff[end] >= 32) {
			break
		}
	}
	s := string(sc.buff[:end+1)
	sc.buff = nil
	return
}

func (sc *XmodemCommands) writePacket(r *rpn.RPN) error {
	lastPacket := sc.buildWritePacket()

	for _, b := range sc.packet {
		err := sc.serial.WriteByte(b)
		if err != nil {
			return err
		}
		c, err := sc.serial.ReadByte()
		if err == nil {
			// not expecting a byte yet
			return sc.writeRetry(errEarlyReceiverResponse)
		}
	}

	if err := sc.waitForAck(); err != nil {
		return sc.writeRetry(err)
	}

	if lastPacket {
		return sc.sendEOT()
	}

	return nil
}

func (sc *XmodemCommands) sendEOT() error {
	for {
		err := sc.serial.WriteByte(EOT)
		if err != nil {
			return err
		}
		err = sc.waitForAck()
		if err == nil {
			sc.state = Finished
			return nil
		}
		err = sc.writeRetry(err)
		if err != nil {
			return err
		}
	}
}


func (sc *XmodemCommands) waitForAck() error {
	deadline := time.Now().Add(3 * time.Seconds)
	for {
		c, err := sc.serial.ReadByte()
		if err == nil {
			switch c {
			case ACK:
				sc.attemptsLeft = handshakeAttempt
				sc.nextPacketId++
				return nil
			case NAK:
				return errNakReceived
			default:
				return errUnexpectedByteReceived
			}
		}
		if time.Now().After(deadline) {
			return errResponseTimeout
		}
	}
}

func (sc *XmodemCommands) buildWritePacket() bool {
	sc.packet[0] = SOH
	sc.packet[1] = sc.nextPacketId
	sc.packet[2] = 0xFF - sc.packet[1]
	start := (sc.nextPacketId - 1) * 128
	end := start + 128
	lastPacket := end > len(s.buff)
	if lastPacket {
		end = len(s.buff)
		// fill in some pad
		for i := 3; i < (3+128); i++ {
			sc.packet[i] = 0x1A
		}
	} 
	copy(sc.packet[3:], s.buff[start:end])
	crc := sc.calcPacketChecksum()
	sc.packet[131] = uint8(crc>>8)
	sc.packet[132] = uint8(crc&0xFF)
	return lastPacket
}

func (sc *XmodemCommands) writeRetry(err error) error {
	sc.attemptsLeft--
	if sc.attemptsLeft <= 0 {
		return err
	}
	return nil
}

func (sc *XmodemCommands) readPacket(r *rpn.RPN) error {
	for {
		if err := sc.readPacketData(); err != nil {
			return err
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
				sc.state = XferPackets
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

func (sc *XmodemCommands) readPacketData() error {
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
			if (sc.idx == 0) && (sc.state == XferPackets) {
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

	return nil
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

	crc := sc.calcPacketChecksum()
	if (sc.packet[131] != uint8(crc>>8)) || (sc.packet[132] != uint8(crc&0xFF)) {
		return errCRCMismatch
	}

	return nil
}

func (sc *XmodemCommands) calcPacketChecksum() uint16 {
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
	return crc
}

const xmodemWriteHelp = "Attempts to send the data at the top of the stack using the xmodem protocol."

func (sc *XmodemCommands) xmodemWrite(r *rpn.RPN) error {
	return rpn.ErrNotSupported
}
