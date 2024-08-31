package lexer

type TokenType int

func (t TokenType) String() string {
	return []string{
		"EOF",
		"NEWLINE",
		"WHITESPACE",
		"QUOTE",
		"LBRACE",
		"RBRACE",
		"LBRACKET",
		"RBRACKET",
		"COMMA",
		"COLON",
		"NUMBER",
		"STRING",
		"TRUE",
		"FALSE",
		"NULL",
	}[t]
}

const (
	EOF TokenType = iota
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

// structures is a map of runes that correspond to JSON structure elements.
var structures = map[rune]TokenType{
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
	',': COMMA,
	':': COLON,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	// Col is the UTF-8 character position on the token's line.
	Col int
	// Pos is the byte position in the input where this token starts.
	Pos int
}
