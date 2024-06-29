package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Expr represents an expression statement
type Expr struct {
	Expr     ast.Expression
	EndToken token.Token
}

func (e *Expr) Location() token.Location {
	return token.LocationBetween(e.Expr, &e.EndToken)
}

func (e *Expr) String() string {
	return fmt.Sprintf("%s;", e.Expr.String())
}

func (e *Expr) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	e.Expr, err = p.ParseExpr()
	if err != nil {
		return err
	}

	e.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	return err
}

func (e *Expr) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := e.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	return t, nil
}
