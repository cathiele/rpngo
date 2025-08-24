package io

import "mattwach/rpngo/io/key"

func getLine(input Input, txtd TextDisplay) (string, error) {
	print(txtd, "> ")
	var line []byte
	for {
		c, err := input.GetChar()
		if err != nil {
			return "", err
		}
		switch c {
		case key.KEY_LEFT:
			print(txtd, "LEFT")
			break
		default:
			b := byte(c)
			putbyte(txtd, b)
			if b == '\n' {
				newLine(txtd)
				return string(line), nil
			}
			line = append(line, b)
		}
	}
}
