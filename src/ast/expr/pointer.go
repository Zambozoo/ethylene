package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

// Reference represents expressions of the form
//
//	'&' EXPR
type Reference struct {
	PrefixedToken
}

func (r *Reference) String() string {
	return fmt.Sprintf("Reference{Expr:%s}", r.Expr.String())
}

func (r *Reference) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := r.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	return &type_.Pointer{Type: t}, nil
}

// Dereference represents expressions of the form
//
//	'*' EXPR
type Dereference struct {
	PrefixedToken
}

func (d *Dereference) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

func (d *Dereference) String() string {
	return fmt.Sprintf("Dereference{Expr:%s}", d.Expr.String())
}
