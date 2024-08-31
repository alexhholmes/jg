package main

import (
	"github.com/alexhholmes/jg/parser"
)

func main() {
	input := []byte(`{
		"key": "value",
		"num": 10000,
		"num2": 10000.0,
		"bool": true
	}`)

	parse := parser.NewParser(input)
	_ = parse.Parse()
}
