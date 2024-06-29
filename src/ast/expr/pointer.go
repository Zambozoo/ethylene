package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Reference represents expressions of the form
//
//	'&' EXPR
type Reference struct {
	PrefixedToken
}

func (r *Reference) String() string {
	return fmt.Sprintf("%s&", r.Expr.String())
}

func (r *Reference) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
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
	t, err := d.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if p, ok := t.(*type_.Pointer); ok {
		t = p.Type
	} else if a, ok := t.(*type_.Array); ok {
		t = a.Type
	} else {
		return nil, io.NewError("deference of non-pointer and non-array type",
			zap.Stringer("location", d.Location()),
			zap.Stringer("type", t),
		)
	}

	return t, nil
}

func (d *Dereference) String() string {
	return fmt.Sprintf("%s*", d.Expr.String())
}
