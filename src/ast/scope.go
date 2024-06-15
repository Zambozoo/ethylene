package ast

import (
	"geth-cody/compile/lexer/token"
)

type Scope struct {
	Token     *token.Token
	Parent    *Scope
	Variables map[string]Variable
	Labels    map[string]Label
	// lambdaDepth is used to keep track of variable captures in closures
	lambdaDepth int
}

type Opt func(*Scope)

func WithLambda() func(*Scope) {
	return func(s *Scope) {
		s.lambdaDepth++
	}
}
func (s *Scope) Wrap(opts ...Opt) {
	*s = Scope{
		Parent:    s,
		Variables: map[string]Variable{},
	}

	for _, opt := range opts {
		opt(s)
	}
}

func (s *Scope) WrapDecl(token *token.Token) {
	*s = Scope{
		Token:     token,
		Parent:    s,
		Variables: map[string]Variable{},
	}
}

func (s *Scope) Unwrap() {
	*s = *s.Parent
}
