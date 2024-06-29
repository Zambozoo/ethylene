package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

type Lambda struct {
	StartToken token.Token
	Type       ast.FunType
	Parameters []string
	Stmt       ast.Statement
}

func (l *Lambda) Location() token.Location {
	return token.LocationBetween(&l.StartToken, l.Stmt)
}

func (l *Lambda) String() string {
	return fmt.Sprintf("Lambda{Type:%s, Parameters:%s, Stmt:%s}", l.Type.String(), l.Parameters, l.Stmt.String())
}

func (l *Lambda) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	var err io.Error
	l.StartToken, err = p.Consume(token.TOK_LAMBDA)
	if err != nil {
		return nil, err
	}

	var ok bool
	t, err := p.ParseType()
	if err != nil {
		return nil, err
	}

	if l.Type, ok = t.(ast.FunType); !ok {
		return nil, io.NewError("expected a function type for lambda", zap.Any("location", t.Location()))
	}

	if _, err := p.Consume(token.TOK_COLON); err != nil {
		return nil, err
	}
	if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
		return nil, err
	}

	if !p.Match(token.TOK_RIGHTPAREN) {
		for {
			tok, err := p.Consume(token.TOK_IDENTIFIER)
			if err != nil {
				return nil, err
			}

			l.Parameters = append(l.Parameters, tok.Value)
			if p.Match(token.TOK_RIGHTPAREN) {
				break
			}

			if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return nil, err
			}
		}
	}

	if l.Type.Arity() != len(l.Parameters) {
		return nil, io.NewError("arity of lambda does not match number of parameters",
			zap.Any("expected", l.Type.Arity()),
			zap.Any("actual", len(l.Parameters)),
			zap.Any("location", l.Location()),
		)
	}

	l.Stmt, err = p.ParseStmt()
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *Lambda) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}
