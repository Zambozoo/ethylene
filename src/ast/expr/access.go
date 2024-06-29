package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Access represents expressions of the form
//
//	EXPR '[' EXPR ']'
type Access struct {
	Left     ast.Expression
	Right    ast.Expression
	EndToken token.Token
}

func (a *Access) Location() token.Location {
	return token.LocationBetween(a.Left, &a.EndToken)
}

func (a *Access) String() string {
	return fmt.Sprintf("%s[%s]", a.Left.String(), a.Right.String())
}

func (a *Access) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := a.Left.Semantic(p)
	var t ast.Type
	if err != nil {
		return nil, err
	} else if ptr, ok := left.(*type_.Pointer); ok {
		t = ptr.Type
	} else if arr, ok := left.(*type_.Array); ok {
		t = arr.Type
	} else {
		return nil, io.NewError("left expression of access must be a pointer or array", zap.Any("location", a.Location()))
	}

	right, err := a.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(p, left, right, &type_.Integer{}); err != nil {
		return nil, err
	}

	return t, nil
}
