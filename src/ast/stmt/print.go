package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Print represents a print statement
type Print struct {
	BoundedStmt
	Expr ast.Expression
}

func (p *Print) String() string {
	return fmt.Sprintf("print(%s);", p.Expr.String())
}

func (p *Print) Syntax(parser ast.SyntaxParser) io.Error {
	var err io.Error
	p.BoundedStmt.StartToken, err = parser.Consume(token.TOK_PRINT)
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

func (p *Print) Semantic(parser ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := p.Expr.Semantic(parser)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(parser, t, &type_.String{}); err != nil {
		return nil, err
	}

	return nil, nil
}
