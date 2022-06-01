package ast

import "github.com/unsafe-risk/protodecl/token"

type Tree struct {
	PackageName string
	FileName    string

	Nodes []Node
}

type Node interface {
	Pos() token.Position
}

type EnumerationValue struct {
	Key   string
	Value Node
}

type EnumerationType struct {
	Position token.Position

	Name       string
	ReturnType Node

	Values []EnumerationValue
}

func (e *EnumerationType) Pos() token.Position {
	return e.Position
}

type PacketField struct {
	Name string
	Type Node
}

type PacketType struct {
	Position token.Position

	Name string

	Parameters []PacketField
	Fields     []PacketField
}

func (p *PacketType) Pos() token.Position {
	return p.Position
}

type ProtocolType struct {
	Position token.Position

	Name string
}

func (p *ProtocolType) Pos() token.Position {
	return p.Position
}

type NumberLiteralType struct {
	Position token.Position

	Value uint64
}

func (n *NumberLiteralType) Pos() token.Position {
	return n.Position
}

type IdentifierType struct {
	Position token.Position

	Value string
}

func (i *IdentifierType) Pos() token.Position {
	return i.Position
}

type CommentType struct {
	Position token.Position

	IsMultiline bool
	Value       string
}

func (c *CommentType) Pos() token.Position {
	return c.Position
}

type TypeType struct {
	Position token.Position

	TypeName  string
	Arguments []Node
}

func (t *TypeType) Pos() token.Position {
	return t.Position
}
