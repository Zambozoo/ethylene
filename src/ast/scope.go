package ast

import (
	"geth-cody/compile/lexer/token"
)

type Scope struct {
	Declaration_ Declaration
	Token        *token.Token
	Parent       *Scope
	Variables    map[string]Variable
	Labels       map[string]Label
	// lambdaDepth is used to keep track of variable captures in closures
	lambdaDepth int
}

func NewScope(d Declaration) *Scope {
	return &Scope{
		Declaration_: d,
		Variables:    map[string]Variable{},
		Labels:       map[string]Label{},
	}
}

type Opt func(*Scope)

func WithLambda() func(*Scope) {
	return func(s *Scope) {
		s.lambdaDepth++
	}
}
func WithDeclaration(d Declaration) func(*Scope) {
	return func(s *Scope) {
		s.Declaration_ = d
	}
}
func (s *Scope) Wrap(opts ...Opt) {
	parent := *s
	*s = Scope{
		Declaration_: s.Declaration_,
		Parent:       &parent,
		Variables:    map[string]Variable{},
		Labels:       map[string]Label{},
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

func (s *Scope) Declaration() Declaration {
	return s.Declaration_
}
