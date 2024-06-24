package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/field"
	"geth-cody/ast/type_"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Enum struct {
	BaseDecl
}

func newEnum() *Enum {
	return &Enum{
		BaseDecl: newDecl(),
	}
}

func (e *Enum) Generics() map[string]ast.GenericConstraint {
	return map[string]ast.GenericConstraint{}
}

func (e *Enum) String() string {
	return fmt.Sprintf("Interface{Name: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		e.Name().Value,
		strings.Join(maps.Keys(e.Methods_), ","),
		strings.Join(maps.Keys(e.Members_), ","),
		strings.Join(maps.Keys(e.StaticMembers_), ","),
		strings.Join(maps.Keys(e.StaticMembers_), ","),
	)
}

func (e *Enum) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if e.BaseDecl.StartToken, err = p.Consume(token.TOK_ENUM); err != nil {
		return err
	}

	e.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	if _, err := p.Consume(token.TOK_LEFTBRACE); err != nil {
		return err
	}

	if !p.Match(token.TOK_SEMICOLON) && p.Peek().Type != token.TOK_RIGHTBRACE {
		for {
			enumField := field.Enum{
				Type_: &type_.Composite{
					Context_: p.TypeContext(),
					Tokens:   []token.Token{e.Name_},
				},
			}

			if err := enumField.Syntax(p); err != nil {
				return err
			}

			e.StaticMembers_[enumField.Name().Value] = &enumField

			if p.Match(token.TOK_SEMICOLON) {
				break
			}

			if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return err
			}
		}
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		f, err := p.ParseField()
		if err != nil {
			return err
		} else if f.HasModifier(ast.MOD_VIRTUAL) {
			return io.NewError("virtual fields are not allowed in enums", zap.Any("field", f.Name()))
		}
		if err := e.AddField(f); err != nil {
			return err
		}
	}
	e.BaseDecl.EndToken = p.Prev()

	return nil
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
