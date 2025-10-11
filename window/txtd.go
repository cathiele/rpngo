package window

func Print(tb *TextBuffer, msg string, updatenow bool) {
	PrintBytes(tb, []byte(msg), updatenow)
}

func PrintErr(tb *TextBuffer, err error, updatenow bool) {
	tb.TextColor(Red)
	Print(tb, err.Error(), updatenow)
	tb.Write('\n', updatenow)
	tb.TextColor(White)
}

func PrintBytes(tb *TextBuffer, msg []byte, updatenow bool) {
	for _, b := range msg {
		if err := tb.Write(b, updatenow); err != nil {
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
