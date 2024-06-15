package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// RightShift represents expressions of the form
//
//	EXPR `>>` EXPR
type RightShift struct{ Binary }

func (r *RightShift) String() string {
	return fmt.Sprintf("RightShift{Left:%s,Right:%s}", r.Left.String(), r.Right.String())
}

func (r *RightShift) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	left, err := r.Left.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := p.TypeContext().MustExtend(left, &type_.Word{}); err != nil {
		return nil, err
	}

	right, err := r.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := p.TypeContext().MustExtend(right, &type_.Integer{}); err != nil {
		return nil, err
	}

	return left, nil
}

// LeftShift represents expressions of the form
//
//	EXPR `<<` EXPR
type LeftShift struct{ Binary }

func (l *LeftShift) String() string {
	return fmt.Sprintf("LeftShift{Left:%s,Right:%s}", l.Left.String(), l.Right.String())
}

func (l *LeftShift) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	left, err := l.Left.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := p.TypeContext().MustExtend(left, &type_.Word{}); err != nil {
		return nil, err
	}

	right, err := l.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := p.TypeContext().MustExtend(right, &type_.Integer{}); err != nil {
		return nil, err
	}

	return left, nil
}

// UnsignedRightShift represents expressions of the form
//
//	EXPR `>>>` EXPR
type UnsignedRightShift struct{ Binary }

func (u *UnsignedRightShift) String() string {
	return fmt.Sprintf("UnsignedRightShift{Left:%s,Right:%s}", u.Left.String(), u.Right.String())
}

func (u *UnsignedRightShift) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	left, err := u.Left.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := p.TypeContext().MustExtend(left, &type_.Word{}); err != nil {
		return nil, err
	}

	right, err := u.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := p.TypeContext().MustExtend(right, &type_.Integer{}); err != nil {
		return nil, err
	}

	return left, nil
}

func syntaxShift(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxTerm(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_SHIFTLEFT:
			p.Next()
			r, err := syntaxTerm(p)
			if err != nil {
				return nil, err
			}

			expr = &LeftShift{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_SHIFTRIGHT:
			p.Next()
			r, err := syntaxTerm(p)
			if err != nil {
				return nil, err
			}

			expr = &RightShift{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_SHIFTUNSIGNEDRIGHT:
			p.Next()
			r, err := syntaxTerm(p)
			if err != nil {
				return nil, err
			}

			expr = &UnsignedRightShift{
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
