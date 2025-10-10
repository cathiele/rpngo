package window

func Print(txtd TextArea, msg string) {
	PrintBytes(txtd, []byte(msg))
}

func PrintErr(txtd TextArea, err error) {
	txtd.TextColor(Red)
	Print(txtd, err.Error())
	PutByte(txtd, '\n')
	txtd.TextColor(White)
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
	for x >= txtd.TextWidth() {
		y += 1
		if y >= txtd.TextHeight() {
			txtd.Scroll(-1)
			y--
		}
		x -= txtd.TextWidth()
	}
	for x < 0 {
		x += txtd.TextWidth()
		y -= 1
		if y < 0 {
			txtd.Scroll(1)
			y = 0
		}
	}
	txtd.SetCursorXY(x, y)
}
