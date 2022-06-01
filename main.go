package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unsafe-risk/protodecl/token"
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
			fmt.Println(ErrorPrint(err, string(file)))
			os.Exit(1)
		}
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	parser := NewParser(absfn, tokens)
	err = parser.Parse()
	if err != nil {
		fmt.Println(ErrorPrint(err, string(file)))
		os.Exit(1)
	}

	b, err := json.MarshalIndent(parser.Result(), "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
