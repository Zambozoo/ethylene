package type_

import (
	"geth-cody/ast"
)

// genericDecl is a generic wrapper around a declaration.
type genericDecl struct {
	ast.Declaration
	mapping map[string]ast.Type
}

func newGenericDecl(g *Generic, d ast.Declaration) ast.Declaration {
	generics := d.Generics()
	mapping := map[string]ast.Type{}
	for i, t := range generics {
		mapping[t.Name().Value] = g.GenericTypes[i]
	}

	if cd, ok := d.(ast.ChildDeclaration); ok {
		return &genericChildDecl{
			ChildDeclaration: cd,
			mapping:          mapping,
		}
	}

	return &genericDecl{
		Declaration: d,
		mapping:     mapping,
	}
}

func genericsMap(d ast.Declaration, mapping map[string]ast.Type) map[string]ast.DeclType {
	m := map[string]ast.DeclType{}
	for k, v := range d.GenericsMap() {
		t := v.Concretize(mapping)
		if c, ok := t.(*Composite); ok && c.IsGeneric() {
			m[k] = c
		}
	}

	return m
}
func generics(d ast.Declaration, mapping map[string]ast.Type) []ast.DeclType {
	s := []ast.DeclType{}
	for _, v := range d.Generics() {
		t := v.Concretize(mapping)
		if c, ok := t.(*Composite); ok && c.IsGeneric() {
			s = append(s, c)
		}
	}

	return s
}
func members(d ast.Declaration, mapping map[string]ast.Type) map[string]ast.Member {
	m := map[string]ast.Member{}
	for k, v := range d.Members() {
		m[k] = &genericMember{Member: v, mapping: mapping}
	}

	return m
}
func methods(d ast.Declaration, mapping map[string]ast.Type) map[string]ast.Method {
	m := map[string]ast.Method{}
	for k, v := range d.Methods() {
		m[k] = &genericMethod{Method: v, mapping: mapping}
	}

	return m
}

func (g *genericDecl) GenericsMap() map[string]ast.DeclType {
	return genericsMap(g.Declaration, g.mapping)
}
func (g *genericDecl) Generics() []ast.DeclType {
	return generics(g.Declaration, g.mapping)
}

func (g *genericDecl) Members() map[string]ast.Member {
	return members(g.Declaration, g.mapping)
}
func (g *genericDecl) Methods() map[string]ast.Method {
	return methods(g.Declaration, g.mapping)
}

type genericMember struct {
	ast.Member
	mapping map[string]ast.Type
}

func (g *genericMember) Type() ast.Type {
	return g.Member.Type().Concretize(g.mapping)
}

type genericMethod struct {
	ast.Method
	mapping map[string]ast.Type
}

func (g *genericMethod) Type() ast.Type {
	return g.Method.Type().Concretize(g.mapping)
}

// genericDecl is a generic wrapper around a declaration.
type genericChildDecl struct {
	ast.ChildDeclaration
	mapping map[string]ast.Type
}

func (g *genericChildDecl) GenericsMap() map[string]ast.DeclType {
	return genericsMap(g.ChildDeclaration, g.mapping)
}
func (g *genericChildDecl) Generics() []ast.DeclType {
	return generics(g.ChildDeclaration, g.mapping)
}

func (g *genericChildDecl) Members() map[string]ast.Member {
	return members(g.ChildDeclaration, g.mapping)
}
func (g *genericChildDecl) Methods() map[string]ast.Method {
	return methods(g.ChildDeclaration, g.mapping)
}

func (g *genericChildDecl) Parents() []ast.Type {
	parents := g.ChildDeclaration.Parents()
	s := make([]ast.Type, len(parents))
	for i, p := range parents {
		s[i] = p.Concretize(g.mapping)
	}

	return s
}
