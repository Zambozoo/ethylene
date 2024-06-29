package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

// Hash represents expressions of the form
//
//	'#' EXPR
type Hash struct {
	PrefixedToken
}

func (h *Hash) String() string {
	return fmt.Sprintf("#%s", h.Expr.String())
}

func (h *Hash) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return &type_.Integer{}, nil
}
