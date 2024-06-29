package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Cast represents expressions of the form
//
//	EXPR '{' EXPR '}'
type Cast struct {
	SuffixedToken
	Type ast.Type
}

func (c *Cast) String() string {
	return fmt.Sprintf("%s{%s}", c.Expr.String(), c.Type.String())
}

func (c *Cast) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := c.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if !type_.CastPrimitive(p, t, c.Type) {
		if extends, err := c.Type.Extends(p, t); err == nil {
			return nil, err
		} else if !extends {
			if extends, err := t.Extends(p, c.Type); err == nil {
				return nil, err
			} else if !extends {
				return nil, io.NewError("cast type mismatch",
					zap.Stringer("source", c.Type),
					zap.Stringer("target", t),
				)
			}
		}
	}

	return c.Type, nil
}
