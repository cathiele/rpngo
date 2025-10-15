// Handles input heystrokes from the user. e.g.an "editor"
//
// Many of the implementation choices below are to avoid heap allocations
// on microcontrollers.  Notably the use of []byte to hold the result
// as well as the extra code to keep from converting between bytes
// and strings (which forces a heap allocation).  Heap allocations are
// not forbidden of course, as this is an interactive session, but since
// line entry is character-based anyway, using a []byte as a baseline
// makes sense even if heap allocations were not a concern.
package input

import (
	"bytes"
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
	txtb           *window.TextBuffer
	history        [MAX_HISTORY_LINES][]byte
	historyCount   int
	historyFile    *os.File
	namesAndValues []rpn.NameAndValues
	// line is the current line.  It's kept here to support entering
	// scrolling mode without losing the current line contents.
	line []byte
}

const histFile = ".rpngo_history"

func initGetLine(input Input, txtb *window.TextBuffer) *getLine {
	gl := &getLine{ // object allocated on the heap: (OK)
		insertMode:     true,
		input:          input,
		txtb:           txtb,
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
		if len(line) == 0 {
			continue
		}
		hidx := gl.historyCount % MAX_HISTORY_LINES
		if len(gl.history[hidx]) > 0 {
			// history has wrapped the internal buffer
			gl.history[hidx] = gl.history[hidx][:0]
		}
		for _, c := range line {
			gl.history[hidx] = append(gl.history[hidx], byte(c))
		}
		gl.historyCount++
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
		line := gl.history[i%MAX_HISTORY_LINES]
		_, err := gl.historyFile.Write(line)
		if err != nil {
			log.Printf("error writing exsiting history: %v", err)
		}
		_, err = gl.historyFile.Write([]byte{'\n'})
		if err != nil {
			log.Printf("error writing exsiting history cr: %v", err)
		}
	}
}

func (gl *getLine) get(r *rpn.RPN) (string, error) {
	gl.txtb.Cursor(true)
	defer gl.txtb.Cursor(false)
	gl.line = gl.line[:0]
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
				gl.txtb.Shift(-1)
			}
		case key.KEY_RIGHT:
			if idx < len(gl.line) {
				idx++
				gl.txtb.Shift(1)
			}
		case key.KEY_UP:
			if historyIdx < gl.historyCount && historyIdx <= MAX_HISTORY_LINES {
				historyIdx++
				gl.replaceLineWithHistory(historyIdx, idx)
				idx = len(gl.line)
			}
		case key.KEY_DOWN:
			if historyIdx > 0 {
				historyIdx--
				gl.replaceLineWithHistory(historyIdx, idx)
				idx = len(gl.line)
			}
		case key.KEY_BACKSPACE:
			if idx > 0 {
				idx--
				gl.line = delete(gl.line, idx)
				gl.txtb.Shift(-1)
				gl.txtb.PrintBytes(gl.line[idx:], true)
				gl.txtb.Write(' ', true)
				gl.txtb.Shift(-(len(gl.line) - idx + 1))
			}
		case key.KEY_DEL:
			if idx < len(gl.line) {
				gl.line = delete(gl.line, idx)
				gl.txtb.PrintBytes(gl.line[idx:], true)
				gl.txtb.Write(' ', true)
				gl.txtb.Shift(-(len(gl.line) - idx + 1))
			}
		case key.KEY_INS:
			gl.insertMode = !gl.insertMode
		case key.KEY_END:
			gl.txtb.Shift(len(gl.line) - idx)
			idx = len(gl.line)
		case key.KEY_HOME:
			gl.txtb.Shift(-idx)
			idx = 0
		case '\t':
			idx = gl.tabComplete(r, idx)
		case 27: // ESCAPE key
			gl.enterScrollingMode(0)
		case key.KEY_PAGEUP:
			gl.enterScrollingMode(-gl.pageDelta())
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
				gl.txtb.Shift(len(gl.line) - idx)
				gl.txtb.Write(b, true)
				gl.addToHistory()
				return string(gl.line), nil
			}
			gl.addChar(idx, b)
			idx++
		}
	}
}

func (gl *getLine) execMacro(r *rpn.RPN, idx int, name string) (string, error) {
	gl.txtb.Shift(-idx)
	for idx > 0 {
		gl.txtb.Write(' ', true)
		idx--
	}
	return name, nil
}

func (gl *getLine) addChar(idx int, b byte) {
	if idx >= len(gl.line) {
		gl.line = append(gl.line, b)
		gl.txtb.Write(b, true)
	} else if gl.insertMode {
		gl.line = append(gl.line, 0) // grow the buffer
		copy(gl.line[idx+1:], gl.line[idx:])
		gl.line[idx] = b
		gl.txtb.Print(string(gl.line[idx:]), true)
		gl.txtb.Shift(-(len(gl.line) - idx - 1))
	} else {
		gl.line[idx] = b
		gl.txtb.Write(b, true)
	}
}

func (gl *getLine) addToHistory() {
	// if the last history element is the same as line, don't repeat it
	if gl.historyCount > 0 && bytes.Equal(gl.history[(gl.historyCount-1)%MAX_HISTORY_LINES], gl.line) {
		return
	}
	if len(gl.line) == 0 {
		// line is empty
		return
	}
	hidx := gl.historyCount % MAX_HISTORY_LINES
	if len(gl.history[hidx]) > 0 {
		gl.history[hidx] = gl.history[hidx][:0]
	}
	for _, b := range gl.line {
		gl.history[hidx] = append(gl.history[hidx], b)
	}
	gl.historyCount++
	if gl.historyFile != nil {
		_, err := gl.historyFile.Write(gl.line)
		if err != nil {
			log.Printf("could not write history line: %v", err) // object allocated on the heap (OK)
		}
		_, err = gl.historyFile.Write([]byte{'\n'})
		if err != nil {
			log.Printf("could not write history cr: %v", err) // object allocated on the heap (OK)
		}
		gl.historyFile.Sync()
	}
}

func (gl *getLine) replaceLineWithHistory(historyIdx int, idx int) {
	oldlen := len(gl.line)
	newl := gl.history[(gl.historyCount-historyIdx)%MAX_HISTORY_LINES]
	// remove the existing line
	gl.txtb.Shift(-idx)
	for i := 0; i < oldlen; i++ {
		gl.txtb.Write(' ', true)
	}
	gl.txtb.Shift(-oldlen)
	gl.txtb.PrintBytes(newl, true)
	gl.line = gl.line[:0]
	for _, b := range newl {
		gl.line = append(gl.line, b)
	}
}

func (gl *getLine) pageDelta() int {
	return gl.txtb.Txtw.TextHeight() * 3 / 4
}

func (gl *getLine) enterScrollingMode(delta int) {
	gl.txtb.Cursor(false)
	scrollDelta, _ := gl.maybeScroll(0, delta)
	if scrollDelta == 0 {
		gl.drawScrollingBanner(true)
	}
	var exit bool
	for {
		c, err := gl.input.GetChar()
		if err != nil {
			return
		}
		switch c {
		case 'k':
			fallthrough
		case key.KEY_UP:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, -1)
		case 'j':
			fallthrough
		case key.KEY_DOWN:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, 1)
		case key.KEY_PAGEUP:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, -gl.pageDelta())
		case ' ':
			fallthrough
		case key.KEY_PAGEDOWN:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, gl.pageDelta())
		case 27:
			fallthrough
		case '\n':
			fallthrough
		case 'q':
			exit = true
		}
		if exit {
			gl.txtb.Scroll(-scrollDelta)
			gl.txtb.Update()
			gl.txtb.Cursor(true)
			gl.drawScrollingBanner(false)
			return
		}
	}
}

func (gl *getLine) maybeScroll(scrollDelta int, delta int) (int, bool) {
	newDelta := scrollDelta + delta
	maxDelta := gl.txtb.BufferLines() - gl.txtb.Txtw.TextHeight()
	if newDelta < -maxDelta {
		newDelta = -maxDelta
	}
	if newDelta > 0 {
		return 0, true
	}
	if newDelta != scrollDelta {
		gl.txtb.Scroll(newDelta - scrollDelta)
		gl.txtb.Update()
		gl.drawScrollingBanner(true)
	}
	return newDelta, false
}

// Here we draw directly on the text window (white text, blue background)
var scrollCol = window.White | (window.Blue >> 4)

const scrollMsg = "Scrolling Mode"

func (gl *getLine) drawScrollingBanner(enable bool) {
	w, h := gl.txtb.Txtw.TextSize()
	// Both DrawStr and RefreshArea can handle negative values
	x := (w - len(scrollMsg)) / 2
	if enable {
		window.DrawStr(gl.txtb.Txtw, x, h-1, scrollMsg, scrollCol)
	} else {
		gl.txtb.RefreshArea(x, h-1, len(scrollMsg), 1)
	}
}

func delete(line []byte, idx int) []byte {
	return append(line[:idx], line[idx+1:]...)
}
