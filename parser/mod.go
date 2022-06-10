package parser

import (
	"os"
	"path/filepath"

	"github.com/unsafe-risk/protodecl/ast"
	"github.com/unsafe-risk/protodecl/token"
)

func ParseFile(filename string) (ast *ast.Tree, file []byte, err error) {
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		return nil, nil, err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	relFilename, err := filepath.Rel(currentDir, absFilename)
	if err != nil {
		return nil, nil, err
	}

	file, err = os.ReadFile(relFilename)
	if err != nil {
		return
	}

	lexer := NewLexer(relFilename, []rune(string(file)))
	var tokens []token.Token
	for {
		tok, err := lexer.NextToken()
		if err != nil {
			return nil, file, err
		}
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	pp := NewParser(relFilename, tokens)
	err = pp.Parse()
	if err != nil {
		return nil, file, err
	}

	tree := pp.Result()

	return &tree, file, nil
}
