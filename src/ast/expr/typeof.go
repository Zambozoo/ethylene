package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

// TypeOf represents expressions of the form
//
//	`@` EXPR
type TypeOf struct {
	PrefixedToken
}

func (t *TypeOf) String() string {
	return fmt.Sprintf("TypeOf{Expr:%s}", t.Expr.String())
}

func (to *TypeOf) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	_, err := to.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	return &type_.TypeID{}, nil
}
