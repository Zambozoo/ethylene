package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"geth-cody/stringers"

	"go.uber.org/zap"
)

// GenericType represents a type with generic type parameters
type Generic struct {
	Constant     bool
	Context_     ast.TypeContext
	Type         *Lookup
	GenericTypes []ast.Type
	EndToken     token.Token
}

func (g *Generic) Name() token.Token {
	return g.Type.Name()
}

func (g *Generic) Context() ast.TypeContext {
	return g.Context_
}

func (g *Generic) Location() *token.Location {
	return token.LocationBetween(g.Type, &g.EndToken)
}

func (g *Generic) String() string {
	return fmt.Sprintf("%s%s", g.Type, stringers.Join(g.GenericTypes, ","))
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
	}

	d, err := g.Declaration(p)
	if err != nil {
		return false, err
	}

	otherD, err := gOther.Declaration(p)
	if err != nil {
		return false, err
	}

	return d.Equals(p, otherD)
}

func (g *Generic) Declaration(p ast.SemanticParser) (ast.Declaration, io.Error) {
	d, err := g.Type.Context_.Declaration(g.Type.Tokens)
	if err != nil {
		return nil, err
	}

	if len(d.Generics()) != len(g.GenericTypes) {
		return nil, io.NewError("number of generic arguments in type differs from the number of generic parameters in declaration",
			zap.Int("expected", len(d.Generics())),
			zap.Int("actual", len(g.GenericTypes)),
			zap.Stringer("location", g.Location()),
		)
	}

	return p.WrapDeclWithGeneric(d, g.GenericTypes), nil
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

func (g *Generic) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	d, err := g.Declaration(parser)
	if err != nil {
		return nil, err
	}
	tid, err := d.TypeID(parser)
	if err != nil {
		return nil, err
	}

	index := tid.Index()
	if g.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, tid.ListIndex()), nil
}

func (g *Generic) IsConcrete() bool {
	for _, t := range g.GenericTypes {
		if !t.IsConcrete() {
			return false
		}
	}

	return true
}

func (*Generic) IsFieldable() bool {
	return false
}
