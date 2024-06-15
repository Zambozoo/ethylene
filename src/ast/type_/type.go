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
		t, err = (&Composite{}).Syntax(p)
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
				return nil, io.NewError("type array size is larger than max signed integer limit", zap.Any("token", tok))
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
			var types []ast.Type
			for {
				t, err := p.ParseType()
				if err != nil {
					return nil, err
				}
				types = append(types, t)
				if p.Match(token.TOK_RIGHTBRACKET) {
					break
				}

				if _, err := p.Consume(token.TOK_COMMA); err != nil {
					return nil, err
				}
			}
			t = &Generic{
				Type:         t,
				GenericTypes: types,
				EndToken:     p.Prev(),
			}
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
				return nil, io.NewError("type array size is larger than max signed integer limit", zap.Any("token", tok))
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
		default:
			return t, nil
		}
	}
}

func MustExtend(ctx ast.TypeContext, child ast.Type, parent ast.Type, parents ...ast.Type) (ast.Type, io.Error) {
	parents = append(parents, parent)
	for _, p := range parents {
		if extends, err := child.Extends(ctx, p); err != nil {
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
