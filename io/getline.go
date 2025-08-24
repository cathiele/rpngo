package io

import (
	"mattwach/rpngo/io/key"
)

type getLine struct {
	insertMode bool
	input      Input
	txtd       TextDisplay
}

func (gl *getLine) get() (string, error) {
	print(gl.txtd, "> ")
	var line []byte
	idx := 0
	for {
		c, err := gl.input.GetChar()
		if err != nil {
			return "", err
		}
		switch c {
		case key.KEY_LEFT:
			if idx > 0 {
				idx--
				shift(gl.txtd, -1)
			}
		case key.KEY_RIGHT:
			if idx < len(line) {
				idx++
				shift(gl.txtd, 1)
			}
		case key.KEY_BACKSPACE:
			if idx > 0 {
				idx--
				line = delete(line, idx)
				shift(gl.txtd, -1)
				printBytes(gl.txtd, line[idx:])
				putByte(gl.txtd, ' ')
				shift(gl.txtd, -(len(line) - idx + 1))
			}
		case key.KEY_DEL:
			if idx < len(line) {
				line = delete(line, idx)
				printBytes(gl.txtd, line[idx:])
				putByte(gl.txtd, ' ')
				shift(gl.txtd, -(len(line) - idx + 1))
			}
		case key.KEY_INS:
			gl.insertMode = !gl.insertMode
		case key.KEY_END:
			shift(gl.txtd, len(line)-idx)
			idx = len(line)
		case key.KEY_HOME:
			shift(gl.txtd, -idx)
			idx = 0
		case key.KEY_EOF:
			return "exit", nil
		default:
			b := byte(c)
			if b == '\n' {
				putByte(gl.txtd, b)
				return string(line), nil
			}
			line = gl.addChar(line, idx, b)
			idx++
		}
		gl.txtd.Refresh()
	}
}

func (gl *getLine) addChar(line []byte, idx int, b byte) []byte {
	if idx >= len(line) {
		line = append(line, b)
		putByte(gl.txtd, b)
	} else if gl.insertMode {
		line = append(line, 0) // grow the buffer
		copy(line[idx+1:], line[idx:])
		line[idx] = b
		print(gl.txtd, string(line[idx:]))
		shift(gl.txtd, -(len(line) - idx - 1))
	} else {
		line[idx] = b
		putByte(gl.txtd, b)
	}
	return line
}

func delete(line []byte, idx int) []byte {
	return append(line[:idx], line[idx+1:]...)
}
