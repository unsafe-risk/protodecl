package main

import (
	"fmt"
	"strconv"

	"github.com/unsafe-risk/protodecl/ast"
	"github.com/unsafe-risk/protodecl/token"
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
		/* // Preserve comments
		value := p.Tokens[p.Position].Value
		IsMultiline := strings.Contains(value, "\n")
		p.Out.Nodes = append(p.Out.Nodes, &ast.CommentType{
			Value:       value,
			Position:    p.Tokens[p.Position].Position,
			IsMultiline: IsMultiline,
		})
		*/
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
	p.Out.PackageName = p.FileName
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
		default:
			if p.Tokens[p.Position].Type == token.EOF {
				return nil
			}
			return p.error(fmt.Sprintf("unexpected token %s", p.Tokens[p.Position]))
		}
	}

	return nil
}

func (p *Parser) parseNumber() (*ast.NumberLiteralType, error) {
	tkn := p.Tokens[p.Position]
	value, err := strconv.ParseUint(tkn.Value, 10, 64)
	if err != nil {
		return nil, p.error(err.Error())
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

	rettype, err := p.parseType()
	if err != nil {
		return nil, err
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
			if tkn.Type != token.Operator || tkn.Value != "=" {
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
			p.skipComments()
			if !p.lenCheck() {
				return nil, p.error("unexpected EOF")
			}
			tkn = p.Tokens[p.Position]
			if tkn.Type != token.Delimiter || tkn.Value != ";" {
				return nil, p.error(fmt.Sprintf("expected ';' but got %s", tkn))
			}
			p.Position++
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

func (p *Parser) parsePacket() (*ast.PacketType, error) {
	var err error
	tkn := p.Tokens[p.Position]
	if tkn.Type != token.Keyword || tkn.Value != "packet" {
		return nil, p.error(fmt.Sprintf("expected \"packet\" but got %s", tkn))
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

	if p.Tokens[p.Position].Type != token.Delimiter || p.Tokens[p.Position].Value != "(" {
		return nil, p.error(fmt.Sprintf("expected '(' but got %s", tkn))
	}
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}
	var args []ast.PacketField

	for {
		tkn = p.Tokens[p.Position]
		if tkn.Type == token.Delimiter && tkn.Value == ")" {
			p.Position++
			break
		}
		if tkn.Type != token.Identifier {
			return nil, p.error(fmt.Sprintf("expected identifier but got %s", tkn))
		}
		arg := ast.PacketField{
			Name: tkn.Value,
		}
		p.Position++
		p.skipComments()
		if !p.lenCheck() {
			return nil, p.error("unexpected EOF")
		}

		tkn = p.Tokens[p.Position]
		if tkn.Type != token.Delimiter || tkn.Value != ":" {
			return nil, p.error(fmt.Sprintf("expected ':' but got %s", tkn))
		}

		p.Position++
		p.skipComments()
		if !p.lenCheck() {
			return nil, p.error("expected <Type> but got EOF")
		}
		arg.Type, err = p.parseType()
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}

	if p.Tokens[p.Position].Type != token.Delimiter || p.Tokens[p.Position].Value != "{" {
		return nil, p.error(fmt.Sprintf("expected '{' but got %s", tkn))
	}
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}

	var fields []ast.PacketField
	for {
		tkn = p.Tokens[p.Position]
		if tkn.Type == token.Delimiter && tkn.Value == "}" {
			p.Position++
			break
		}

		t, err := p.parseType()
		if err != nil {
			return nil, err
		}

		p.skipComments()
		if !p.lenCheck() {
			return nil, p.error("unexpected EOF")
		}

		if p.Tokens[p.Position].Type != token.Identifier && p.Tokens[p.Position].Type != token.Keyword {
			return nil, p.error(fmt.Sprintf("expected identifier but got %s", tkn))
		}
		name := p.Tokens[p.Position].Value
		p.Position++
		p.skipComments()
		if !p.lenCheck() {
			return nil, p.error("unexpected EOF")
		}

		if p.Tokens[p.Position].Type != token.Delimiter || p.Tokens[p.Position].Value != ";" {
			return nil, p.error(fmt.Sprintf("expected ';' but got %s", tkn))
		}
		p.Position++
		p.skipComments()

		fields = append(fields, ast.PacketField{
			Name: name,
			Type: t,
		})
	}

	return &ast.PacketType{
		Position:   tkn.Position,
		Name:       name,
		Parameters: args,
		Fields:     fields,
	}, nil
}

func (p *Parser) parseProtocol() (*ast.ProtocolType, error) {
	tkn := p.Tokens[p.Position]
	if tkn.Type != token.Keyword || tkn.Value != "protocol" {
		return nil, p.error(fmt.Sprintf("expected \"protocol\" but got %s", tkn))
	}
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}

	if p.Tokens[p.Position].Type != token.Identifier {
		return nil, p.error(fmt.Sprintf("expected identifier but got %s", tkn))
	}
	return nil, nil
}

func (p *Parser) parseType() (ast.Node, error) {
	p.skipComments()
	tkn := p.Tokens[p.Position]

	if tkn.Type != token.Keyword && tkn.Type != token.Identifier {
		return nil, p.error(fmt.Sprintf("expected type but got %s", tkn))
	}

	switch tkn.Value {
	case "enum":
		return p.parseEnum()
	case "packet":
		return p.parsePacket()
	case "protocol":
		return p.parseProtocol()
	}

	name := tkn.Value
	p.Position++
	p.skipComments()
	if !p.lenCheck() {
		return nil, p.error("unexpected EOF")
	}
	tkn = p.Tokens[p.Position]

	var args []ast.Node
	if tkn.Type == token.Delimiter && tkn.Value == "(" {
		p.Position++
		p.skipComments()
		if !p.lenCheck() {
			return nil, p.error("unexpected EOF")
		}
		for {
			tkn = p.Tokens[p.Position]
			if tkn.Type == token.Delimiter && tkn.Value == ")" {
				p.Position++
				break
			}

			if tkn.Type == token.Delimiter && tkn.Value == "," {
				p.Position++
				p.skipComments()
				if !p.lenCheck() {
					return nil, p.error("unexpected EOF")
				}
				tkn = p.Tokens[p.Position]
			}

			switch tkn.Type {
			case token.Identifier:
				args = append(args, &ast.IdentifierType{
					Position: tkn.Position,
					Value:    tkn.Value,
				})
				p.Position++
				p.skipComments()
				if !p.lenCheck() {
					return nil, p.error("unexpected EOF")
				}
			case token.Number:
				n, err := p.parseNumber()
				if err != nil {
					return nil, err
				}
				args = append(args, n)
			default:
				return nil, p.error(fmt.Sprintf("expected identifier or number but got %s", tkn))
			}
		}
	}
	return &ast.TypeType{
		Position:  tkn.Position,
		TypeName:  name,
		Arguments: args,
	}, nil
}
