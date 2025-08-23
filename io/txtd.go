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

func puts(txtd TextDisplay, b byte) {
	if b == '\n' {
		newLine(txtd)
		return
	}
	txtd.Write([]byte{b})
}

func newLine(txtd TextDisplay) {
	txtd.SetY(txtd.Y() + 1)
	txtd.SetX(0)
}
