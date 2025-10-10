package input

import (
	"log"
	"mattwach/rpngo/key"
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
	"os"
	"path/filepath"
	"strings"
)

const MAX_HISTORY_LINES = 500

type getLine struct {
	insertMode     bool
	input          Input
	txtd           window.TextWindow
	history        [MAX_HISTORY_LINES]string
	historyCount   int
	historyFile    *os.File
	namesAndValues []rpn.NameAndValues
}

const histFile = ".rpngo_history"

func initGetLine(input Input, txtd window.TextWindow) *getLine {
	gl := &getLine{ // object allocated on the heap: object size 4028 exceeds maximum stack allocation size 256
		insertMode:     true,
		input:          input,
		txtd:           txtd,
		historyCount:   0,
		namesAndValues: make([]rpn.NameAndValues, 0, 16),
	}
	gl.loadHistory()
	gl.prepareHistory()
	return gl
}

func historyPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, histFile), nil
}

func (gl *getLine) loadHistory() {
	path, err := historyPath()
	if err != nil {
		log.Printf("Could not generate history path for load: %v", err) // object allocated on the heap (OK)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Could not read hitory file: %v", err) // object allocated on the heap (OK)
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) > 0 {
			gl.history[gl.historyCount%MAX_HISTORY_LINES] = line
			gl.historyCount++
		}
	}
}

func (gl *getLine) prepareHistory() {
	path, err := historyPath()
	if err != nil {
		log.Printf("Could not generate history path for prepare: %v", err) // object allocated on the heap (OK)
		return
	}
	gl.historyFile, err = os.Create(path)
	if err != nil {
		log.Printf("Could not create history path: %v", err) // object allocated on the heap (OK)
		return
	}
	mini := gl.historyCount - MAX_HISTORY_LINES
	if mini < 0 {
		mini = 0
	}
	for i := mini; i < gl.historyCount; i++ {
		line := gl.history[i%MAX_HISTORY_LINES] + "\n"
		_, err := gl.historyFile.Write([]byte(line))
		if err != nil {
			log.Printf("error writing exsiting history: %v", err) // object allocated on the heap: escapes at line 85
		}
	}
}

func (gl *getLine) get(r *rpn.RPN) (string, error) {
	gl.txtd.Cursor(true)
	defer gl.txtd.Cursor(false)
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
				window.Shift(gl.txtd, -1)
			}
		case key.KEY_RIGHT:
			if idx < len(line) {
				idx++
				window.Shift(gl.txtd, 1)
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
				window.Shift(gl.txtd, -1)
				window.PrintBytes(gl.txtd, line[idx:])
				window.PutByte(gl.txtd, ' ')
				window.Shift(gl.txtd, -(len(line) - idx + 1))
			}
		case key.KEY_DEL:
			if idx < len(line) {
				line = delete(line, idx)
				window.PrintBytes(gl.txtd, line[idx:])
				window.PutByte(gl.txtd, ' ')
				window.Shift(gl.txtd, -(len(line) - idx + 1))
			}
		case key.KEY_INS:
			gl.insertMode = !gl.insertMode
		case key.KEY_END:
			window.Shift(gl.txtd, len(line)-idx)
			idx = len(line)
		case key.KEY_HOME:
			window.Shift(gl.txtd, -idx)
			idx = 0
		case '\t':
			line, idx = gl.tabComplete(r, line, idx)
		case key.KEY_EOF:
			return "exit", nil
		case key.KEY_F1:
			return gl.execMacro(r, idx, "@.f1")
		case key.KEY_F2:
			return gl.execMacro(r, idx, "@.f2")
		case key.KEY_F3:
			return gl.execMacro(r, idx, "@.f3")
		case key.KEY_F4:
			return gl.execMacro(r, idx, "@.f4")
		case key.KEY_F5:
			return gl.execMacro(r, idx, "@.f5")
		case key.KEY_F6:
			return gl.execMacro(r, idx, "@.f6")
		case key.KEY_F7:
			return gl.execMacro(r, idx, "@.f7")
		case key.KEY_F8:
			return gl.execMacro(r, idx, "@.f8")
		case key.KEY_F9:
			return gl.execMacro(r, idx, "@.f9")
		case key.KEY_F10:
			return gl.execMacro(r, idx, "@.f10")
		case key.KEY_F11:
			return gl.execMacro(r, idx, "@.f11")
		case key.KEY_F12:
			return gl.execMacro(r, idx, "@.f12")
		default:
			b := byte(c)
			if b == '\n' {
				window.Shift(gl.txtd, len(line)-idx)
				window.PutByte(gl.txtd, b)
				s := string(line)
				gl.addToHistory(s)
				return s, nil
			}
			line = gl.addChar(line, idx, b)
			idx++
		}
	}
}

func (gl *getLine) execMacro(r *rpn.RPN, idx int, name string) (string, error) {
	window.Shift(gl.txtd, -idx)
	for idx > 0 {
		window.PutByte(gl.txtd, ' ')
		idx--
	}
	return name, nil
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
	if gl.historyFile != nil {
		line = line + "\n"
		_, err := gl.historyFile.WriteString(line)
		if err != nil {
			log.Printf("could not write history line: %v", err) // object allocated on the heap (OK)
		}
		gl.historyFile.Sync()
	}
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
