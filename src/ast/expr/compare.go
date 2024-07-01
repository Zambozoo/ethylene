package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// LessThan represents expressions of the form
//
//	EXPR '<' EXPR
type LessThan struct{ Binary }

func (e *LessThan) String() string {
	return fmt.Sprintf("%s < %s", e.Left.String(), e.Right.String())
}

func (l *LessThan) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := l.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := l.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if err := compareTypeCheck(p, left, right); err != nil {
		return nil, err
	}

	return type_.NewBoolean(), nil
}

// LessThanOrEqual represents expressions of the form
//
//	EXPR '<=' EXPR
type LessThanOrEqual struct{ Binary }

func (e *LessThanOrEqual) String() string {
	return fmt.Sprintf("%s <= %s", e.Left.String(), e.Right.String())
}

func (l *LessThanOrEqual) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := l.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := l.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if err := compareTypeCheck(p, left, right); err != nil {
		return nil, err
	}

	return type_.NewBoolean(), nil
}

// GreaterThan represents expressions of the form
//
//	EXPR '>=' EXPR
type GreaterThan struct{ Binary }

func (e *GreaterThan) String() string {
	return fmt.Sprintf("%s > %s", e.Left.String(), e.Right.String())
}

func (g *GreaterThan) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := g.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := g.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if err := compareTypeCheck(p, left, right); err != nil {
		return nil, err
	}

	return type_.NewBoolean(), nil
}

// GreaterThanOrEqual represents expressions of the form
//
//	EXPR '>=' EXPR
type GreaterThanOrEqual struct{ Binary }

func (e *GreaterThanOrEqual) String() string {
	return fmt.Sprintf("%s >= %s", e.Left.String(), e.Right.String())
}

func (g *GreaterThanOrEqual) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := g.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := g.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if err := compareTypeCheck(p, left, right); err != nil {
		return nil, err
	}

	return type_.NewBoolean(), nil
}

// Spaceship represents expressions of the form
//
//	EXPR '<=>' EXPR
type Spaceship struct{ Binary }

func (e *Spaceship) String() string {
	return fmt.Sprintf("%s <=> %s", e.Left.String(), e.Right.String())
}

func (s *Spaceship) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := s.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := s.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if err := compareTypeCheck(p, left, right); err != nil {
		return nil, err
	}

	return type_.NewInteger(), nil
}

func syntaxCompare(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxShift(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_SUBTYPE:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &SuperType{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_SUPERTYPE:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &SubType{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_SUBTYPE:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &SuperType{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_SUPERTYPE:
			p.Next()
			r, err := syntaxShift(p)
			if err != nil {
				return nil, err
			}

			expr = &SubType{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
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

func compareTypeCheck(p ast.SemanticParser, left, right ast.Type) io.Error {
	if eq, err := left.Equals(p, right); err != nil {
		return err
	} else if !eq {
		return io.NewError("compare type mismatch",
			zap.Stringer("left", left),
			zap.Stringer("right", right),
		)
	}

	return nil
}
