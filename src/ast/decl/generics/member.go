package generics

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/stringers"
)

type genericMember struct {
	ast.Member
	mapping []ast.Type
}

func (m *genericMember) String() string {
	return fmt.Sprintf("[%s]:%s", stringers.Join(m.mapping, ","), m.Member.String())
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
