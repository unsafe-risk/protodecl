package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("example.protodecl")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lexer := NewLexer("example.protodecl", []rune(string(file)))
	for {
		token, err := lexer.NextToken()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(token)
	}
}
