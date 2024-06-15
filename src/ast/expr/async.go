package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Variable represents a variable declaration
type Async struct {
	StartToken token.Token
	CallExpr   *Call
}

func (a *Async) Location() token.Location {
	return token.LocationBetween(&a.StartToken, a.CallExpr)
}

func (a *Async) String() string {
	return fmt.Sprintf("Async{Expr:%s}", a.CallExpr.String())
}

func (a *Async) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	var err io.Error
	a.StartToken, err = p.Consume(token.TOK_ASYNC)
	if err != nil {
		return nil, err
	}

	expr, err := p.ParseExpr()
	if err != nil {
		return nil, err
	} else if callExpr, ok := expr.(*Call); !ok {
		return nil, io.NewError("async argument must be a function call", zap.Any("location", a.Location()))
	} else {
		a.CallExpr = callExpr
	}

	return a, nil
}

func (a *Async) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := a.CallExpr.Semantic(p)
	if err != nil {
		return nil, err
	}

	return &type_.Thread{Type: t}, nil
}
