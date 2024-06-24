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

func NewScope() *Scope {
	return &Scope{
		Variables: map[string]Variable{},
		Labels:    map[string]Label{},
	}
}

type Opt func(*Scope)

func WithLambda() func(*Scope) {
	return func(s *Scope) {
		s.lambdaDepth++
	}
}
func (s *Scope) Wrap(opts ...Opt) {
	parent := *s
	*s = Scope{
		Parent:    &parent,
		Variables: map[string]Variable{},
		Labels:    map[string]Label{},
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
