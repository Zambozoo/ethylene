package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

type Primitive[T any] token.Token

func (p *Primitive[T]) String() string {
	return (*token.Token)(p).String()
}

func (p *Primitive[T]) Key() string {
	return (*token.Token)(p).Value
}

func (p *Primitive[T]) Location() token.Location {
	return (*token.Token)(p).Location()
}

type Integer struct{ Primitive[Integer] }

func (p *Integer) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Integer) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Integer) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type Float struct{ Primitive[Float] }

func (p *Float) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Float) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Float) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type Word struct{ Primitive[Word] }

func (p *Word) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Word) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Word) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type Character struct{ Primitive[Character] }

func (p *Character) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Character) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Character) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type String struct{ Primitive[String] }

func (p *String) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *String) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *String) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type Boolean struct{ Primitive[Boolean] }

func (p *Boolean) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Boolean) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Boolean) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type Void struct{ Primitive[Void] }

func (p *Void) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Void) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Void) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type TypeID struct{ Primitive[TypeID] }

func (p *TypeID) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *TypeID) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *TypeID) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

type Null struct{ Primitive[Null] }

func (p *Null) ExtendsAsPointer(parent ast.Type) (bool, io.Error) { return p.Equals(parent) }
func (p *Null) Extends(parent ast.Type) (bool, io.Error)          { return p.Equals(parent) }
func (p *Null) Equals(other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}

func syntaxPrimitive(p ast.SyntaxParser) (ast.Type, io.Error) {
	switch t := p.Peek(); t.Type {
	case token.TOK_TYPEINT:
		return &Integer{Primitive: Primitive[Integer](p.Next())}, nil
	case token.TOK_TYPEFLT:
		return &Float{Primitive: Primitive[Float](p.Next())}, nil
	case token.TOK_TYPEWORD:
		return &Word{Primitive: Primitive[Word](p.Next())}, nil
	case token.TOK_TYPE:
		return &TypeID{Primitive: Primitive[TypeID](p.Next())}, nil
	case token.TOK_TYPESTR:
		return &String{Primitive: Primitive[String](p.Next())}, nil
	case token.TOK_TYPECHAR:
		return &Character{Primitive: Primitive[Character](p.Next())}, nil
	case token.TOK_TYPEBOOL:
		return &Boolean{Primitive: Primitive[Boolean](p.Next())}, nil
	case token.TOK_TYPEVOID:
		return &Void{Primitive: Primitive[Void](p.Next())}, nil
	default:
		return nil, io.NewError("expected type", zap.String("token", t.String()))
	}
}

func isCastablePrimitive(t ast.Type) bool {
	switch t.(type) {
	case *Integer, *Float, *Word, *Character:
		return true
	default:
		return false
	}
}

func CastPrimitive(p ast.SemanticParser, src, dst ast.Type) bool {
	if !isCastablePrimitive(src) || !isCastablePrimitive(dst) {
		return false
	}

	// TODO: bytecode
	switch src.(type) {
	case *Integer:
		switch dst.(type) {
		case *Float:
			// TODO: add cast bytecode
		case *Character:
			// TODO: add cast bytecode
		}
	case *Float:
		if _, toInt := dst.(*Integer); !toInt {
			return false
		}
		// TODO: add cast bytecode
	case *Word, *Character:
		if _, toInt := dst.(*Integer); !toInt {
			return false
		}
	}

	return true
}
