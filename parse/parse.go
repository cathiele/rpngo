// Package parse turns a (potentially multiline) string into a set of tokens.
//
// Parse rules:
//
//  1. Spaces, tabs, or carriage returns can be used to separate tokens
//  2. Single or double quotes can be used to desgniate strings
//  3. A backslash can be used to cancel the meaning of proceeded character.
//     e.g.  'It\'s good' (or just say "It's good")
//
// 4 A # can be used for comments, which last until the end of the line
//
// Implementation is via a finite state machine
package parse

import (
	"errors"
	"fmt"
)

var (
	ErrUnterminatedDouble      = errors.New("unterminated double quote")
	ErrUnterminatedSingleQuote = errors.New("unterminated single quote")
)

type State uint8

const (
	WHITESPACE State = iota
	TOKEN
	STRING_DOUBLE
	STRING_SINGLE
	COMMENT
)

type parseData struct {
	s             State
	t             []rune
	nextIsLiteral bool
}

const defaultStaticRunes = 64

// Use a static instance to avoid heap allocations on every Fields call
var parse parseData = parseData{
	t: make([]rune, defaultStaticRunes),
}

func (p *parseData) init() {
	p.s = WHITESPACE
	if cap(p.t) > defaultStaticRunes {
		p.t = make([]rune, defaultStaticRunes) // object allocated on the heap: (OK)
	}
	p.t = p.t[:0]
}

func (p *parseData) whitespace(c rune) {
	if isWhitespace(c) {
		return
	}
	if c == '\\' {
		p.nextIsLiteral = true
		p.s = TOKEN
		return
	}
	if c == '#' {
		p.s = COMMENT
		return
	}
	if len(p.t) > 0 {
		p.t = p.t[:1]
		p.t[0] = c
	} else {
		p.t = append(p.t, c)
	}
	switch c {
	case '\'':
		p.s = STRING_SINGLE
	case '"':
		p.s = STRING_DOUBLE
	default:
		p.s = TOKEN
	}
}

func (p *parseData) token(c rune, fn func(string) error) error {
	if p.nextIsLiteral {
		p.t = append(p.t, c)
		p.nextIsLiteral = false
		return nil
	}
	if c == '\\' {
		p.nextIsLiteral = true
		return nil
	}
	if isWhitespace(c) {
		token := string(p.t)
		p.t = p.t[:0]
		p.s = WHITESPACE

		return fn(token)
	}
	p.t = append(p.t, c)
	return nil
}

func (p *parseData) str(c rune, quoteChar rune, fn func(string) error) error {
	if p.nextIsLiteral {
		switch c {
		case 'n':
			c = '\n'
		case 't':
			c = '\t'
		}
		p.t = append(p.t, c)
		p.nextIsLiteral = false
		return nil
	}
	if c == '\\' {
		p.nextIsLiteral = true
		return nil
	}
	p.t = append(p.t, c)
	if c == quoteChar {
		token := string(p.t)
		p.t = p.t[:0]
		p.s = WHITESPACE
		return fn(token)
	}
	return nil
}

func (p *parseData) comment(c rune) {
	if c == '\n' {
		p.s = WHITESPACE
	}
}

// Fields breaks a string into fields, calling fn for each parsed field
// (to avoid memory allocation)
func Fields(m string, fn func(string) error) error {
	parse.init()
	var err error
	starti := 0
	for i, c := range m {
		switch parse.s {
		case WHITESPACE:
			parse.whitespace(c)
			if parse.s != WHITESPACE {
				starti = i
			}
		case TOKEN:
			err = parse.token(c, fn)
		case STRING_SINGLE:
			err = parse.str(c, '\'', fn)
		case STRING_DOUBLE:
			err = parse.str(c, '"', fn)
		case COMMENT:
			parse.comment(c)
		}
		if err != nil {
			return fmt.Errorf(
				"%s: %w",
				buildContextString(m, starti, i),
				err)
		}
	}

	switch parse.s {
	case TOKEN:
		err = parse.token('\n', fn)
	case STRING_SINGLE:
		err = ErrUnterminatedSingleQuote
	case STRING_DOUBLE:
		err = ErrUnterminatedDouble
	}

	if err != nil {
		return fmt.Errorf(
			"%s: %w",
			buildContextString(m, starti, len(m)),
			err)
	}
	return nil
}

// set a max context length so that the error location is not buried
const maxContextLength = 80

func buildContextString(m string, starti, endi int) string {
	if len(m) > maxContextLength {
		m, starti, endi = truncateString(m, starti, endi, maxContextLength)
	}
	if starti >= len(m) {
		return m + "<-"
	}
	return m[:starti] + "->" + m[starti:endi] + "<-" + m[endi:]
}

func truncateString(m string, starti, endi, maxlen int) (string, int, int) {
	if endi - starti >= maxlen {
		return m[starti:starti+maxlen], 0, maxlen
	}

	span := endi - starti
	center := (starti + endi) / 2

	starts := center - span
	ends := starts + span

	if starts < 0 {
		// shift window right
		starti -= starts
		endi -= starts
		ends -= starts
		starts = 0
	}

	if ends >= len(m) {
		// shift window left
		delta := ends - len(m)
		starti -= delta
		endi -= delta
		starts -= delta
		ends -= delta
	}

	return m[starts:ends], starti, endi
}

func isWhitespace(c rune) bool {
	return (c == ' ') || (c == '\t') || (c == '\n')
}
