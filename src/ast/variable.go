package ast

import (
	"fmt"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

type Variable interface {
	fmt.Stringer
	Name() *token.Token
	Type() Type
}

func (s *Scope) AddVariable(variable Variable) io.Error {
	if oldVar, ok := s.Variables[variable.Name().Value]; ok {
		return io.NewError("variable already declared in this scope",
			zap.Stringer("old variable", oldVar),
			zap.Stringer("new variable", variable.Name()),
		)
	}

	s.Variables[variable.Name().Value] = variable
	return nil
}

func (s *Scope) GetVariable(tok token.Token) (variable Variable, isCaptured bool, err io.Error) {
	for scope := s; scope != nil; scope = scope.Parent {
		if info, ok := scope.Variables[tok.Value]; ok {
			return info, scope.lambdaDepth == s.lambdaDepth, nil
		}
	}

	return nil, false, io.NewError("variable not found", zap.Stringer("token", &tok))
}
