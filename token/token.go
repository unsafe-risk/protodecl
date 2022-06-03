package token

import (
	"fmt"
	"strconv"
)

type TType uint8

const (
	Identifier TType = iota
	Boolean
	Number
	Operator
	Delimiter
	Keyword
	Comment
	EOF
)

type Position struct {
	File string
	Line int
	Col  int
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}

type TokenType struct {
	Type  TType
	Value string
}

type Token struct {
	TokenType
	Position
}

func (t TType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case Boolean:
		return "Boolean"
	case Number:
		return "Number"
	case Operator:
		return "Operator"
	case Delimiter:
		return "Delimiter"
	case Keyword:
		return "Keyword"
	case Comment:
		return "Comment"
	case EOF:
		return "EOF"
	default:
		return "Unknown"
	}
}

func (t TokenType) String() string {
	return t.Type.String() + "<" + strconv.Quote(t.Value) + ">"
}

func (t Token) String() string {
	return t.TokenType.String()
}

func NewToken(t TType, v string, p Position) Token {
	return Token{
		TokenType{
			Type:  t,
			Value: v,
		},
		p,
	}
}
