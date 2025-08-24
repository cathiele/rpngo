package io

import (
	"mattwach/rpngo/io/key"
)

func getLine(input Input, txtd TextDisplay) (string, error) {
	print(txtd, "> ")
	var line []byte
	idx := 0
	for {
		c, err := input.GetChar()
		if err != nil {
			return "", err
		}
		switch c {
		case key.KEY_LEFT:
			if idx > 0 {
				idx--
				shift(txtd, -1)
			}
		case key.KEY_RIGHT:
			if idx < len(line) {
				idx++
				shift(txtd, 1)
			}
		case key.KEY_BACKSPACE:
			if idx > 0 {
				idx--
				line = delete(line, idx)
				shift(txtd, -1)
				printBytes(txtd, line[idx:])
				putByte(txtd, ' ')
				shift(txtd, -(len(line) - idx + 1))
			}
		case key.KEY_DEL:
			if idx < len(line) {
				line = delete(line, idx)
				printBytes(txtd, line[idx:])
				putByte(txtd, ' ')
				shift(txtd, -(len(line) - idx + 1))
			}
		default:
			b := byte(c)
			if b == '\n' {
				putByte(txtd, b)
				return string(line), nil
			}
			line = addOrInsert(line, idx, b, txtd)
			idx++
		}
		txtd.Refresh()
	}
}

func addOrInsert(line []byte, idx int, b byte, txtd TextDisplay) []byte {
	if idx >= len(line) {
		line = append(line, b)
		putByte(txtd, b)
	} else {
		line = append(line, 0) // grow the buffer
		copy(line[idx+1:], line[idx:])
		line[idx] = b
		print(txtd, string(line[idx:]))
		shift(txtd, -(len(line) - idx - 1))
	}
	return line
}

func delete(line []byte, idx int) []byte {
	return append(line[:idx], line[idx+1:]...)
}
