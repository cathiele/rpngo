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
	for _, c := range m {
		switch parse.s {
		case WHITESPACE:
			parse.whitespace(c)
		case TOKEN:
			if err := parse.token(c, fn); err != nil {
				return err
			}
		case STRING_SINGLE:
			if err := parse.str(c, '\'', fn); err != nil {
				return err
			}
		case STRING_DOUBLE:
			if err := parse.str(c, '"', fn); err != nil {
				return err
			}
		case COMMENT:
			parse.comment(c)
		}
	}
	switch parse.s {
	case TOKEN:
		if err := parse.token('\n', fn); err != nil {
			return err
		}
	case STRING_SINGLE:
		return ErrUnterminatedSingleQuote
	case STRING_DOUBLE:
		return ErrUnterminatedDouble
	}
	return nil
}

func isWhitespace(c rune) bool {
	return (c == ' ') || (c == '\t') || (c == '\n')
}
