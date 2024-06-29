package generics

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Arg is the variable value: `int` in `List[int]`
type Arg struct {
	ast.Type
	Index int
}

// Decl represents the base struct for a generic declaration
type Decl struct {
	ast.Declaration
	SymbolSlice []ast.Type
}

func NewDecl(d ast.Declaration, symbolSlice []ast.Type, symbolMap map[string]ast.Type) *Decl {
	b := &Decl{
		Declaration: d,
		SymbolSlice: symbolSlice,
	}

	for _, t := range symbolSlice {
		if param, ok := t.(*type_.Param); ok {
			param.Decl = b
		}
	}

	return b
}

func (b *Decl) Arity() int {
	return len(b.SymbolSlice)
}

func (d *Decl) Generics() []ast.Type {
	return d.SymbolSlice
}
func (b *Decl) Concretize_(p ast.SemanticParser, d ast.Declaration, args []ast.Type) (ast.Declaration, io.Error) {
	resultSlice := make([]ast.Type, 0, len(b.SymbolSlice))
	if b.Arity() != len(args) {
		return nil, io.NewError("Incorrect number of generic arguments",
			zap.Int("expected", b.Arity()),
			zap.Int("actual", len(args)),
			zap.Any("location", d.Location()),
		)
	}
	for i, arg := range args {
		// TODO: CHECK IF SUPER/SUB WORKS HERE
		switch t := arg.(type) {
		case *type_.Param:
			arg = &type_.Param{
				Token: t.Token,
				Decl:  t.Decl,
				Index: i,
			}
		case *Arg:
			arg = &Arg{
				Type:  t.Type,
				Index: i,
			}
		default:
			arg = &Arg{
				Type:  t,
				Index: i,
			}
		}
		resultSlice = append(resultSlice, arg)
	}

	return &Decl{
		Declaration: d,
		SymbolSlice: resultSlice,
	}, nil
}

func (g *Decl) Parents() (data.Set[ast.DeclType], io.Error) {
	decl, ok := g.Declaration.(ast.ChildDeclaration)
	if !ok {
		return nil, nil
	}
	parents := data.Set[ast.DeclType]{}
	for _, parent := range decl.Parents() {
		parents.Set(parent.Concretize(g.SymbolSlice).(ast.DeclType))
	}

	return parents, nil
}

// TODO: ADD SUPPORT FOR SUB/SUPER
func Syntax(d ast.Declaration, p ast.SyntaxParser) (*Decl, io.Error) {
	b := &Decl{
		Declaration: d,
		SymbolSlice: []ast.Type{},
	}
	if p.Match(token.TOK_LEFTBRACKET) {
		for {
			tok, err := p.Consume(token.TOK_IDENTIFIER)
			if err != nil {
				return nil, err
			}

			t := &type_.Param{
				Context_: p.TypeContext(),
				Token:    tok,
				Index:    len(b.SymbolSlice),
				Decl:     b,
			}
			b.SymbolSlice = append(b.SymbolSlice, t)

			if p.Match(token.TOK_RIGHTBRACKET) {
				break
			} else if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, nil
	}

	return b, nil
}

func (d *Decl) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return d.Equals(p, parent)
}

func (d *Decl) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	ps, err := d.Parents()
	if err != nil {
		return false, err
	}
	for _, p := range ps {
		if equals, err := p.Equals(parser, parent); equals || err != nil {
			return equals, err
		}
	}

	return d.Equals(parser, parent)
}

func (d *Decl) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherDecl, ok := other.(*Decl); ok {
		return d == otherDecl, nil
	} else if otherDeclType, ok := other.(ast.DeclType); ok {
		otherDeclaration, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		}
		if otherDecl, ok := otherDeclaration.(*Decl); ok {
			return d == otherDecl, nil
		}
	}

	return false, nil
}

func (d *Decl) Key(p ast.SemanticParser) (string, io.Error) {
	l := d.Location()
	var keys string
	var spacer string
	for _, t := range d.SymbolSlice {
		k, err := t.Key(p)
		if err != nil {
			return "", err
		}
		keys += spacer + k
	}
	return fmt.Sprintf("%s:%s[%s]", l.String(), d.Name().Value, keys), nil
}

func (d *Decl) Concretize(mapping []ast.Type) ast.Type {
	return &Decl{
		Declaration: d,
		SymbolSlice: mapping,
	}
}

type ParentProviderFunc func() (data.Set[ast.DeclType], io.Error)
type ParentSetterFunc func(parent ast.DeclType) io.Error
type MethodProviderFunc func() map[string]ast.Method

type childDecl interface {
	LinkParentsWithProvider(p ast.SemanticParser, parentProviderFunc ParentProviderFunc, parentSetterFunc ParentSetterFunc, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error)
	LinkFieldsWithProvider(p ast.SemanticParser, parentProviderFunc ParentProviderFunc, methodProviderFunc MethodProviderFunc, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error
	ExtendsParent(p ast.SemanticParser, methodProviderFunc MethodProviderFunc, parent ast.Declaration, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error
}

func (d *Decl) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	if cd, ok := d.Declaration.(childDecl); ok {
		return cd.LinkParentsWithProvider(p,
			func() (data.Set[ast.DeclType], io.Error) {
				return d.Parents()
			},
			func(parent ast.DeclType) io.Error {
				ps, err := d.Parents()
				if err != nil {
					return err
				}
				ps.Set(parent.Concretize(d.SymbolSlice).(ast.DeclType))
				return nil
			},
			visitedDecls,
			cycleMap)
	}
	return d.Declaration.LinkParents(p, visitedDecls, cycleMap)
}

func (d *Decl) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if cd, ok := d.Declaration.(childDecl); ok {
		return cd.LinkFieldsWithProvider(p,
			func() (data.Set[ast.DeclType], io.Error) {
				return d.Parents()
			},
			func() map[string]ast.Method {
				return d.Methods()
			},
			visitedDecls,
		)
	}
	return d.Declaration.LinkFields(p, visitedDecls)
}

func (d *Decl) GenericParamIndex(name string) (int, bool) {
	for _, t := range d.SymbolSlice {
		if p, ok := t.(*type_.Param); ok && p.Token.Value == name {
			return p.Index, true
		}
	}
	return 0, false
}
