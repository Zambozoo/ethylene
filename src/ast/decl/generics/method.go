package generics

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/stringers"
)

type genericMethod struct {
	ast.Method
	mapping []ast.Type
}

func (m *genericMethod) String() string {
	return fmt.Sprintf("[%s]:%s", stringers.Join(m.mapping, ","), m.Method.String())
}

func (g *genericMethod) Type() ast.Type {
	return g.Method.Type().Concretize(g.mapping)
}

func (g *Decl) Methods() map[string]ast.Method {
	m := map[string]ast.Method{}
	for k, v := range g.Declaration.Methods() {
		m[k] = &genericMethod{Method: v, mapping: g.SymbolSlice}
	}

	return m
}
