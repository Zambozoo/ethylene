package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Return represents a return statement
type Return struct {
	BoundedStmt
	Expr ast.Expression
}

func (p *Return) String() string {
	return fmt.Sprintf("Return{Expr:%s}", p.Expr.String())
}

func (r *Return) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	r.BoundedStmt.StartToken, err = p.Consume(token.TOK_RETURN)
	if err != nil {
		return err
	}

	if p.Peek().Type != token.TOK_SEMICOLON {
		r.Expr, err = p.ParseExpr()
		if err != nil {
			return err
		}
	}

	r.BoundedStmt.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	return err
}

func (r *Return) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	if r.Expr == nil {
		return nil, nil
	}

	t, err := r.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	return t, nil
}
