//go:build pico || pico2

// Package serialconsile provides a common mechanism for processsing
// characters from the UART
package serial

import (
	"machine"
	"mattwach/rpngo/key"
	"time"
)

type TermState int

const (
	NORMAL TermState = iota
	ESC
	ARROW
	PAGEUP
	PAGEDOWN
)

type Serial struct {
	state       TermState
	ignoreUntil time.Time
}

func (sc *Serial) Init() {
	sc.state = NORMAL
	sc.ignoreUntil = time.Now().Add(time.Second)
}

func (sc *Serial) Open(path string) error {
	return nil
}

func (sc *Serial) Close() error {
	return nil
}

func (sc *Serial) ReadByte() (byte, error) {
	return machine.Serial.ReadByte()
}

func (sc *Serial) WriteByte(c byte) error {
	return machine.Serial.WriteByte(c)
}

func (sc *Serial) GetChar() key.Key {
	c, err := machine.Serial.ReadByte()
	if err != nil {
		// nothing available
		return 0
	}
	if time.Now().Before(sc.ignoreUntil) {
		// This might be some random junk from starting up cold
		return 0
	}
	switch sc.state {
	case NORMAL:
		switch c {
		case 13:
			return '\n'
		case 27:
			sc.state = ESC
		case 127:
			return key.KEY_BACKSPACE
		default:
			return key.Key(c)
		}
	case ESC:
		switch c {
		case '[':
			sc.state = ARROW
		default:
			sc.state = NORMAL
		}
	case ARROW:
		sc.state = NORMAL
		switch c {
		case 'A':
			return key.KEY_UP
		case 'B':
			return key.KEY_DOWN
		case 'C':
			return key.KEY_RIGHT
		case 'D':
			return key.KEY_LEFT
		case '5':
			sc.state = PAGEUP
		case '6':
			sc.state = PAGEDOWN
		}
	case PAGEUP:
		sc.state = NORMAL
		if c == 126 {
			return key.KEY_PAGEUP
		}
	case PAGEDOWN:
		sc.state = NORMAL
		if c == 126 {
			return key.KEY_PAGEDOWN
		}
	}
	return key.Key(c)
}
