package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// LessThan represents expressions of the form
//
//	EXPR '<' EXPR
type LessThan struct{ Binary }

func (e *LessThan) String() string {
	return fmt.Sprintf("%s < %s", e.Left.String(), e.Right.String())
}

func (l *LessThan) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// LessThanOrEqual represents expressions of the form
//
//	EXPR '<=' EXPR
type LessThanOrEqual struct{ Binary }

func (e *LessThanOrEqual) String() string {
	return fmt.Sprintf("%s <= %s", e.Left.String(), e.Right.String())
}

func (l *LessThanOrEqual) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// GreaterThan represents expressions of the form
//
//	EXPR '>=' EXPR
type GreaterThan struct{ Binary }

func (e *GreaterThan) String() string {
	return fmt.Sprintf("%s > %s", e.Left.String(), e.Right.String())
}

func (g *GreaterThan) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// GreaterThanOrEqual represents expressions of the form
//
//	EXPR '>=' EXPR
type GreaterThanOrEqual struct{ Binary }

func (e *GreaterThanOrEqual) String() string {
	return fmt.Sprintf("%s >= %s", e.Left.String(), e.Right.String())
}

func (g *GreaterThanOrEqual) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// Spaceship represents expressions of the form
//
//	EXPR '<=>' EXPR
type Spaceship struct{ Binary }

func (e *Spaceship) String() string {
	return fmt.Sprintf("%s <=> %s", e.Left.String(), e.Right.String())
}

func (s *Spaceship) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

func syntaxCompare(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxShift(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_SPACESHIP:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &Spaceship{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_LESSTHAN:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &LessThan{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_LESSTHANEQUAL:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &LessThanOrEqual{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_GREATERTHAN:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &GreaterThan{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_GREATERTHANEQUAL:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &GreaterThanOrEqual{
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
