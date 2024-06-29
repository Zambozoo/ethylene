package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Delete represents a delete statement
type Delete struct {
	BoundedStmt
	Expr ast.Expression
}

func (d *Delete) String() string {
	return fmt.Sprintf("delete(%s);", d.Expr.String())
}

func (d *Delete) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	d.BoundedStmt.StartToken, err = p.Consume(token.TOK_DELETE)
	if err != nil {
		return err
	}

	d.Expr, err = p.ParseExpr()
	if err != nil {
		return err
	}

	d.BoundedStmt.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	return err
}

func (d *Delete) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return nil, nil
}
