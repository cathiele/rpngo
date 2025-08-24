package io

import "strings"

func print(txtd TextDisplay, msg string) {
	if len(msg) == 0 {
		return
	}
	idx := strings.Index(msg, "\n")
	if idx >= 0 {
		print(txtd, msg[0:idx])
		newLine(txtd)
		print(txtd, msg[idx+1:])
		return
	}
	txtd.Write([]byte(msg))
}

func putbyte(txtd TextDisplay, b byte) {
	if b == '\n' {
		newLine(txtd)
		return
	}
	txtd.Write([]byte{b})
}

func newLine(txtd TextDisplay) {
	y := txtd.Y()
	h := txtd.Height()
	if y >= (h - 1) {
		txtd.Scroll(1)
		txtd.SetY(h - 1)
	} else {
		txtd.SetY(y + 1)
	}
	txtd.SetX(0)
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
