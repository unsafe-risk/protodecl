package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	for {
		token, err := lexer.NextToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(token)
	}
}
