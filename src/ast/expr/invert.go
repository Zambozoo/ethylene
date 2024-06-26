package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

type Negation struct {
	PrefixedUnary
}

func (n *Negation) String() string {
	return fmt.Sprintf("-%s", n.Expr.String())
}

func (n *Negation) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := n.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(p, t, type_.NewInteger(), type_.NewFloat()); err != nil {
		return nil, err
	}

	return t, nil
}

type Bang struct {
	PrefixedUnary
}

func (b *Bang) String() string {
	return fmt.Sprintf("!%s", b.Expr.String())
}

func (b *Bang) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := b.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(p, t, &type_.Word{}, type_.NewBoolean()); err != nil {
		return nil, err
	}

	return t, nil
}
