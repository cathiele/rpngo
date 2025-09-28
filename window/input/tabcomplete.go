package input

import (
	"mattwach/rpngo/rpn"
	"mattwach/rpngo/window"
	"sort"
	"strings"
)

func (gl *getLine) tabComplete(r *rpn.RPN, line []byte, idx int) ([]byte, int) {
	if (idx <= 0) || (idx > len(line)) {
		return line, idx
	}
	if (idx < len(line)) && line[idx] != ' ' {
		return line, idx
	}

	startIdx := findStartOfWord(line, idx)
	if startIdx == idx {
		return line, idx
	}

	word := string(line[startIdx:idx])
	//log.Printf("found word: %v", word)
	newWord := findNewWord(r, word)

	if len(newWord) == 0 {
		return line, idx
	}

	//log.Printf("newword: %v", newWord)

	startLine := string(line[:startIdx])
	endLine := string(line[idx:])
	line = []byte(startLine + newWord + endLine)

	// update the line
	window.Shift(gl.txtd, startIdx-idx)
	window.PrintBytes(gl.txtd, line[startIdx:])
	numSpaces := len(word) - len(newWord)
	if numSpaces > 0 {
		for i := 0; i < numSpaces; i++ {
			window.PutByte(gl.txtd, ' ')
		}
		window.Shift(gl.txtd, -numSpaces)
	}
	window.Shift(gl.txtd, -len(endLine))

	idx = idx + len(newWord) - len(word)
	//log.Printf("idx=%v line=%v", idx, string(line))
	return line, idx
}

func findStartOfWord(line []byte, idx int) int {
	startIdx := idx
	for {
		if startIdx == 0 {
			break
		}
		lastChar := line[startIdx-1]
		if (lastChar == ' ') || (lastChar == '\'') || (lastChar == '"') {
			break
		}
		startIdx--
	}
	return startIdx
}

func findNewWord(r *rpn.RPN, word string) string {
	var wordList []string
	var varPrefix string
	if word[0] == '$' {
		varPrefix = "$"
		word = word[1:]
		wordList = allVariableNames(r)
	} else if word[0] == '@' {
		varPrefix = "@"
		word = word[1:]
		wordList = allStringVariables(r)
	} else {
		wordList = r.AllFunctionNames()
	}

	// Look for an exact match of the word
	var newWord string
	for wordIdx := 0; wordIdx < len(wordList); wordIdx++ {
		if wordList[wordIdx] == word {
			newWord = wordList[(wordIdx+1)%len(wordList)]
			break
		}
	}

	if len(newWord) == 0 {
		// look for a partial match
		for wordIdx := 0; wordIdx < len(wordList); wordIdx++ {
			if strings.HasPrefix(wordList[wordIdx], word) {
				newWord = wordList[wordIdx]
				break
			}
		}
	}

	return varPrefix + newWord
}

func allVariableNames(r *rpn.RPN) []string {
	var wordList []string
	for _, nv := range r.AllVariableNamesAndValues() {
		wordList = append(wordList, nv.Name)
	}
	sort.Strings(wordList)
	return wordList
}

func allStringVariables(r *rpn.RPN) []string {
	var wordList []string
	for _, nv := range r.AllVariableNamesAndValues() {
		if nv.Values[len(nv.Values)-1].Type == rpn.STRING_FRAME {
			wordList = append(wordList, nv.Name)
		}
	}
	sort.Strings(wordList)
	return wordList
}
