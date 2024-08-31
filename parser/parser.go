package parser

import (
	"log"

	"github.com/alexhholmes/jg/lexer"
)

type Parser struct {
	lexer lexer.Lexer
}

func NewParser(input []byte) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(input),
	}
}

func (p *Parser) Parse() Document {
	for p.lexer.Next() != 0 {
		log.Printf("%#v\n", p.lexer.Tokens())
	}

	d := Document{}
	for p.lexer.Next() != 0 {
		t := p.lexer.Token()
		p.lexer.Next()
		switch t.Type {
		case lexer.LBRACE:
			d.elements = append(d.elements, CreateObject(&p.lexer))
		case lexer.LBRACKET:
			// TODO
		case lexer.EOF:
			return d
		default:
			log.Fatalf(
				"unexpected token at line %d column %d: must be object or array",
				t.Line,
				t.Pos,
			)
		}
	}

	return d
}
