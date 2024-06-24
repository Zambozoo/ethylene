package ast

import (
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

type Label struct {
	token *token.Token
}

func (l *Label) Token() *token.Token {
	return l.token
}

func (s *Scope) AddLabel(tok *token.Token) io.Error {
	if oldVar, ok := s.Variables[tok.Value]; ok {
		return io.NewError("label already declared in this scope",
			zap.Any("old variable", oldVar),
			zap.Any("new variable", tok),
		)
	}

	s.Labels[tok.Value] = Label{
		token: tok,
	}
	return nil
}

func (s *Scope) GetLabel(tok token.Token) (*Label, io.Error) {
	for scope := s; scope != nil; scope = scope.Parent {
		if scope.lambdaDepth != s.lambdaDepth {
			break
		}

		if label, ok := scope.Labels[tok.Value]; ok {
			return &label, nil
		}
	}

	return nil, io.NewError("label not found", zap.String("token", tok.String()))
}
