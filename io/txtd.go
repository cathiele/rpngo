package io

func print(txtd TextDisplay, msg string) {
	printBytes(txtd, []byte(msg))
}

func printErr(txtd TextDisplay, err error) {
	txtd.Color(31, 0, 0, 0, 0, 0)
	print(txtd, "ERROR: ")
	print(txtd, err.Error())
	putByte(txtd, '\n')
	txtd.Color(31, 31, 31, 0, 0, 0)
}

func printBytes(txtd TextDisplay, msg []byte) {
	for _, b := range msg {
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
