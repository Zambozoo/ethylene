package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// UnsignedWaitRightShift represents expressions of the form
//
//	`wait` EXPR
type Wait struct {
	PrefixedToken
}

func (w *Wait) String() string {
	return fmt.Sprintf("Wait{Expr:%s}", w.Expr.String())
}

func (w *Wait) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	var err io.Error
	w.Token, err = p.Consume(token.TOK_WAIT)
	if err != nil {
		return nil, err
	}

	w.Expr, err = p.ParseExpr()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wait) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := w.Expr.Semantic(p)
	if err != nil {
		return nil, err
	} else if thread, ok := t.(*type_.Thread); !ok {
		return nil, io.NewError("wait requires a thread argument", zap.String("actual", t.String()))
	} else {
		return thread.Type, nil
	}
}
