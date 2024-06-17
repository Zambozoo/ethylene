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

func (g *Generic) Name() token.Token {
	return g.Type.Name()
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

func (g *Generic) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	cDecl, err := g.Declaration()
	if err != nil {
		return false, err
	}
	cChildDecl, ok := cDecl.(ast.ChildDeclaration)
	if !ok {
		return false, nil
	}

	mapping := map[string]ast.Type{}
	for i, t := range cChildDecl.Generics() {
		gt := g.GenericTypes[i]
		if c, ok := gt.(*Composite); !ok || !c.IsGeneric() {
			mapping[t.Name().Value] = gt
		}
	}

	pDecl, ok := parent.(ast.DeclType)
	if ok {
		return false, nil
	}

	for _, parentType := range cChildDecl.Parents() {
		concreteParentType := parentType.Concretize(mapping)
		if ok, err := concreteParentType.ExtendsAsPointer(pDecl); err != nil || ok {
			return ok, err
		}
	}

	return false, nil
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

func (g *Generic) Syntax(p ast.SyntaxParser) io.Error {
	for {
		t, err := p.ParseType()
		if err != nil {
			return err
		}
		g.GenericTypes = append(g.GenericTypes, t)

		if p.Match(token.TOK_RIGHTBRACKET) {
			break
		}

		if _, err := p.Consume(token.TOK_COMMA); err != nil {
			return err
		}
	}

	g.Context_ = p.TypeContext()
	g.EndToken = p.Prev()
	return nil
}

func (g *Generic) Concretize(mapping map[string]ast.Type) ast.Type {
	genericTypes := make([]ast.Type, len(g.GenericTypes))
	for i, t := range g.GenericTypes {
		genericTypes[i] = t.Concretize(mapping)
	}

	return &Generic{
		Context_:     g.Context_,
		Type:         g.Type,
		GenericTypes: genericTypes,
		EndToken:     g.EndToken,
	}
}
