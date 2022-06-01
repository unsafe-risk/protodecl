package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kr/pretty"
	"github.com/lemon-mint/protodecl/token"
)

func main() {
	absfn, err := filepath.Abs("example.protodecl")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.ReadFile(absfn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lexer := NewLexer(absfn, []rune(string(file)))
	var tokens []token.Token
	for {
		tok, err := lexer.NextToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(tok)
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	parser := NewParser(absfn, tokens)
	err = parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pretty.Println(parser.Result())
}
