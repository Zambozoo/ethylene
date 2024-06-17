package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Ternary represents expressions of the form
//
//	EXPR `?` EXPR `:` EXPR
type Ternary struct {
	Condition ast.Expression
	Then      ast.Expression
	Else      ast.Expression
}

func (t *Ternary) Location() token.Location {
	return token.LocationBetween(t.Condition, t.Else)
}

func (t *Ternary) String() string {
	return fmt.Sprintf("Ternary{Condition:%s,Then:%s,Else:%s}",
		t.Condition.String(),
		t.Then.String(),
		t.Else.String(),
	)
}

func (t *Ternary) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&Nullary{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_QUESTIONMARK) {
			left, err := (&Nullary{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			if _, err := p.Consume(token.TOK_COLON); err != nil {
				return nil, err
			}

			right, err := (&Nullary{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			expr = &Ternary{
				Condition: expr,
				Then:      left,
				Else:      right,
			}
		} else {
			return expr, nil
		}
	}
}

func (t *Ternary) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	cond, err := t.Condition.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(cond, &type_.Boolean{}); err != nil {
		return nil, err
	}

	then, err := t.Then.Semantic(p)
	if err != nil {
		return nil, err
	}

	else_, err := t.Else.Semantic(p)
	if err != nil {
		return nil, err
	}

	return type_.Union{then, else_}, nil
}
