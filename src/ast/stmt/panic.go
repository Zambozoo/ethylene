package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Panic represents a panic statement
type Panic struct {
	BoundedStmt
	Expr ast.Expression
}

func (p *Panic) String() string {
	return fmt.Sprintf("Panic{Expr:%s}", p.Expr.String())
}

func (p *Panic) Syntax(parser ast.SyntaxParser) io.Error {
	var err io.Error
	p.BoundedStmt.StartToken, err = parser.Consume(token.TOK_PANIC)
	if err != nil {
		return err
	}

	p.Expr, err = parser.ParseExpr()
	if err != nil {
		return err
	}

	p.BoundedStmt.EndToken, err = parser.Consume(token.TOK_SEMICOLON)
	return err
}

func (p *Panic) Semantic(parser ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := p.Expr.Semantic(parser)
	if err != nil {
		return nil, err
	} else if _, err := parser.TypeContext().MustExtend(t, &type_.String{}); err != nil {
		return nil, err
	}

	return nil, nil
}
