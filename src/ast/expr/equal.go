package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Equal represents expressions of the form
//
//	EXPR '==' EXPR
type Equal struct {
	Binary
}

func (e *Equal) String() string {
	return fmt.Sprintf("Equal{Left:%s,Right:%s}", e.Left.String(), e.Right.String())
}

func (e *Equal) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// BangEqual represents expressions of the form
//
//	EXPR '!=' EXPR
type BangEqual struct {
	Binary
}

func (e *BangEqual) String() string {
	return fmt.Sprintf("BangEqual{Left:%s,Right:%s}", e.Left.String(), e.Right.String())
}

func (b *BangEqual) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

func syntaxEqual(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxCompare(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_EQUAL:
			p.Next()
			r, err := syntaxCompare(p)
			if err != nil {
				return nil, err
			}

			expr = &Equal{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_BANGEQUAL:
			p.Next()
			r, err := syntaxCompare(p)
			if err != nil {
				return nil, err
			}

			expr = &BangEqual{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		default:
			return expr, nil
		}
	}
}
