package main

import (
	"fmt"
	"strconv"
	"strings"
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
			Col:  l.Col - 1,
		},
	}
}

func (l *Lexer) readChar() bool {
	if l.Cursor >= len(l.Data) {
		l.CurrentChar = 0
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
	//fmt.Printf("%c at %d:%d\n", l.CurrentChar, l.Line, l.Col)
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

type LexerError struct {
	Message  string
	Filename string
	Line     int
	Col      int
	Index    int
}

func (e LexerError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.Filename, e.Line, e.Col, e.Message)
}

func (l *Lexer) dumpError(msg string) *LexerError {
	return &LexerError{
		Message:  msg,
		Filename: l.FileName,
		Line:     l.Line,
		Col:      l.Col,
		Index:    l.Position,
	}
}

func (l *Lexer) NextToken() (t token.Token, err error) {
	if !l.skipWhitespace() {
		return l.newToken(token.TokenType{Type: token.EOF}), nil
	}

	//fmt.Printf("CurrentChar: %x\n", l.CurrentChar)
	switch l.CurrentChar {
	case '/':
		line, col := l.Line, l.Col
		nextC, ok := l.nextChar()
		if !ok {
			return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("unexpected EOF")
		}
		switch nextC {
		case '/':
			if !l.readChar() {
				return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("Expected '\\n' but got EOF")
			}
			if !l.readChar() {
				return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("Expected '\\n' but got EOF")
			}
			position := l.Position
			for l.CurrentChar != '\n' {
				if !l.readChar() {
					break
				}
			}
			commentStr := string(l.Data[position:l.Position])
			t := l.newToken(token.TokenType{Type: token.Comment, Value: commentStr})
			t.Position.Line, t.Position.Col = line, col
			return t, nil
		case '*':
			if !l.readChar() {
				return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("Expected END_OF_COMMENT but got EOF")
			}
			if !l.readChar() {
				return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("Expected END_OF_COMMENT but got EOF")
			}
			position := l.Position
			for {
				if !l.readChar() {
					return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("Expected END_OF_COMMENT but got EOF")
				}
				nextC, ok := l.nextChar()
				if !ok {
					return l.newToken(token.TokenType{Type: token.EOF}), l.dumpError("Expected END_OF_COMMENT but got EOF")
				}
				if l.CurrentChar == '*' && nextC == '/' {
					break
				}
			}
			commentStr := string(l.Data[position:l.Position])
			l.readChar()
			l.readChar()
			t := l.newToken(token.TokenType{Type: token.Comment, Value: commentStr})
			t.Position.Line, t.Position.Col = line, col
			return t, nil
		default:
			t := l.newToken(token.TokenType{Type: token.Operator, Value: "/"})
			l.readChar()
			return t, nil
		}
	case '+', '-', '*', '%', '=', '<', '>', '!', '&', '|', '^', '~':
		t := l.newToken(token.TokenType{Type: token.Operator, Value: string(l.CurrentChar)})
		l.readChar()
		return t, nil
	case '{', '}', '(', ')', '[', ']', ';', ':', '.', ',':
		t := l.newToken(token.TokenType{Type: token.Delimiter, Value: string(l.CurrentChar)})
		l.readChar()
		return t, nil

	default:
		line, col := l.Line, l.Col-1
		id := l.readIdentifier()
		switch id {
		case "enum", "packet", "protocol", "message", "field",
			"bool", "u8", "u16", "u32", "u64", "u128", "i8", "i16", "i32", "i64", "i128",
			"CString", "String",
			"Cbytes", "Bytes",
			"LongString", "LongBytes",
			"Bytes8le", "Bytes16le", "Bytes32le", "Bytes64le",
			"Bytes8be", "Bytes16be", "Bytes32be", "Bytes64be",
			"String8le", "String16le", "String32le", "String64le",
			"String8be", "String16be", "String32be", "String64be",
			"Array", "Padding", "Bits",
			"f32", "f64",
			"true", "false":
			t := l.newToken(token.TokenType{Type: token.Keyword, Value: id})
			t.Line, t.Col = line, col
			return t, nil
		default:
			// t := l.newToken(token.TokenType{Type: token.Identifier, Value: id})
			// t.Line, t.Col = line, col

			if len(id) <= 0 {
				t := l.newToken(token.TokenType{Type: token.Identifier, Value: id})
				t.Line, t.Col = line, col
			}

			// Parse number

			if id[0] >= '0' && id[0] <= '9' {
				if strings.HasPrefix(id, "0x") {
					// parse hex number
					num, err := strconv.ParseUint(id[2:], 16, 64)
					if err != nil {
						return l.newToken(token.TokenType{Type: token.Number}), l.dumpError("invalid hex number (Error: " + strconv.Quote(err.Error()) + ")")
					}
					t := l.newToken(token.TokenType{Type: token.Number, Value: id[2:]})
					t.Line, t.Col = line, col
					t.Value = strconv.FormatUint(num, 10)
					return t, nil
				} else if strings.HasPrefix(id, "0b") {
					// parse binary number
					num, err := strconv.ParseUint(id[2:], 2, 64)
					if err != nil {
						return l.newToken(token.TokenType{Type: token.Number}), l.dumpError("invalid binary number (Error: " + strconv.Quote(err.Error()) + ")")
					}
					t := l.newToken(token.TokenType{Type: token.Number, Value: id[2:]})
					t.Line, t.Col = line, col
					t.Value = strconv.FormatUint(num, 10)
					return t, nil
				} else {
					// parse decimal number
					num, err := strconv.ParseUint(id, 10, 64)
					if err != nil {
						return l.newToken(token.TokenType{Type: token.Number}), l.dumpError("invalid decimal number (Error: " + strconv.Quote(err.Error()) + ")")
					}
					t := l.newToken(token.TokenType{Type: token.Number, Value: id})
					t.Line, t.Col = line, col
					t.Value = strconv.FormatUint(num, 10)
					return t, nil
				}
			}

			// Parse string

			t := l.newToken(token.TokenType{Type: token.Identifier, Value: id})
			t.Line, t.Col = line, col
			return t, nil
		}
	}
	//panic("unreachable")
}
