package lexer

import (
	"fmt"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input  []byte
	tokens []Token
	line   int
	// col is the current UTF-8 character on the current line of input.
	col int
	// pos is the current byte position in the input.
	pos int
}

func NewLexer(input []byte) Lexer {
	return Lexer{
		input: input,
		line:  1,
	}
}

// Reset sets points to the beginning of the input and clears the tokens.
func (l *Lexer) Reset() {
	l.tokens = []Token{}
	l.pos = 0
	l.col = 0
	l.line = 1
}

// Token returns the last token that was read.
func (l *Lexer) Token() Token {
	return l.tokens[len(l.tokens)-1]
}

// Tokens returns a copy of all tokens that have been read.
func (l *Lexer) Tokens() []Token {
	out := make([]Token, len(l.tokens))
	copy(out, l.tokens)
	return out
}

// Len returns the number of tokens that have been read.
func (l *Lexer) Len() int {
	return len(l.tokens)
}

func (l *Lexer) String() string {
	builder := strings.Builder{}
	for _, t := range l.tokens {
		builder.WriteString(fmt.Sprintf("%s: %s\n", t.Type, t.Literal))
	}
	return builder.String()
}

// Next reads the next token from the input and returns the number of bytes read.
func (l *Lexer) Next() (int, error) {
	if l.pos >= len(l.input) {
		l.line++
		l.col = 0
		l.appendToken(EOF, EOF.String())
		return 0, nil
	}

	// Read the next rune from the input
	r, size := utf8.DecodeRune(l.input[l.pos:])
	if r == utf8.RuneError {
		log.Fatalf("Invalid UTF-8 character at %d:%d", l.line, l.pos)
	} else if size == 0 {
		l.line++
		l.col = 0
		l.appendToken(EOF, EOF.String())
		return 0, nil
	}

	if token, ok := structures[r]; ok {
		// All the structure elements are one byte characters
		l.col++
		l.appendToken(token, token.String())
	} else {
		switch r {
		case '\n':
			l.line++
			l.col = 0
			l.appendToken(NEWLINE, NEWLINE.String())
		case ' ', '\t', '\r':
			if r == '\t' {
				l.col += 4
			} else {
				l.col++
			}
			l.appendToken(WHITESPACE, WHITESPACE.String())
		case 't', 'f', 'n':
			// Read true, false, and null literals
			if s, err := l.readReserved(); err != nil {
				return 0, err
			} else {
				size += s
			}
		case '"':
			// String literals; read until closing quote
			if s, err := l.readString(); err != nil {
				return 0, err
			} else {
				size += s
			}
		default:
			// Number literals
			if r == '-' || unicode.IsDigit(r) {
				if s, err := l.readNumber(size); err != nil {
					return 0, err
				} else {
					size += s
				}
			} else {
				return 0, fmt.Errorf("unexpected rune %s at %d:%d", string(r), l.line, l.col)
			}
		}
	}

	l.pos += size
	return size, nil
}

func (l *Lexer) readReserved() (int, error) {
	// Boolean and null literals; read until non-letter character
	slice := l.input[l.pos:]
	if len(slice) >= 5 && string(slice[:5]) == "false" {
		l.appendToken(FALSE, "false")
		l.col += 4
		return 4, nil
	} else if len(slice) >= 4 && string(slice[:4]) == "true" || string(slice[:4]) == "null" {
		l.appendToken(TRUE, string(slice[:4]))
		l.col += 3
		return 3, nil
	}

	return 0, fmt.Errorf("unexpected identifier at %d:%d", l.line, l.col)
}

func (l *Lexer) readString() (int, error) {
	// String literals; read until closing quote
	var size int
	for {
		r, s := utf8.DecodeRune(l.input[l.pos:])
		if r == utf8.RuneError {
			return 0, fmt.Errorf("invalid UTF-8 character at %d:%d", l.line, l.col)
		} else if structures[r] == EOF {
			return 0, fmt.Errorf("unexpected EOF at %d:%d", l.line, l.col)
		} else if r == '"' {
			// Read the closing quote
			if size != 0 {
				l.col--
				l.appendToken(STRING, string(l.input[l.pos:l.pos+size]))
				l.col++
			}
			l.appendToken(QUOTE, "\"")
			l.col++
			size += s

			return size, nil
		}

		l.col++
		size += s
	}
}

func (l *Lexer) readNumber(size int) (int, error) {
	hasDecimal := false
	for {
		if l.pos+size >= len(l.input) {
			return 0, fmt.Errorf("unexpected EOF at %d:%d", l.line, l.col)
		}

		r, s := utf8.DecodeRune(l.input[l.pos+size:])
		if r == utf8.RuneError {
			return 0, fmt.Errorf("invalid UTF-8 character at %d:%d", l.line, l.col)
		} else if t := structures[r]; t == WHITESPACE || t == NEWLINE || t == COMMA {
			break
		} else if r == '.' {
			if hasDecimal {
				log.Fatalf("Unexpected lexer %s at %d:%d", string(r), l.line, l.pos)
			}
			hasDecimal = true
		} else if !unicode.IsDigit(r) {
			log.Fatalf("Unexpected lexer %s at %d:%d", string(r), l.line, l.pos)
		}

		l.col++
		size += s
	}

	return size, nil
}

func (l *Lexer) appendToken(t TokenType, literal string) {
	l.tokens = append(
		l.tokens,
		Token{
			Type:    t,
			Literal: literal,
			Line:    l.line,
			Col:     l.pos,
		})
}

func (l *Lexer) peek() (rune, int) {
	r, size := utf8.DecodeRune(l.input[l.pos:])
	if l.pos+size >= len(l.input) {
		return 0, 0
	}

	r, size = utf8.DecodeRune(l.input[l.pos+size:])

	return r, size
}
