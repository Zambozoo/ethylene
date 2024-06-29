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
	Constant     bool
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
	return fmt.Sprintf("%s[%s]", g.Type, strs.Strings(g.GenericTypes, ","))
}

func (g *Generic) Key(p ast.SemanticParser) (string, io.Error) {
	decl, err := g.Declaration(p)
	if err != nil {
		return "", err
	}

	return decl.Concretize(g.Context_.TopScope().Generics()).Key(p)
}

func (g *Generic) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	decl, err := g.Declaration(p)
	if err != nil {
		return false, err
	}

	return decl.ExtendsAsPointer(p, parent)
}

func (g *Generic) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return g.Equals(p, parent)
}

func (g *Generic) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	gOther, ok := other.(*Generic)
	if !ok {
		return false, nil
	} else if ok, err := g.Type.Equals(p, gOther.Type); err != nil || !ok {
		return false, err
	}

	for i, childGenericArg := range g.GenericTypes {
		parentGenericArg := gOther.GenericTypes[i]
		if ok, err := childGenericArg.Equals(p, parentGenericArg); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (g *Generic) Declaration(p ast.SemanticParser) (ast.Declaration, io.Error) {
	d, err := g.Type.Declaration(p)
	if err != nil {
		return nil, err
	}

	mapping := map[string]ast.Type{}
	for _, t := range g.GenericTypes {
		k, err := t.Key(p)
		if err != nil {
			return nil, err
		}
		mapping[k] = t
	}

	return p.NewGenericDecl(d, g.GenericTypes, mapping), nil
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

func (g *Generic) Concretize(mapping []ast.Type) ast.Type {
	genericTypes := make([]ast.Type, len(g.GenericTypes))
	for i, t := range g.GenericTypes {
		genericTypes[i] = t.Concretize(mapping)
	}

	return &Generic{
		Constant:     g.Constant,
		Context_:     g.Context_,
		Type:         g.Type,
		GenericTypes: genericTypes,
		EndToken:     g.EndToken,
	}
}

func (g *Generic) IsConstant() bool {
	return g.Constant
}
func (g *Generic) SetConstant() {
	g.Constant = true
}
