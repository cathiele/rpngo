package window

func Print(txtd TextArea, msg string) {
	PrintBytes(txtd, []byte(msg))
}

func PrintErr(txtd TextArea, err error) {
	txtd.Color(31, 0, 0, 0, 0, 0)
	Print(txtd, err.Error())
	PutByte(txtd, '\n')
	txtd.Color(31, 31, 31, 0, 0, 0)
}

func PrintBytes(txtd TextArea, msg []byte) {
	for _, b := range msg {
		if err := txtd.Write(b); err != nil {
			return
		}
	}
}

func PutByte(txtd TextArea, b byte) {
	txtd.Write(b)
}

func Shift(txtd TextArea, n int) {
	x, y := txtd.CursorXY()
	x += n
	if x < 0 {
		x += txtd.TextWidth()
		y -= 1
		if y < 0 {
			y = 0
		}
	}
	txtd.SetCursorXY(x, y)
}
