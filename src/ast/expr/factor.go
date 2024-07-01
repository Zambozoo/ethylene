package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Multiply represents expressions of the form
//
//	EXPR '*' EXPR
type Multiply struct{ Binary }

func (m *Multiply) String() string {
	return fmt.Sprintf("%s * %s", m.Left.String(), m.Right.String())
}

func (m *Multiply) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return m.Binary.MustBothExtendOne(p, type_.NewInteger(), type_.NewFloat())
}

// Divide represents expressions of the form
//
//	EXPR '/' EXPR
type Divide struct{ Binary }

func (d *Divide) String() string {
	return fmt.Sprintf("%s / %s", d.Left.String(), d.Right.String())
}

func (d *Divide) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return d.Binary.MustBothExtendOne(p, type_.NewInteger(), type_.NewFloat())
}

// Modulo represents expressions of the form
//
//	EXPR '%' EXPR
type Modulo struct{ Binary }

func (m *Modulo) String() string {
	return fmt.Sprintf("%s %% %s", m.Left.String(), m.Right.String())
}

func (m *Modulo) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	return m.Binary.MustBothExtendOne(p, type_.NewInteger())
}

func syntaxFactor(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxUnaryPost(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_STAR:
			p.Next()
			r, err := syntaxUnaryPost(p)
			if err != nil {
				return nil, err
			}

			expr = &Multiply{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_DIVIDE:
			p.Next()
			r, err := syntaxUnaryPost(p)
			if err != nil {
				return nil, err
			}

			expr = &Divide{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_MODULO:
			p.Next()
			r, err := syntaxUnaryPost(p)
			if err != nil {
				return nil, err
			}

			expr = &Modulo{
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
