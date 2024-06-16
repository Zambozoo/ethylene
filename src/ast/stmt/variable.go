package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Variable represents a variable declaration
type Var struct {
	BoundedStmt
	Type_ ast.Type
	Name_ token.Token
	Expr  ast.Expression
}

func (v *Var) Name() *token.Token {
	return &v.Name_
}

func (v *Var) Type() ast.Type {
	return v.Type_
}

func (v *Var) String() string {
	return fmt.Sprintf("Variable{Name:%s,Expr:%s}", v.Name_.Value, v.Expr.String())
}

func (v *Var) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	v.BoundedStmt.StartToken, err = p.Consume(token.TOK_VAR)
	if err != nil {
		return err
	}

	v.Type_, err = p.ParseType()
	if err != nil {
		return err
	}

	v.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	if _, err := p.Consume(token.TOK_ASSIGN); err != nil {
		return err
	}

	if p.Peek().Type != token.TOK_SEMICOLON {
		v.Expr, err = p.ParseExpr()
		if err != nil {
			return err
		}
	}

	v.BoundedStmt.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	return err
}

func (v *Var) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	if v.Expr == nil {
		return nil, nil
	}

	t, err := v.Expr.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(t, v.Type_); err != nil {
		return nil, err
	}

	if err := p.Scope().AddVariable(v); err != nil {
		return nil, err
	}

	return nil, nil
}
