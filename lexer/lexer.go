package lexer

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	EOF = iota
	NEWLINE
	WHITESPACE
	QUOTE
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	COMMA
	COLON
	NUMBER
	STRING
	TRUE
	FALSE
	NULL
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Pos     int
}

type Lexer struct {
	input    []byte
	tokens   []Token
	pos      int
	index    int
	line     int
	tokenMap map[string]TokenType
}

func NewLexer(input []byte) Lexer {
	return Lexer{
		input: input,
		line:  1,
		tokenMap: map[string]TokenType{
			"{":     LBRACE,
			"}":     RBRACE,
			"[":     LBRACKET,
			"]":     RBRACKET,
			",":     COMMA,
			":":     COLON,
			"true":  TRUE,
			"false": FALSE,
			"null":  NULL,
			"\"":    QUOTE,
			" ":     WHITESPACE,
			"\t":    WHITESPACE,
			"\r":    WHITESPACE,
			"\n":    NEWLINE,
		},
	}
}

func (l *Lexer) Reset() {
	l.pos = 0
	l.index = 0
	l.line = 1
}

func (l *Lexer) Token() Token {
	return l.tokens[l.index]
}

func (l *Lexer) Next() int {
	if l.input == nil || len(l.input) == 0 {
		l.appendToken(EOF, "")
		return 0
	}

	if l.pos >= len(l.input) {
		l.line++
		l.appendToken(EOF, "")
		return 0
	}

	r, size := utf8.DecodeRune(l.input[l.pos:])
	if r == utf8.RuneError {
		log.Fatalf("Invalid UTF-8 character at %d:%d", l.line, l.pos)
	}

	if token, ok := l.tokenMap[string(r)]; ok {
		size = l.lexToken(token, r, size)
	} else if r == 't' || r == 'f' || r == 'n' {
		// Boolean and null literals; read until non-letter character
		if string(l.input[l.pos:l.pos+5]) == "false" {
			l.appendToken(FALSE, "false")
			l.index += 4
			size += 4
		} else if string(l.input[l.pos:l.pos+4]) == "true" {
			l.appendToken(TRUE, "true")
			l.index += 3
			size += 3
		} else if string(l.input[l.pos:l.pos+4]) == "null" {
			l.appendToken(NULL, "null")
			l.index += 3
			size += 3
		}
	} else if unicode.IsDigit(r) || r == '-' {
		// Number literals; read until non-digit or decimal character
		hasDecimal := false
		for {
			r, s := utf8.DecodeRune(l.input[l.pos+size:])
			if r == utf8.RuneError {
				log.Fatalf("Invalid UTF-8 character at %d:%d", l.line, l.pos)
			} else if r == EOF {
				log.Fatalf("Unexpected EOF at %d:%d", l.line, l.pos)
			} else if t := l.tokenMap[string(r)]; t == WHITESPACE || t == NEWLINE || t == COMMA {
				break
			} else if r == '.' {
				if hasDecimal {
					log.Fatalf("Unexpected token %s at %d:%d", string(r), l.line, l.pos)
				}
				hasDecimal = true
			} else if !unicode.IsDigit(r) {
				log.Fatalf("Unexpected token %s at %d:%d", string(r), l.line, l.pos)
			}

			l.index++
			size += s
		}

		l.appendToken(NUMBER, string(l.input[l.pos:l.pos+size]))
	} else {
		log.Fatalf("Unexpected token %s at %d:%d", string(r), l.line, l.pos)
	}

	l.pos += size
	return size
}

func (l *Lexer) Tokens() []Token {
	return l.tokens
}

func (l *Lexer) String() string {
	builder := strings.Builder{}

	for _, t := range l.tokens {
		builder.WriteString(fmt.Sprintf("%s: %s\n", t.Type, t.Literal))
	}

	return builder.String()
}

func (l *Lexer) lexToken(token TokenType, r rune, size int) int {
	switch token {
	case EOF:
		l.appendToken(EOF, "")
		return size
	case NEWLINE:
		l.line++
		return size
	case WHITESPACE:
		// Ignore whitespace
		l.index++
		return size
	case QUOTE:
		// String literals; read until closing quote
		l.appendToken(QUOTE, string(r))
		pos := l.pos + size
		for {
			var s int
			r, s = utf8.DecodeRune(l.input[pos:])
			if r == utf8.RuneError {
				log.Fatalf("Invalid UTF-8 character at %d:%d", l.line, l.pos)
			} else if r == EOF {
				log.Fatalf("Unexpected EOF at %d:%d", l.line, l.pos)
			} else if r == '\\' {
				// TODO escape chars
			} else if r == '"' {
				literal := l.input[l.pos : l.pos+s]
				l.appendToken(STRING, string(l.input[l.pos+size:pos]))
				l.appendToken(QUOTE, string(literal))
				break
			}
			pos += s
		}
		size += pos - l.pos
	default:
		l.appendToken(token, string(r))
	}

	return size
}

func (l *Lexer) appendToken(t TokenType, literal string) {
	l.tokens = append(
		l.tokens,
		Token{
			Type:    t,
			Literal: literal,
			Line:    l.line,
			Pos:     l.pos,
		})
}

func (l *Lexer) peak() (rune, int) {
	r, size := utf8.DecodeRune(l.input[l.pos:])
	if l.pos+size >= len(l.input) {
		return 0, 0
	}

	r, size = utf8.DecodeRune(l.input[l.pos+size:])

	return r, size
}
