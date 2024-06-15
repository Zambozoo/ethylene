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
	Type         ast.Type
	GenericTypes []ast.Type
	EndToken     token.Token
}

func (g *Generic) Location() token.Location {
	return token.LocationBetween(g.Type, &g.EndToken)
}

func (g *Generic) String() string {
	return fmt.Sprintf("Generic{Type:%s,GenericParameters:%s}", g.Type, strs.Strings(g.GenericTypes))
}

func (g *Generic) ExtendsAsPointer(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (g *Generic) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return g.Equals(ctx, parent)
}

func (g *Generic) Equals(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	if _, ok := other.(*Generic); ok {
		panic("not implemented")
	}

	return false, nil
}

func (g *Generic) Declaration(ctx ast.TypeContext) (ast.Declaration, io.Error) {
	panic("not implemented")
}
