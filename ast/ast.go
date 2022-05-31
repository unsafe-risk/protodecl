package ast

import "github.com/lemon-mint/protodecl/token"

type Tree struct {
	PackageName string
	FileName    string

	Root Node
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
	Type Node
	Name string
}

type PacketType struct {
	Position token.Position

	Name   string
	Fields []PacketField
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
