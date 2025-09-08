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
	ret           []string
	nextIsLiteral bool
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
	p.t = p.t[:1]
	p.t[0] = c
	switch c {
	case '\'':
		p.s = STRING_SINGLE
	case '"':
		p.s = STRING_DOUBLE
	default:
		p.s = TOKEN
	}
	return
}

func (p *parseData) token(c rune) {
	if p.nextIsLiteral {
		p.t = append(p.t, c)
		p.nextIsLiteral = false
		return
	}
	if c == '\\' {
		p.nextIsLiteral = true
		return
	}
	if isWhitespace(c) {
		p.ret = append(p.ret, string(p.t))
		p.t = p.t[:0]
		p.s = WHITESPACE
		return
	}
	p.t = append(p.t, c)
	return
}

func (p *parseData) str(c rune, quoteChar rune) {
	if p.nextIsLiteral {
		p.t = append(p.t, c)
		p.nextIsLiteral = false
		return
	}
	if c == '\\' {
		p.nextIsLiteral = true
		return
	}
	p.t = append(p.t, c)
	if c == quoteChar {
		p.ret = append(p.ret, string(p.t))
		p.t = p.t[:0]
		p.s = WHITESPACE
	}
}

func (p *parseData) comment(c rune) {
	if c == '\n' {
		p.s = WHITESPACE
	}
}

func Fields(m string) ([]string, error) {
	var p parseData = parseData{t: make([]rune, 'x')}
	for _, c := range m {
		switch p.s {
		case WHITESPACE:
			p.whitespace(c)
		case TOKEN:
			p.token(c)
		case STRING_SINGLE:
			p.str(c, '\'')
		case STRING_DOUBLE:
			p.str(c, '"')
		case COMMENT:
			p.comment(c)
		}
	}
	if p.s == TOKEN {
		p.token('\n')
	} else if p.s == STRING_SINGLE {
		return nil, ErrUnterminatedSingleQuote
	} else if p.s == STRING_DOUBLE {
		return nil, ErrUnterminatedDouble
	}
	return p.ret, nil
}

func isWhitespace(c rune) bool {
	return (c == ' ') || (c == '\t') || (c == '\n')
}
