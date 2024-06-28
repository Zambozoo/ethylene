package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

// Cast represents expressions of the form
//
//	EXPR '{' EXPR '}'
type Cast struct {
	SuffixedToken
	Type ast.Type
}

func (c *Cast) String() string {
	return fmt.Sprintf("Cast{Expr:%s}", c.Expr.String())
}

func (c *Cast) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := c.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if !type_.CastPrimitive(p, t, c.Type) {
		if _, err := type_.MustExtend(p, c.Type, t); err == nil {
			return nil, err
		}
	}

	return c.Type, nil
}
