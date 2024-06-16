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
	return fmt.Sprintf("Negation{Expr:%s}", n.Expr.String())
}

func (n *Negation) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecodes
	t, err := n.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(t, &type_.Integer{}, &type_.Float{}); err != nil {
		return nil, err
	}

	return t, nil
}

type Bang struct {
	PrefixedUnary
}

func (b *Bang) String() string {
	return fmt.Sprintf("Bang{Expr:%s}", b.Expr.String())
}

func (b *Bang) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := b.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := type_.MustExtend(t, &type_.Word{}, &type_.Boolean{}); err != nil {
		return nil, err
	}

	return t, nil
}
