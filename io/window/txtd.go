package window

func Print(txtd TextWindow, msg string) {
	PrintBytes(txtd, []byte(msg))
}

func PrintErr(txtd TextWindow, err error) {
	txtd.Color(31, 0, 0, 0, 0, 0)
	Print(txtd, "ERROR: ")
	Print(txtd, err.Error())
	PutByte(txtd, '\n')
	txtd.Color(31, 31, 31, 0, 0, 0)
}

func PrintBytes(txtd TextWindow, msg []byte) {
	for _, b := range msg {
		if err := txtd.Write(b); err != nil {
			return
		}
	}
}

func PutByte(txtd TextWindow, b byte) {
	txtd.Write(b)
}

func Shift(txtd TextWindow, n int) {
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
