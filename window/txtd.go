package window

func Print(tb *TextBuffer, msg string, updatenow bool) {
	for _, b := range msg {
		if err := tb.Write(byte(b), updatenow); err != nil {
			return
		}
	}
}

func PrintErr(tb *TextBuffer, err error, updatenow bool) {
	tb.TextColor(Red)
	Print(tb, err.Error(), updatenow)
	tb.Write('\n', updatenow)
	tb.TextColor(White)
}

func DrawStr(tw TextWindow, x, y int, msg string, col ColorChar) {
	w, h := tw.TextSize()
	for _, c := range msg {
		if (x >= 0) && (x < w) && (y >= 0) && (y < h) {
			tw.DrawChar(x, y, col|ColorChar(c))
		}
		x++
		if (x >= w) || c == '\n' {
			x = 0
			y++
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
