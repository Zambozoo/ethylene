package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

// GenericType represents a type with generic type parameters
type Generic struct {
	Context_     ast.TypeContext
	Type         ast.DeclType
	GenericTypes []ast.Type
	EndToken     token.Token
}

func (g *Generic) Context() ast.TypeContext {
	return g.Context_
}

func (g *Generic) Location() token.Location {
	return token.LocationBetween(g.Type, &g.EndToken)
}

func (g *Generic) String() string {
	return fmt.Sprintf("Generic{Type:%s,GenericParameters:%s}", g.Type, strs.Strings(g.GenericTypes))
}

func (g *Generic) Key() string {
	var s string
	var spacer string
	for _, t := range g.GenericTypes {
		s += spacer + t.Key()
		spacer = ","
	}

	return fmt.Sprintf("%s[%s]", g.Type.Key(), s)
}

func (g *Generic) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (g *Generic) Extends(parent ast.Type) (bool, io.Error) {
	return g.Equals(parent)
}

func (g *Generic) Equals(other ast.Type) (bool, io.Error) {
	gOther, ok := other.(*Generic)
	if !ok {
		return false, nil
	} else if ok, err := g.Type.Equals(gOther.Type); err != nil || !ok {
		return false, err
	}

	for i, childGenericArg := range g.GenericTypes {
		parentGenericArg := gOther.GenericTypes[i]
		if ok, err := childGenericArg.Equals(parentGenericArg); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (g *Generic) Declaration() (ast.Declaration, io.Error) {
	return g.Type.Declaration()
}
