package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

// IncrementPrefix represents expressions of the form
//
//	'++' EXPR
type IncrementPrefix struct {
	PrefixedUnary
}

func (i *IncrementPrefix) String() string {
	return fmt.Sprintf("++%s", i.Expr.String())
}

func (i *IncrementPrefix) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := i.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(p, t, &type_.Integer{}); err != nil {
		return nil, err
	}

	return t, nil
}

// IncrementSuffix represents expressions of the form
//
//	EXPR '++'
type IncrementSuffix struct {
	SuffixedUnary
}

func (i *IncrementSuffix) String() string {
	return fmt.Sprintf("%s++", i.Expr.String())
}

func (i *IncrementSuffix) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := i.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(p, t, &type_.Integer{}); err != nil {
		return nil, err
	}

	return t, nil
}

// DecrementPrefix represents expressions of the form
//
//	'--' EXPR
type DecrementPrefix struct {
	PrefixedUnary
}

func (d *DecrementPrefix) String() string {
	return fmt.Sprintf("--%s", d.Expr.String())
}

func (d *DecrementPrefix) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := d.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(p, t, &type_.Integer{}); err != nil {
		return nil, err
	}

	return t, nil
}

// DecrementSuffix represents expressions of the form
//
//	EXPR '--'
type DecrementSuffix struct {
	SuffixedUnary
}

func (d *DecrementSuffix) String() string {
	return fmt.Sprintf("%s--", d.Expr.String())
}

func (d *DecrementSuffix) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := d.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(p, t, &type_.Integer{}); err != nil {
		return nil, err
	}

	return t, nil
}
