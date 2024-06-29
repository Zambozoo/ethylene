package field

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Enum struct {
	StartToken token.Token
	EndToken   token.Token
	Type_      ast.Type
	Expression ast.Expression
}

func (e *Enum) Name() *token.Token {
	return &e.StartToken
}

func (e *Enum) HasModifier(m ast.Modifier) bool {
	switch m {
	case ast.MOD_PUBLIC, ast.MOD_STATIC:
		return true
	default:
		return false
	}
}

func (e *Enum) Type() ast.Type {
	return e.Type_
}

func (e *Enum) Location() token.Location {
	return token.LocationBetween(&e.StartToken, &e.EndToken)
}

func (e *Enum) String() string {
	return fmt.Sprintf("EnumField{Name:%s,Expr:%s}",
		e.Name().String(),
		e.Expression.String(),
	)
}

func (e *Enum) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error

	e.StartToken = p.Peek()
	if e.StartToken.Type != token.TOK_IDENTIFIER {
		_, err = p.Consume(token.TOK_IDENTIFIER)
		return err
	}

	if e.Expression, err = p.ParseExpr(); err != nil {
		return err
	}

	e.EndToken = p.Peek()

	return nil
}

func (e *Enum) Semantic(p ast.SemanticParser) io.Error {
	t, err := e.Expression.Semantic(p)
	if err != nil {
		return err
	}

	_, err = type_.MustExtend(p, t, &type_.Void{})
	return err
}
