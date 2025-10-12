package input

import (
	"mattwach/rpngo/rpn"
	"sort"
	"strings"
)

func (gl *getLine) tabComplete(r *rpn.RPN, idx int) int {
	if (idx <= 0) || (idx > len(gl.line)) {
		return idx
	}
	if (idx < len(gl.line)) && gl.line[idx] != ' ' {
		return idx
	}

	startIdx := gl.findStartOfWord(idx)
	if startIdx == idx {
		return idx
	}

	word := string(gl.line[startIdx:idx])
	//log.Printf("found word: %v", word)
	newWord := gl.findNewWord(r, word)

	if len(newWord) == 0 {
		return idx
	}

	//log.Printf("newword: %v", newWord)

	startLine := string(gl.line[:startIdx])
	endLine := string(gl.line[idx:])
	gl.line = gl.line[:0]
	for _, c := range startLine {
		gl.line = append(gl.line, byte(c))
	}
	for _, c := range newWord {
		gl.line = append(gl.line, byte(c))
	}
	for _, c := range endLine {
		gl.line = append(gl.line, byte(c))
	}

	// update the line
	gl.txtb.Shift(startIdx - idx)
	gl.txtb.PrintBytes(gl.line[startIdx:], true)
	numSpaces := len(word) - len(newWord)
	if numSpaces > 0 {
		for i := 0; i < numSpaces; i++ {
			gl.txtb.Write(' ', true)
		}
		gl.txtb.Shift(-numSpaces)
	}
	gl.txtb.Shift(-len(endLine))

	idx = idx + len(newWord) - len(word)
	return idx
}

func (gl *getLine) findStartOfWord(idx int) int {
	startIdx := idx
	for {
		if startIdx == 0 {
			break
		}
		lastChar := gl.line[startIdx-1]
		if (lastChar == ' ') || (lastChar == '\'') || (lastChar == '"') {
			break
		}
		startIdx--
	}
	return startIdx
}

func (gl *getLine) findNewWord(r *rpn.RPN, word string) string {
	var wordList []string
	var varPrefix string
	if word[0] == '$' {
		varPrefix = "$"
		word = word[1:]
		wordList = gl.allVariableNames(r)
	} else if word[0] == '@' {
		varPrefix = "@"
		word = word[1:]
		wordList = gl.allStringVariables(r)
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

func (gl *getLine) allVariableNames(r *rpn.RPN) []string {
	var wordList []string
	gl.namesAndValues = r.AppendAllVariableNamesAndValues(gl.namesAndValues[:0])
	for _, nv := range gl.namesAndValues {
		wordList = append(wordList, nv.Name)
	}
	sort.Strings(wordList)
	return wordList
}

func (gl *getLine) allStringVariables(r *rpn.RPN) []string {
	var wordList []string
	gl.namesAndValues = r.AppendAllVariableNamesAndValues(gl.namesAndValues[:0])
	for _, nv := range gl.namesAndValues {
		if nv.Values[len(nv.Values)-1].Type == rpn.STRING_FRAME {
			wordList = append(wordList, nv.Name)
		}
	}
	sort.Strings(wordList)
	return wordList
}
