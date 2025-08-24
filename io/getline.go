package io

import (
	"mattwach/rpngo/io/key"
)

const MAX_HISTORY_LINES = 100

type getLine struct {
	insertMode   bool
	input        Input
	txtd         TextDisplay
	history      [MAX_HISTORY_LINES]string
	historyCount int
}

func initGetLine(input Input, txtd TextDisplay) *getLine {
	return &getLine{
		insertMode:   true,
		input:        input,
		txtd:         txtd,
		historyCount: 0,
	}
}

func (gl *getLine) get() (string, error) {
	print(gl.txtd, "> ")
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
		case key.KEY_LEFT:
			if idx > 0 {
				idx--
				shift(gl.txtd, -1)
			}
		case key.KEY_RIGHT:
			if idx < len(line) {
				idx++
				shift(gl.txtd, 1)
			}
		case key.KEY_UP:
			if historyIdx < gl.historyCount && historyIdx <= MAX_HISTORY_LINES {
				historyIdx++
				line = gl.replaceLineWithHistory(
					historyIdx,
					len(line),
					idx)
				idx = len(line)
			}
		case key.KEY_DOWN:
			if historyIdx > 0 {
				historyIdx--
				line = gl.replaceLineWithHistory(
					historyIdx,
					len(line),
					idx)
				idx = len(line)
			}
		case key.KEY_BACKSPACE:
			if idx > 0 {
				idx--
				line = delete(line, idx)
				shift(gl.txtd, -1)
				printBytes(gl.txtd, line[idx:])
				putByte(gl.txtd, ' ')
				shift(gl.txtd, -(len(line) - idx + 1))
			}
		case key.KEY_DEL:
			if idx < len(line) {
				line = delete(line, idx)
				printBytes(gl.txtd, line[idx:])
				putByte(gl.txtd, ' ')
				shift(gl.txtd, -(len(line) - idx + 1))
			}
		case key.KEY_INS:
			gl.insertMode = !gl.insertMode
		case key.KEY_END:
			shift(gl.txtd, len(line)-idx)
			idx = len(line)
		case key.KEY_HOME:
			shift(gl.txtd, -idx)
			idx = 0
		case key.KEY_EOF:
			return "exit", nil
		default:
			b := byte(c)
			if b == '\n' {
				putByte(gl.txtd, b)
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
		putByte(gl.txtd, b)
	} else if gl.insertMode {
		line = append(line, 0) // grow the buffer
		copy(line[idx+1:], line[idx:])
		line[idx] = b
		print(gl.txtd, string(line[idx:]))
		shift(gl.txtd, -(len(line) - idx - 1))
	} else {
		line[idx] = b
		putByte(gl.txtd, b)
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
	shift(gl.txtd, -idx)
	for i := 0; i < oldlen; i++ {
		putByte(gl.txtd, ' ')
	}
	shift(gl.txtd, -oldlen)
	print(gl.txtd, newl)
	return []byte(newl)
}

func delete(line []byte, idx int) []byte {
	return append(line[:idx], line[idx+1:]...)
}
