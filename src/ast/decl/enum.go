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

func (*Enum) GenericsMap() map[string]ast.DeclType {
	return map[string]ast.DeclType{}
}

func (*Enum) Generics() []ast.DeclType {
	return nil
}

func (*Enum) GenericsCount() int {
	return 0
}

func (e *Enum) PutGeneric(name string, generic ast.DeclType) io.Error {
	return io.NewError("enums cannot have generic type parameters", zap.Any("location", e.Name_.Location()))
}

func (e *Enum) SetTailed() io.Error {
	return io.NewError("enums cannot be tailed", zap.Any("location", e.Name_.Location()))
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
		e.AddField(f)
	}
	e.BaseDecl.EndToken = p.Prev()

	return nil
}

func (e *Enum) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) io.Error {
	if _, exists := visitedDecls.Get(e); exists {
		return nil
	}
	defer visitedDecls.Set(e)

	return e.BaseDecl.LinkParents(p, visitedDecls, cycleMap)
}

func (*Enum) LinkMethods(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return nil
}

func (e *Enum) Semantic(p ast.SemanticParser) io.Error {
	return e.BaseDecl.Semantic(p)
}
