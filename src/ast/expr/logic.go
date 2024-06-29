package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// And represents expressions of the form
//
//	EXPR '&&' EXPR
type And struct {
	Binary
}

func (a *And) String() string {
	return fmt.Sprintf("%s && %s", a.Left.String(), a.Right.String())
}

func (a *And) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&BitwiseOr{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_AND) {
			r, err := (&BitwiseOr{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			expr = &And{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		} else {
			return expr, nil
		}
	}
}

func (a *And) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return a.Binary.MustBothExtendOne(p, &type_.Boolean{})
}

// Or represents expressions of the form
//
//	EXPR '||' EXPR
type Or struct {
	Binary
}

func (o *Or) String() string {
	return fmt.Sprintf("%s || %s", o.Left.String(), o.Right.String())
}
func (o *Or) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&And{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_OR) {
			r, err := (&And{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			expr = &Or{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		} else {
			return expr, nil
		}
	}
}

func (o *Or) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return o.Binary.MustBothExtendOne(p, &type_.Boolean{})
}
