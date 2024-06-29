package generics

import (
	"geth-cody/ast"
)

type genericMethod struct {
	ast.Method
	mapping []ast.Type
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
