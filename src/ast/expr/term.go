package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Add represents expressions of the form
//
//	EXPR + EXPR
type Add struct{ Binary }

func (a *Add) String() string {
	return fmt.Sprintf("Add{Left:%s,Right:%s}", a.Left.String(), a.Right.String())
}

func (a *Add) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return a.Binary.MustBothExtendOne(p, &type_.Integer{}, &type_.Float{})
}

// Substract represents expressions of the form
//
//	EXPR - EXPR
type Subtract struct{ Binary }

func (s *Subtract) String() string {
	return fmt.Sprintf("Substract{Left:%s,Right:%s}", s.Left.String(), s.Right.String())
}

func (s *Subtract) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return s.Binary.MustBothExtendOne(p, &type_.Integer{}, &type_.Float{})
}

func syntaxTerm(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxFactor(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_PLUS:
			p.Next()
			r, err := syntaxFactor(p)
			if err != nil {
				return nil, err
			}

			expr = &Add{
				Binary: Binary{
					Left:  expr,
					Right: r,
				},
			}
		case token.TOK_MINUS:
			p.Next()
			r, err := syntaxFactor(p)
			if err != nil {
				return nil, err
			}

			expr = &Subtract{
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
