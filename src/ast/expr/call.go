package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
	"geth-cody/strs"

	"go.uber.org/zap"
)

// Call represents expressions of the form
//
//	EXPR '(' EXPR ')'
type Call struct {
	SuffixedToken
	Exprs []ast.Expression
}

func (c *Call) String() string {
	return fmt.Sprintf("%s(%s)", c.Expr.String(), strs.Strings(c.Exprs, ","))
}

func (c *Call) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	left, err := c.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	ft, ok := left.(*type_.Function)
	if !ok {
		return nil, io.NewError("left expression of call must be a function", zap.Any("location", c.Location()))
	}

	if ft.Arity() != len(c.Exprs) {
		return nil, io.NewError("number of arguments does not match function signature", zap.Any("location", c.Location()))
	}

	for i, expr := range c.Exprs {
		t, err := expr.Semantic(p)

		if err != nil {
			return nil, err
		}

		expectedType := ft.ParameterTypes()[i]
		if _, err := type_.MustExtend(p, t, expectedType); err != nil {
			return nil, err
		}
	}

	return ft.ReturnType(), nil
}
