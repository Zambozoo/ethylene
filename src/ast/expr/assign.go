package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Assign represents expressions of the form
//
//	EXPR '=' EXPR
type Assign struct {
	Binary
}

func (a *Assign) String() string {
	return fmt.Sprintf("Assign{Left:%s, Right:%s}", a.Left.String(), a.Right.String())
}
func (a *Assign) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&Ternary{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_ASSIGN) {
			r, err := (&Ternary{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			a := &Assign{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}

			if !assignable(a.Left) {
				return nil, io.NewError("invalid target for assign expresion", zap.Any("location", a.Left.Location()))
			}

			expr = a
		} else {
			return expr, nil
		}
	}
}

func assignable(e ast.Expression) bool {
	switch e.(type) {
	case *Dereference, *Identifier, *Field, *Access:
		return true
	default:
		return false
	}
}

func (a *Assign) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	left, err := a.Left.Semantic(p)
	var t ast.Type
	if err != nil {
		return nil, err
	}

	right, err := a.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(right, left); err != nil {
		return nil, err
	}

	return t, nil
}
