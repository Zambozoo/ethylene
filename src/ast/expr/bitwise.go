package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// BitwiseXor represents expressions of the form
//
//	EXPR '^' EXPR
type BitwiseXor struct {
	Binary
}

func (b *BitwiseXor) String() string {
	return fmt.Sprintf("%s ^ %s", b.Left.String(), b.Right.String())
}

func (b *BitwiseXor) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&BitwiseAnd{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_BITXOR) {
			r, err := (&BitwiseAnd{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			expr = &BitwiseXor{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		} else {
			return expr, nil
		}
	}
}

func (b *BitwiseXor) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return b.Binary.MustBothExtendOne(p, &type_.Boolean{}, &type_.Word{})
}

// BitwiseAnd represents expressions of the form
//
//	EXPR '&' EXPR
type BitwiseAnd struct {
	Binary
}

func (b *BitwiseAnd) String() string {
	return fmt.Sprintf("%s & %s", b.Left.String(), b.Right.String())
}

func (b *BitwiseAnd) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxEqual(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_BITAND) {
			r, err := syntaxEqual(p)
			if err != nil {
				return nil, err
			}

			expr = &BitwiseAnd{Binary: Binary{
				Left:  expr,
				Right: r,
			},
			}
		} else {
			return expr, nil
		}
	}
}

func (b *BitwiseAnd) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return b.Binary.MustBothExtendOne(p, &type_.Boolean{}, &type_.Word{})
}

// BitwiseOr represents expressions of the form
//
//	EXPR '|' EXPR
type BitwiseOr struct {
	Binary
}

func (b *BitwiseOr) String() string {
	return fmt.Sprintf("%s | %s", b.Left.String(), b.Right.String())
}

func (b *BitwiseOr) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := (&BitwiseXor{}).Syntax(p)
	if err != nil {
		return nil, err
	}

	for {
		if p.Match(token.TOK_BITOR) {
			r, err := (&BitwiseXor{}).Syntax(p)
			if err != nil {
				return nil, err
			}

			expr = &BitwiseOr{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		} else {
			return expr, nil
		}
	}
}

func (b *BitwiseOr) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return b.Binary.MustBothExtendOne(p, &type_.Boolean{}, &type_.Word{})
}
