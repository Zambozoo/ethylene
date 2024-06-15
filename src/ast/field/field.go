package field

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type field interface {
	ast.Field
	Syntax(p ast.SyntaxParser) io.Error
}
type Modifiers map[ast.Modifier]struct{}

func (ms *Modifiers) HasModifier(m ast.Modifier) bool {
	_, ok := (*ms)[m]
	return ok
}

func (ms *Modifiers) String() string {
	return strs.Strings(maps.Keys(*ms))
}

func Syntax(p ast.SyntaxParser) (ast.Field, io.Error) {
	var f field

	startToken := p.Peek()
	modifiers, err := ast.SyntaxModifiers(p)
	if err != nil {
		return nil, err
	}

	switch t := p.Peek(); t.Type {
	case token.TOK_VAR:
		f = &Member{Modifiers: modifiers, StartToken: startToken}
	case token.TOK_FUN:
		f = &Method{Modifiers: modifiers, StartToken: startToken}
	case token.TOK_CLASS, token.TOK_ABSTRACT, token.TOK_INTERFACE, token.TOK_STRUCT, token.TOK_ENUM:
		f = &Decl{Modifiers: modifiers, StartToken: startToken}
	default:
		return nil, io.NewError("expected field", zap.Any("token", t))
	}

	if err := f.Syntax(p); err != nil {
		return nil, err
	}

	return f, nil
}
