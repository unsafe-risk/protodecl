package main

import (
	"log"

	"github.com/unsafe-risk/protodecl/compile"
	"github.com/unsafe-risk/protodecl/parser"
)

func main() {
	ast, file, err := parser.ParseFile("example.protodecl")
	if err != nil {
		log.Fatalln(parser.ErrorPrint(err, string(file)))
	}
	compile.Compile(ast)
}
