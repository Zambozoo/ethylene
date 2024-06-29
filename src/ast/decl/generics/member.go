package generics

import (
	"geth-cody/ast"
)

type genericMember struct {
	ast.Member
	mapping []ast.Type
}

func (g *genericMember) Type() ast.Type {
	return g.Member.Type().Concretize(g.mapping)
}

func (g *Decl) Members() map[string]ast.Member {
	m := map[string]ast.Member{}
	for k, v := range g.Declaration.Members() {
		m[k] = &genericMember{Member: v, mapping: g.SymbolSlice}
	}

	return m
}
