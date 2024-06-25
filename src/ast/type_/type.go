package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"math"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Syntax(p ast.SyntaxParser) (ast.Type, io.Error) {
	var t ast.Type
	var err io.Error
	if p.Peek().Type == token.TOK_IDENTIFIER {
		t, err = (&Composite{Context_: p.TypeContext()}).Syntax(p)
		if err != nil {
			return nil, err
		}
	} else {
		t, err = syntaxPrimitive(p)
		if err != nil {
			return nil, err
		}
	}

	// generic or array
	if p.Match(token.TOK_LEFTBRACKET) {
		switch tok := p.Peek(); tok.Type {
		case token.TOK_INTEGER:
			p.Next()
			if tok.Integer > math.MaxInt {
				return nil, io.NewError("type array size is larger than max signed integer limit", zap.String("token", tok.String()))
			}
			endTok, err := p.Consume(token.TOK_RIGHTBRACKET)
			if err != nil {
				return nil, err
			}

			t = &Array{
				Type:     t,
				Size:     int64(tok.Integer),
				EndToken: endTok,
			}
		case token.TOK_TILDE:
			p.Next()
			tok, err := p.Consume(token.TOK_RIGHTBRACKET)
			if err != nil {
				return nil, err
			}

			t = &Array{
				Type:     t,
				Size:     -1,
				EndToken: tok,
			}
		default:
			declType, ok := t.(ast.DeclType)
			if !ok {
				return t, nil
			}

			g := &Generic{Type: declType}
			g.Syntax(p)
			t = g
		}
	}

	if p.Match(token.TOK_TILDE) {
		declType, ok := t.(ast.DeclType)
		if !ok {
			return nil, io.NewError("invalid tailed type.",
				zap.Any("type", t),
				zap.Any("location", t.Location()),
			)
		}

		size := int64(-1)
		if p.Match(token.TOK_INTEGER) {
			tok := p.Prev()
			if tok.Integer > math.MaxInt {
				return nil, io.NewError("type tail size is larger than max signed integer limit", zap.String("token", tok.String()))
			}
			size = int64(tok.Integer)
		}
		t = &Tailed{
			Type:     declType,
			Size:     size,
			EndToken: p.Prev(),
		}
	}

	for {
		switch tok := p.Peek(); tok.Type {
		case token.TOK_STAR:
			t = &Pointer{
				Type:     t,
				EndToken: p.Next(),
			}
		case token.TOK_LEFTPAREN:
			p.Next()
			var types []ast.Type
			if !p.Match(token.TOK_RIGHTPAREN) {
				for {
					t, err := p.ParseType()
					if err != nil {
						return nil, err
					}

					types = append(types, t)
					if p.Match(token.TOK_RIGHTPAREN) {
						break
					}

					if _, err := p.Consume(token.TOK_COMMA); err != nil {
						return nil, err
					}
				}
			}
			t = &Function{
				ReturnType_:     t,
				ParameterTypes_: types,
				EndToken:        p.Prev(),
			}
		case token.TOK_LEFTBRACKET:
			p.Next()
			tok, err := p.Consume(token.TOK_INTEGER)
			if err != nil {
				return nil, err
			} else if tok.Integer > math.MaxInt {
				return nil, io.NewError("type array size is larger than max signed integer limit", zap.String("token", tok.String()))
			}
			endTok, err := p.Consume(token.TOK_RIGHTBRACKET)
			if err != nil {
				return nil, err
			}

			t = &Array{
				Type:     t,
				Size:     int64(tok.Integer),
				EndToken: endTok,
			}
		case token.TOK_DOLLAR:
			p.Next()
			if t.IsConstant() {
				return nil, io.NewError("invalid double constant type",
					zap.Any("type", t),
					zap.Any("location", t.Location()),
				)
			}
			t.SetConstant()
		default:
			return t, nil
		}
	}
}

func MustExtend(child ast.Type, parent ast.Type, parents ...ast.Type) (ast.Type, io.Error) {
	parents = append(parents, parent)
	for _, p := range parents {
		if extends, err := child.Extends(p); err != nil {
			return nil, err
		} else if extends {
			return p, nil
		}
	}

	var expectedField zapcore.Field
	if len(parents) > 0 {
		expectedField = zap.Any("expected one of", parents)
	} else {
		expectedField = zap.Any("expected", parent)
	}

	return nil, io.NewError("type mismatch",
		expectedField,
		zap.Any("actual", child),
		zap.Any("location", child.Location()),
	)
}
