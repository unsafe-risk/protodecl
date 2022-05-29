package main

import (
	"fmt"
	"io"
	"unicode"

	"github.com/lemon-mint/protodecl/token"
)

type Lexer struct {
	FileName string
	Data     []rune
	Line     int
	Col      int

	Position int
	Cursor   int

	CurrentChar rune
	LastToken   *token.Token
}

func NewLexer(filename string, data []rune) *Lexer {
	l := &Lexer{
		FileName:    filename,
		Data:        data,
		Line:        1,
		Col:         1,
		Position:    0,
		Cursor:      0,
		CurrentChar: '\n',
		LastToken:   nil,
	}
	l.readChar()
	return l
}

func (l *Lexer) newToken(t token.TokenType) token.Token {
	return token.Token{
		TokenType: t,
		Position: token.Position{
			File: l.FileName,
			Line: l.Line,
			Col:  l.Col,
		},
	}
}

func (l *Lexer) readChar() bool {
	if l.Cursor >= len(l.Data) {
		l.CurrentChar = '\n'
		l.LastToken = new(token.Token)
		*l.LastToken = l.newToken(token.TokenType{Type: token.EOF})
		return false
	}

	l.CurrentChar = l.Data[l.Cursor]
	if l.CurrentChar == '\n' {
		l.Line++
		l.Col = 0
	}
	l.Position = l.Cursor
	l.Cursor++
	l.Col++

	return true
}

func (l *Lexer) skipWhitespace() bool {
	for l.CurrentChar == ' ' || l.CurrentChar == '\t' || l.CurrentChar == '\n' || l.CurrentChar == '\r' {
		if !l.readChar() {
			return false
		}
	}

	return true
}

func (l *Lexer) readIdentifier() string {
	position := l.Position
	for unicode.IsLetter(l.CurrentChar) || unicode.IsDigit(l.CurrentChar) || l.CurrentChar == '_' {
		if !l.readChar() {
			break
		}
	}
	return string(l.Data[position:l.Position])
}

func (l *Lexer) nextChar() (c rune, ok bool) {
	if l.Cursor >= len(l.Data) {
		return '\n', false
	}
	return l.Data[l.Cursor], true
}

func (l *Lexer) NextToken() (t token.Token, err error) {
	if !l.skipWhitespace() {
		return token.Token{}, io.EOF
	}

	fmt.Printf("CurrentChar: %x\n", l.CurrentChar)
	switch l.CurrentChar {
	case '/':
		nextC, ok := l.nextChar()
		if !ok {
			return l.newToken(token.TokenType{Type: token.EOF}), io.EOF
		}
		switch nextC {
		case '/':
			if !l.readChar() || !l.readChar() {
				return l.newToken(token.TokenType{Type: token.EOF}), io.EOF
			}
			position := l.Position
			for l.CurrentChar != '\n' {
				if !l.readChar() {
					return l.newToken(token.TokenType{Type: token.EOF}), io.EOF
				}
			}
			commentStr := string(l.Data[position:l.Position])
			l.readChar()
			return l.newToken(token.TokenType{Type: token.Comment, Value: commentStr}), nil
		case '*':
			if !l.readChar() || !l.readChar() {
				return l.newToken(token.TokenType{Type: token.EOF}), io.EOF
			}
			position := l.Position
			for {
				if !l.readChar() {
					return l.newToken(token.TokenType{Type: token.EOF}), io.EOF
				}
				nextC, ok := l.nextChar()
				if !ok {
					return l.newToken(token.TokenType{Type: token.EOF}), io.EOF
				}
				if l.CurrentChar == '*' && nextC == '/' {
					break
				}
			}
			commentStr := string(l.Data[position:l.Position])
			l.readChar()
			l.readChar()
			return l.newToken(token.TokenType{Type: token.Comment, Value: commentStr}), nil
		default:
			t := l.newToken(token.TokenType{Type: token.Operator, Value: "/"})
			l.readChar()
			return t, nil
		}
	case '+', '-', '*', '%', '=', '<', '>', '!', '&', '|', '^', '~':
		t := l.newToken(token.TokenType{Type: token.Operator, Value: string(l.CurrentChar)})
		l.readChar()
		return t, nil
	case '{', '}', '(', ')', '[', ']', ';':
		t := l.newToken(token.TokenType{Type: token.Delimiter, Value: string(l.CurrentChar)})
		l.readChar()
		return t, nil

	default:
		id := l.readIdentifier()
		switch id {
		case "enum", "packet", "protocol", "message", "field",
			"bool", "u8", "u16", "u32", "u64", "u128", "i8", "i16", "i32", "i64", "i128",
			"CString", "String",
			"Cbytes", "Bytes",
			"Bytes8le", "Bytes16le", "Bytes32le", "Bytes64le",
			"Bytes8be", "Bytes16be", "Bytes32be", "Bytes64be",
			"String8le", "String16le", "String32le", "String64le",
			"String8be", "String16be", "String32be", "String64be",
			"Array", "Padding", "Bits",
			"f32", "f64",
			"true", "false":
			return l.newToken(token.TokenType{Type: token.Keyword, Value: id}), nil
		default:
			return l.newToken(token.TokenType{Type: token.Identifier, Value: id}), nil
		}
	}
	//panic("unreachable")
}
