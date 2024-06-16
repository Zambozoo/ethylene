package expr

import (
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

type Binary struct {
	Left  ast.Expression
	Right ast.Expression
}

func (b *Binary) Location() token.Location {
	return token.LocationBetween(b.Left, b.Right)
}

func (b *Binary) MustBothExtendOne(p ast.SemanticParser, parent ast.Type, parents ...ast.Type) (ast.Type, io.Error) {
	left, err := b.Left.Semantic(p)
	if err != nil {
		return nil, err
	}
	t, err := type_.MustExtend(left, parent, parents...)
	if err != nil {
		return nil, err
	}

	right, err := b.Right.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(right, t); err != nil {
		return nil, err
	}

	return t, nil
}

type Unary struct {
	Token token.Token
	Expr  ast.Expression
}

type PrefixedUnary Unary

func (s *PrefixedUnary) Location() token.Location {
	return token.LocationBetween(&s.Token, s.Expr)
}

type SuffixedUnary Unary

func (s *SuffixedUnary) Location() token.Location {
	return token.LocationBetween(s.Expr, &s.Token)
}

type tokenExpr struct {
	Token token.Token
	Expr  ast.Expression
}

type SuffixedToken tokenExpr

func (s *SuffixedToken) Location() token.Location {
	return token.LocationBetween(s.Expr, &s.Token)
}

type PrefixedToken tokenExpr

func (p *PrefixedToken) Location() token.Location {
	return token.LocationBetween(&p.Token, p.Expr)
}

func Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	return (&Assign{}).Syntax(p)
}

// syntaxUnaryPost parses postfix unary expressions, including increment and decrement
// operators, type casts, array access, function calls, and field access. It starts with an
// expression parsed by syntaxUnaryPre and then applies any postfix operations to it. The function
// handles errors by returning them along with a nil expression. This function is crucial for
// parsing expressions that modify or access the value of an expression after the expression itself.
func syntaxUnaryPost(p ast.SyntaxParser) (ast.Expression, io.Error) {
	expr, err := syntaxUnaryPre(p)
	if err != nil {
		return nil, err
	}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_INC:
			if !assignable(expr) {
				return nil, io.NewError("invalid target for increment expresion", zap.Any("location", expr.Location()))
			}
			expr = &IncrementSuffix{
				SuffixedUnary: SuffixedUnary{
					Expr:  expr,
					Token: p.Next(),
				},
			}
		case token.TOK_DEC:
			if !assignable(expr) {
				return nil, io.NewError("invalid target for decrement expresion", zap.Any("location", expr.Location()))
			}
			expr = &DecrementSuffix{
				SuffixedUnary: SuffixedUnary{
					Expr:  expr,
					Token: p.Next(),
				},
			}
		case token.TOK_LEFTBRACE:
			p.Next()
			t, err := p.ParseType()
			if err != nil {
				return nil, err
			}

			tok, err := p.Consume(token.TOK_RIGHTBRACE)
			if err != nil {
				return nil, err
			}

			expr = &Cast{
				Type: t,
				SuffixedToken: SuffixedToken{
					Token: tok,
					Expr:  expr,
				},
			}
		case token.TOK_LEFTBRACKET:
			p.Next()
			r, err := p.ParseExpr()
			if err != nil {
				return nil, err
			}

			tok, err := p.Consume(token.TOK_RIGHTBRACKET)
			if err != nil {
				return nil, err
			}

			expr = &Access{
				Left:     expr,
				Right:    r,
				EndToken: tok,
			}
		case token.TOK_LEFTPAREN:
			p.Next()
			var exprs []ast.Expression
			if !p.Match(token.TOK_RIGHTPAREN) {
				for {
					expr, err := p.ParseExpr()
					if err != nil {
						return nil, err
					}

					exprs = append(exprs, expr)
					if p.Match(token.TOK_RIGHTPAREN) {
						break
					}

					if _, err := p.Consume(token.TOK_COMMA); err != nil {
						return nil, err
					}
				}
			}
			expr = &Call{
				SuffixedToken: SuffixedToken{
					Token: p.Prev(),
					Expr:  expr,
				},
				Exprs: exprs,
			}
		case token.TOK_PERIOD:
			p.Next()
			tok, err := p.Consume(token.TOK_IDENTIFIER)
			if err != nil {
				return nil, err
			}

			expr = &Field{
				SuffixedToken: SuffixedToken{
					Token: tok,
					Expr:  expr,
				},
			}
		default:
			return expr, nil
		}
	}
}

func syntaxUnaryPre(p ast.SyntaxParser) (ast.Expression, io.Error) {
	switch t := p.Peek(); t.Type {
	case token.TOK_MINUS:
		tok := p.Next()
		expr, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		return &Negation{
			PrefixedUnary: PrefixedUnary{
				Token: tok,
				Expr:  expr,
			},
		}, nil
	case token.TOK_BANG:
		tok := p.Next()
		expr, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		return &Bang{
			PrefixedUnary: PrefixedUnary{
				Token: tok,
				Expr:  expr,
			},
		}, nil
	case token.TOK_INC:
		tok := p.Next()
		e, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		expr := &IncrementPrefix{
			PrefixedUnary: PrefixedUnary{
				Token: tok,
				Expr:  e,
			},
		}

		if !assignable(expr.Expr) {
			return nil, io.NewError("invalid target for increment expresion", zap.Any("location", expr.Location()))
		}

		return expr, nil
	case token.TOK_DEC:
		tok := p.Next()
		e, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		expr := &DecrementPrefix{
			PrefixedUnary: PrefixedUnary{
				Token: tok,
				Expr:  e,
			},
		}
		if !assignable(expr.Expr) {
			return nil, io.NewError("invalid target for decrement expresion", zap.Any("location", expr.Location()))
		}

		return expr, nil
	case token.TOK_AT:
		tok := p.Next()
		expr, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		return &TypeOf{
			PrefixedToken: PrefixedToken{
				Token: tok,
				Expr:  expr,
			},
		}, nil
	case token.TOK_HASHTAG:
		tok := p.Next()
		expr, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		return &Hash{
			PrefixedToken: PrefixedToken{
				Token: tok,
				Expr:  expr,
			},
		}, nil
	case token.TOK_BITAND:
		tok := p.Next()
		e, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		expr := &Reference{
			PrefixedToken: PrefixedToken{
				Token: tok,
				Expr:  e,
			},
		}
		if !assignable(expr.Expr) {
			return nil, io.NewError("invalid target for reference expresion", zap.Any("location", expr.Location()))
		}

		return expr, nil
	case token.TOK_STAR:
		tok := p.Next()
		expr, err := syntaxUnaryPre(p)
		if err != nil {
			return nil, err
		}

		return &Dereference{
			PrefixedToken: PrefixedToken{
				Token: tok,
				Expr:  expr,
			},
		}, nil
	default:
		return syntaxPrimary(p)
	}
}

func syntaxPrimary(p ast.SyntaxParser) (ast.Expression, io.Error) {
	switch t := p.Peek(); t.Type {
	case token.TOK_INTEGER:
		return &Integer{Token: p.Next()}, nil
	case token.TOK_FLOAT:
		return &Float{Token: p.Next()}, nil
	case token.TOK_CHARACTER:
		return &Character{Token: p.Next()}, nil
	case token.TOK_TYPE:
		return (&TypeField{}).Syntax(p)
	case token.TOK_NEW:
		return (&New{}).Syntax(p)
	case token.TOK_ASYNC:
		return (&Async{}).Syntax(p)
	case token.TOK_WAIT:
		return (&Wait{}).Syntax(p)
	case token.TOK_LAMBDA:
		return (&Lambda{}).Syntax(p)
	case token.TOK_LEFTPAREN:
		if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
			return nil, err
		}

		expr, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}

		if _, err := p.Consume(token.TOK_RIGHTPAREN); err != nil {
			return nil, err
		}

		return expr, nil
	case token.TOK_STRING:
		return &String{Token: p.Next()}, nil
	case token.TOK_IDENTIFIER:
		return (&Identifier{}).Syntax(p)
	case token.TOK_FALSE:
		return &False{Token: p.Next()}, nil
	case token.TOK_TRUE:
		return &True{Token: p.Next()}, nil
	case token.TOK_NULL:
		return &Null{Token: p.Next()}, nil
	case token.TOK_THIS:
		return &This{Token: p.Next()}, nil
	case token.TOK_SUPER:
		return &Super{Token: p.Next()}, nil
	default:
		return nil, io.NewError("unexpected token", zap.Any("token", t))
	}
}
