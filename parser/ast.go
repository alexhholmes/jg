package parser

import (
	"log"

	"github.com/alexhholmes/jg/lexer"
)

type Document struct {
	// index is a multi-level map of elements in this Document for quicker
	// lookups of leaf elements.
	index    map[string]element
	elements []element
}

type element interface {
	AddElement(name string, e element)
}

type object struct {
	parent element
	m      map[string]element
}

func CreateObject(l *lexer.Lexer) *object {
	for l.Token().Type != lexer.RBRACE {
		l.Next()
		if l.Token().Type == lexer.EOF {
			log.Fatalf("unexpected EOF at line %d column %d", l.Token().Line, l.Token().Col)
		} else if l.Token().Type == lexer.STRING {
			_ = l.Token().Literal // TODO
			l.Next()
			if l.Token().Type != lexer.COLON {
				log.Fatalf("unexpected lexer at line %d column %d: expected colon", l.Token().Line, l.Token().Col)
			}
			l.Next()
		}

	}

	return &object{m: make(map[string]element)}
}

func (o *object) AddElement(name string, e element) {
	if _, ok := o.m[name]; ok {
		log.Fatalf("element with name %s already exists", name)
	}
	o.m[name] = e
}

type array struct {
	parent element
	e      []element
}

func (a *array) AddElement(_ string, e element) {
	a.e = append(a.e, e)
}

type value struct {
	parent element
	t      lexer.Token
}

func (v *value) AddElement(name string, e element) {
	v.parent.AddElement(name, e)
}
