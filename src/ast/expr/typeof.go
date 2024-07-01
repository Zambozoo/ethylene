package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"

	"go.uber.org/zap"
)

// TypeOf represents expressions of the form
//
//	`@` EXPR
type TypeOf struct {
	PrefixedToken
}

func (t *TypeOf) String() string {
	return fmt.Sprintf("@%s", t.Expr.String())
}

func (to *TypeOf) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := to.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if ok, err := t.Extends(p, &type_.Null{}); err != nil {
		return nil, err
	} else if ok {
		return nil, io.NewError("typeof expression cannot be called on null type", zap.Stringer("location", to.Location()))
	}

	return type_.NewTypeID(), nil
}
