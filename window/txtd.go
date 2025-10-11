package window

func Print(tb *TextBuffer, msg string) {
	PrintBytes(tb, []byte(msg))
}

func PrintErr(tb *TextBuffer, err error) {
	tb.TextColor(Red)
	Print(tb, err.Error())
	tb.Write('\n')
	tb.TextColor(White)
}

func PrintBytes(tb *TextBuffer, msg []byte) {
	for _, b := range msg {
		if err := tb.Write(b); err != nil {
			return
		}
	}
}

func Shift(tb *TextBuffer, n int) {
	x, y := tb.CursorXY()
	x += n
	for x >= tb.Txtw.TextWidth() {
		y += 1
		if y >= tb.Txtw.TextHeight() {
			tb.Scroll(-1)
			y--
		}
		x -= tb.Txtw.TextWidth()
	}
	for x < 0 {
		x += tb.Txtw.TextWidth()
		y -= 1
		if y < 0 {
			tb.Scroll(1)
			y = 0
		}
	}
	tb.SetCursorXY(x, y)
}
