package main

import (
	"fmt"
	"strconv"

	"github.com/lemon-mint/protodecl/ast"
	"github.com/lemon-mint/protodecl/token"
)

type Parser struct {
	FileName string
	Tokens   []token.Token
	Out      ast.Tree

	Position int
}

func NewParser(filename string, data []token.Token) *Parser {
	return &Parser{
		FileName: filename,
		Tokens:   data,
		Position: 0,
	}
}

func (p *Parser) Result() ast.Tree {
	return p.Out
}

func (p *Parser) skipComments() {
	for p.Position < len(p.Tokens) && p.Tokens[p.Position].Type == token.Comment {
		p.Position++
	}
}

type ParserError struct {
	Tokens []token.Token
	Pos    int

	Message  string
	Position token.Position
}

func (p *ParserError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", p.Position.File, p.Position.Line, p.Position.Col, p.Message)
}

func newParserError(tkns []token.Token, pos int, msg string) *ParserError {
	return &ParserError{
		Tokens:   tkns,
		Pos:      pos,
		Message:  msg,
		Position: tkns[pos].Position,
	}
}

func (p *Parser) error(msg string) error {
	return newParserError(p.Tokens, p.Position, msg)
}

func (p *Parser) lenCheck() bool {
	return p.Position < len(p.Tokens) || p.Tokens[p.Position].Type == token.EOF
}

func (p *Parser) Parse() error {
	p.Out = ast.Tree{}
	p.Out.FileName = p.FileName
	p.Out.Nodes = p.Out.Nodes[:0]

	for p.Position < len(p.Tokens) {
		p.skipComments()
		switch p.Tokens[p.Position].Type {
		case token.Number:
			return p.error("unexpected numberLiteral " + p.Tokens[p.Position].Value)
		case token.Identifier:
			return p.error(fmt.Sprintf("unexpected identifier %s", p.Tokens[p.Position].Value))
		case token.Keyword:
			switch p.Tokens[p.Position].Value {
			case "enum", "packet", "protocol":
				n, err := p.parseType()
				if err != nil {
					return err
				}
				p.Out.Nodes = append(p.Out.Nodes, n)
			default:
				return p.error(fmt.Sprintf("unexpected keyword %s", p.Tokens[p.Position].Value))
			}
		}
	}

	return nil
}

func (p *Parser) parseNumber() (*ast.NumberLiteralType, error) {
	tkn := p.Tokens[p.Position]
	value, err := strconv.ParseUint(tkn.Value, 10, 64)
	if err != nil {
		return nil, err
	}
	p.Position++

	return &ast.NumberLiteralType{
		Position: tkn.Position,
		Value:    value,
	}, nil
}

func (p *Parser) parseEnum() (*ast.EnumerationType, error) {
	var err error
	tkn := p.Tokens[p.Position]
	if tkn.Type != token.Keyword || tkn.Value != "enum" {
		return nil, p.error(fmt.Sprintf("expected \"enum\" but got %s", tkn))
	}
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}
	if p.Tokens[p.Position].Type != token.Identifier {
		return nil, p.error(fmt.Sprintf("expected identifier but got %s", tkn))
	}
	name := p.Tokens[p.Position].Value
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}

	var rettype *ast.IdentifierType = new(ast.IdentifierType)
	switch p.Tokens[p.Position].Type {
	case token.Keyword:
		rettype.Value = p.Tokens[p.Position].Value
	case token.Identifier:
		rettype.Value = p.Tokens[p.Position].Value
	default:
		return nil, p.error(fmt.Sprintf("expected type but got %s", tkn))
	}
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}

	if p.Tokens[p.Position].Type != token.Delimiter || p.Tokens[p.Position].Value != "{" {
		return nil, p.error(fmt.Sprintf("expected '{' but got %s", tkn))
	}
	p.Position++

	var values []ast.EnumerationValue

L:
	for {
		p.skipComments()
		if !p.lenCheck() {
			return nil, p.error("unexpected EOF")
		}
		tkn = p.Tokens[p.Position]
		switch {
		case tkn.Type == token.Delimiter && tkn.Value == "}":
			p.Position++
			break L
		case tkn.Type == token.Identifier:
			v := ast.EnumerationValue{Key: tkn.Value}
			p.Position++
			p.skipComments()
			if !p.lenCheck() {
				return nil, p.error("unexpected EOF")
			}
			tkn = p.Tokens[p.Position]
			if tkn.Type != token.Delimiter || tkn.Value != "=" {
				return nil, p.error(fmt.Sprintf("expected '=' but got %s", tkn))
			}
			p.Position++
			p.skipComments()
			if !p.lenCheck() {
				return nil, p.error("unexpected EOF")
			}
			v.Value, err = p.parseNumber()
			if err != nil {
				return nil, err
			}
			values = append(values, v)
		default:
			return nil, p.error(fmt.Sprintf("expected identifier but got %s", tkn))
		}
	}

	return &ast.EnumerationType{
		Position:   tkn.Position,
		Name:       name,
		ReturnType: rettype,
		Values:     values,
	}, nil
}

func (p *Parser) parseType() (ast.Node, error) {
	p.skipComments()
	tkn := p.Tokens[p.Position]
	if tkn.Type != token.Keyword {
		return nil, p.error(fmt.Sprintf("expected keyword but got %s", tkn))
	}

	switch tkn.Value {
	case "enum":
		return p.parseEnum()
	default:
		return nil, p.error(fmt.Sprintf("unexpected keyword %s", tkn))
	}
	// panic("unreachable")
}