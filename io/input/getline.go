package input

import (
	"mattwach/rpngo/io/window"
)

const MAX_HISTORY_LINES = 100

type getLine struct {
	insertMode   bool
	input        Input
	txtd         window.TextWindow
	history      [MAX_HISTORY_LINES]string
	historyCount int
}

func initGetLine(input Input, txtd window.TextWindow) *getLine {
	return &getLine{
		insertMode:   true,
		input:        input,
		txtd:         txtd,
		historyCount: 0,
	}
}

func (gl *getLine) get() (string, error) {
	window.Print(gl.txtd, "> ")
	var line []byte
	idx := 0
	// how many steps back into history, with 0 being not in history
	historyIdx := 0
	for {
		c, err := gl.input.GetChar()
		if err != nil {
			return "", err
		}
		switch c {
		case KEY_LEFT:
			if idx > 0 {
				idx--
				window.Shift(gl.txtd, -1)
			}
		case KEY_RIGHT:
			if idx < len(line) {
				idx++
				window.Shift(gl.txtd, 1)
			}
		case KEY_UP:
			if historyIdx < gl.historyCount && historyIdx <= MAX_HISTORY_LINES {
				historyIdx++
				line = gl.replaceLineWithHistory(
					historyIdx,
					len(line),
					idx)
				idx = len(line)
			}
		case KEY_DOWN:
			if historyIdx > 0 {
				historyIdx--
				line = gl.replaceLineWithHistory(
					historyIdx,
					len(line),
					idx)
				idx = len(line)
			}
		case KEY_BACKSPACE:
			if idx > 0 {
				idx--
				line = delete(line, idx)
				window.Shift(gl.txtd, -1)
				window.PrintBytes(gl.txtd, line[idx:])
				window.PutByte(gl.txtd, ' ')
				window.Shift(gl.txtd, -(len(line) - idx + 1))
			}
		case KEY_DEL:
			if idx < len(line) {
				line = delete(line, idx)
				window.PrintBytes(gl.txtd, line[idx:])
				window.PutByte(gl.txtd, ' ')
				window.Shift(gl.txtd, -(len(line) - idx + 1))
			}
		case KEY_INS:
			gl.insertMode = !gl.insertMode
		case KEY_END:
			window.Shift(gl.txtd, len(line)-idx)
			idx = len(line)
		case KEY_HOME:
			window.Shift(gl.txtd, -idx)
			idx = 0
		case KEY_EOF:
			return "exit", nil
		default:
			b := byte(c)
			if b == '\n' {
				window.PutByte(gl.txtd, b)
				s := string(line)
				gl.addToHistory(s)
				return s, nil
			}
			line = gl.addChar(line, idx, b)
			idx++
		}
		gl.txtd.Refresh()
	}
}

func (gl *getLine) addChar(line []byte, idx int, b byte) []byte {
	if idx >= len(line) {
		line = append(line, b)
		window.PutByte(gl.txtd, b)
	} else if gl.insertMode {
		line = append(line, 0) // grow the buffer
		copy(line[idx+1:], line[idx:])
		line[idx] = b
		window.Print(gl.txtd, string(line[idx:]))
		window.Shift(gl.txtd, -(len(line) - idx - 1))
	} else {
		line[idx] = b
		window.PutByte(gl.txtd, b)
	}
	return line
}

func (gl *getLine) addToHistory(line string) {
	// if the last history element is the same as line, don't repeat it
	if gl.historyCount > 0 && gl.history[(gl.historyCount-1)%MAX_HISTORY_LINES] == line {
		return
	}
	gl.history[gl.historyCount%MAX_HISTORY_LINES] = line
	gl.historyCount++
}

func (gl *getLine) replaceLineWithHistory(historyIdx int, oldlen int, idx int) []byte {
	newl := gl.history[(gl.historyCount-historyIdx)%MAX_HISTORY_LINES]
	// remove the existing line
	window.Shift(gl.txtd, -idx)
	for i := 0; i < oldlen; i++ {
		window.PutByte(gl.txtd, ' ')
	}
	window.Shift(gl.txtd, -oldlen)
	window.Print(gl.txtd, newl)
	return []byte(newl)
}

func delete(line []byte, idx int) []byte {
	return append(line[:idx], line[idx+1:]...)
}
