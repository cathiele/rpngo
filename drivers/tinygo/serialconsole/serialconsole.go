// Package serialconsile provides a common mechanism for processsing
// characters from the UART
package serialconsole

import (
	"machine"
	"mattwach/rpngo/key"
)

type TermState int

const (
	NORMAL TermState = iota
	ESC
	ARROW
	PAGEUP
	PAGEDOWN
)

type SerialConsole struct {
	state TermState
}

func (sc *SerialConsole) GetChar() key.Key {
	c, err := machine.Serial.ReadByte()
	if err != nil {
		// nothing available
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
