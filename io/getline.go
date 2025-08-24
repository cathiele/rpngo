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
				idx = idx - 1
				shift(txtd, -1)
			}
		default:
			b := byte(c)
			if b == '\n' {
				newLine(txtd)
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
		putbyte(txtd, b)
	} else {
		line = append(line, 0) // grow the buffer
		copy(line[idx+1:], line[idx:])
		line[idx] = b
		print(txtd, string(line[idx:]))
		shift(txtd, -(len(line) - idx - 1))
	}
	return line
}
