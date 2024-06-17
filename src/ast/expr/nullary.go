package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Nullary represents expressions of the form
//
//	EXPR '??' EXPR
type Nullary struct {
	Binary
}

func (n *Nullary) String() string {
	return fmt.Sprintf("Nullary{Left:%s,Right:%s}", n.Left.String(), n.Right.String())
}

func (n *Nullary) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&Or{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_NULLARY) {
			r, err := (&Or{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			expr = &Nullary{
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

func (n *Nullary) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	left, err := n.Left.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, ok := left.(*type_.Pointer); !ok {
		return nil, io.NewError("left operand of nullary operator must be a pointer", zap.Any("location", n.Left.Location()))
	}

	right, err := n.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, ok := right.(*type_.Pointer); !ok {
		return nil, io.NewError("right operand of nullary operator must be a pointer", zap.Any("location", n.Left.Location()))
	}

	return type_.Union{right, left}, nil
}
