package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/field"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Enum struct {
	BaseDecl
}

func (*Enum) IsInterface() bool {
	return false
}
func (*Enum) IsAbstract() bool {
	return false
}
func (*Enum) IsClass() bool {
	return false
}
func (*Enum) IsConstant() bool {
	return false
}
func (*Enum) SetConstant() {}

func newEnum() *Enum {
	return &Enum{
		BaseDecl: newDecl(),
	}
}

func (e *Enum) String() string {
	return fmt.Sprintf("enum %s {\n%s\n%s\n%s\n%s}",
		e.Name().Value,
		strings.Join(maps.Keys(e.Methods_), "\n"),
		strings.Join(maps.Keys(e.Members_), "\n"),
		strings.Join(maps.Keys(e.StaticMembers_), "\n"),
		strings.Join(maps.Keys(e.StaticMembers_), "\n"),
	)
}

func (e *Enum) Syntax(p ast.SyntaxParser) (ast.Declaration, io.Error) {
	var err io.Error
	if e.BaseDecl.StartToken, err = p.Consume(token.TOK_ENUM); err != nil {
		return nil, err
	}

	e.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(token.TOK_LEFTBRACE); err != nil {
		return nil, err
	}

	if !p.Match(token.TOK_SEMICOLON) && p.Peek().Type != token.TOK_RIGHTBRACE {
		for {
			enumField := field.Enum{Type_: e}
			if err := enumField.Syntax(p); err != nil {
				return nil, err
			}

			e.StaticMembers_[enumField.Name().Value] = &enumField

			if p.Match(token.TOK_SEMICOLON) {
				break
			}

			if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return nil, err
			}
		}
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		f, err := p.ParseField()
		if err != nil {
			return nil, err
		} else if f.HasModifier(ast.MOD_VIRTUAL) {
			return nil, io.NewError("virtual fields are not allowed in enums", zap.Stringer("field", f.Name()))
		}
		if err := e.AddField(f); err != nil {
			return nil, err
		}
	}
	e.BaseDecl.EndToken = p.Prev()

	e.BaseDecl.Index, err = p.Types().NextEnumIndex(e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Enum) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], _ map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	if _, exists := visitedDecls.Get(e); exists {
		return nil, nil
	}
	visitedDecls.Set(e)

	return nil, e.BaseDecl.LinkParents(p, visitedDecls)
}

func (e *Enum) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(e); exists {
		return nil
	}
	defer visitedDecls.Set(e)

	return e.BaseDecl.LinkFields(p, visitedDecls)
}

func (e *Enum) Semantic(p ast.SemanticParser) io.Error {
	return e.BaseDecl.Semantic(p)
}

func (e *Enum) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return e.Equals(p, parent)
}

func (e *Enum) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return e.Equals(p, parent)
}

func (e *Enum) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherEnum, ok := other.(*Enum); ok {
		return e == otherEnum, nil
	} else if otherDeclType, ok := other.(ast.DeclType); ok {
		otherDeclaration, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		} else if otherEnum, ok := otherDeclaration.(*Enum); ok {
			return e == otherEnum, nil
		}
	}

	return false, nil
}

func (e *Enum) Concretize(mapping []ast.Type) ast.Type {
	return e
}
func (e *Enum) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	return typeid.NewTypeID(parser.Types().EnumIndex(e.Index), 0), nil
}
