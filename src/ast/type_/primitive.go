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

func (p *Primitive[T]) Location() token.Location {
	return (*token.Token)(p).Location()
}

func (p *Primitive[T]) ExtendsAsPointer(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return p.Equals(ctx, parent)
}

func (p *Primitive[T]) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return p.Equals(ctx, parent)
}

func (p *Primitive[T]) Equals(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Primitive[T])
	return ok, nil
}

type Integer struct{ Primitive[Integer] }
type Float struct{ Primitive[Float] }
type Word struct{ Primitive[Word] }
type Character struct{ Primitive[Character] }
type String struct{ Primitive[String] }
type Boolean struct{ Primitive[Boolean] }
type Void struct{ Primitive[Void] }
type TypeID struct{ Primitive[TypeID] }
type Null struct{ Primitive[Null] }

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
		return nil, io.NewError("expected type", zap.Any("token", t))
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
