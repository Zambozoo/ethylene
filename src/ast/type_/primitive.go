package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"

	"go.uber.org/zap"
)

type Primitive[T any] struct {
	Constant bool
	token.Token
}

func (p *Primitive[T]) IsConstant() bool {
	return p.Constant
}
func (p *Primitive[T]) SetConstant() {
	p.Constant = true
}

func (p *Primitive[T]) String() string {
	return p.Token.Value
}

func (p *Primitive[T]) Location() *token.Location {
	return p.Token.Location()
}

func (p *Primitive[T]) IsConcrete() bool {
	return true
}

type Integer struct{ Primitive[Integer] }

func (p *Integer) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Integer) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Integer) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Integer)
	return ok, nil
}
func (p *Integer) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Integer) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Int.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type Float struct{ Primitive[Float] }

func (p *Float) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Float) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Float) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Float)
	return ok, nil
}
func (p *Float) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Float) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Float.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type Word struct{ Primitive[Word] }

func (p *Word) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Word) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Word) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Word)
	return ok, nil
}
func (p *Word) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Word) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Word.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type Character struct{ Primitive[Character] }

func (p *Character) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Character) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Character) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Character)
	return ok, nil
}
func (p *Character) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Character) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Char.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type String struct{ Primitive[String] }

func (p *String) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *String) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *String) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*String)
	return ok, nil
}
func (p *String) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *String) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Str.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type Boolean struct{ Primitive[Boolean] }

func (p *Boolean) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Boolean) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Boolean) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Boolean)
	return ok, nil
}
func (p *Boolean) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Boolean) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Bool.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type Void struct{ Primitive[Void] }

func (p *Void) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Void) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Void) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Void)
	return ok, nil
}
func (p *Void) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Void) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_Void.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type TypeID struct{ Primitive[TypeID] }

func (p *TypeID) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *TypeID) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *TypeID) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*TypeID)
	return ok, nil
}
func (p *TypeID) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *TypeID) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	index := typeid.ID_TypeID.Index()
	if p.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, 0), nil
}

type Null struct{ Primitive[Null] }

func (p *Null) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Null) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	_, ok := parent.(*Pointer)
	return ok, nil
}
func (p *Null) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	_, ok := other.(*Null)
	return ok, nil
}
func (p *Null) Concretize([]ast.Type) ast.Type {
	return p
}
func (p *Null) TypeID(ast.SemanticParser) (ast.TypeID, io.Error) {
	return nil, io.NewError("null does not have a type id", zap.Stringer("location", p.Location()))
}

func syntaxPrimitive(p ast.SyntaxParser) (ast.Type, io.Error) {
	switch t := p.Peek(); t.Type {
	case token.TOK_TYPEINT:
		return &Integer{Primitive: Primitive[Integer]{Token: p.Next()}}, nil
	case token.TOK_TYPEFLT:
		return &Float{Primitive: Primitive[Float]{Token: p.Next()}}, nil
	case token.TOK_TYPEWORD:
		return &Word{Primitive: Primitive[Word]{Token: p.Next()}}, nil
	case token.TOK_TYPE:
		return &TypeID{Primitive: Primitive[TypeID]{Token: p.Next()}}, nil
	case token.TOK_TYPESTR:
		return &String{Primitive: Primitive[String]{Token: p.Next()}}, nil
	case token.TOK_TYPECHAR:
		return &Character{Primitive: Primitive[Character]{Token: p.Next()}}, nil
	case token.TOK_TYPEBOOL:
		return &Boolean{Primitive: Primitive[Boolean]{Token: p.Next()}}, nil
	case token.TOK_TYPEVOID:
		return &Void{Primitive: Primitive[Void]{Token: p.Next()}}, nil
	default:
		return nil, io.NewError("expected type", zap.Stringer("token", &t))
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
