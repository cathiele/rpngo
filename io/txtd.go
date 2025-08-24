package io

import "strings"

func print(txtd TextDisplay, msg string) {
	if len(msg) == 0 {
		return
	}
	idx := strings.Index(msg, "\n")
	if idx >= 0 {
		print(txtd, msg[0:idx])
		putByte(txtd, '\n')
		print(txtd, msg[idx+1:])
		return
	}
	for _, b := range []byte(msg) {
		if err := txtd.Write(b); err != nil {
			return
		}
	}
}

func putByte(txtd TextDisplay, b byte) {
	txtd.Write(b)
}

func shift(txtd TextDisplay, n int) {
	x, y := txtd.XY()
	x += n
	if x < 0 {
		x += txtd.Width()
		y -= 1
		if y < 0 {
			y = 0
		}
	}
	txtd.SetXY(x, y)
}
