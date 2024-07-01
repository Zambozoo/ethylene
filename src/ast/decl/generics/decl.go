package generics

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"geth-cody/stringers"

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

func (d *Decl) String() string {
	return fmt.Sprintf("[%s]:%s", stringers.Join(d.SymbolSlice, ","), d.Declaration.String())
}

func NewDecl(d ast.Declaration, symbolSlice []ast.Type) *Decl {
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
			zap.Stringer("location", d.Location()),
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

func (g *Decl) Parents() data.Set[ast.DeclType] {
	decl, ok := g.Declaration.(ast.ChildDeclaration)

	if !ok {
		return nil
	}
	parents := data.Set[ast.DeclType]{}
	for _, parent := range decl.Parents() {
		parents.Set(parent.Concretize(g.SymbolSlice).(ast.DeclType))
	}

	return parents
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
	for _, p := range d.Parents() {
		if equals, err := p.Equals(parser, parent); equals || err != nil {
			return equals, err
		}
	}

	return d.Equals(parser, parent)
}

func (d *Decl) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	otherDecl, ok := other.(*Decl)
	if !ok {
		otherDeclType, ok := other.(*type_.Generic)
		if !ok {
			return false, nil
		}

		otherD, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		}
		otherDecl = otherD.(*Decl)
	}

	for i, t := range d.SymbolSlice {
		otherT := otherDecl.SymbolSlice[i]
		eq, err := t.Equals(p, otherT)
		if err != nil || !eq {
			return false, err
		}
	}

	root := d.rootDeclaration()
	otherRoot := otherDecl.rootDeclaration()
	return root == otherRoot, nil
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
				return d.Parents(), nil
			},
			func(parent ast.DeclType) io.Error {
				d.Parents().Set(parent.Concretize(d.SymbolSlice).(ast.DeclType))
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
				return d.Parents(), nil
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

func (d *Decl) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	ids := make([]uint64, len(d.SymbolSlice))
	for _, t := range d.SymbolSlice {
		id, err := t.TypeID(parser)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id.ID())
	}

	tid, err := d.Declaration.TypeID(parser)
	if err != nil {
		return nil, err
	}

	return typeid.NewTypeID(tid.Index(), parser.Types().ListIndex(ids)), nil
}

func (d *Decl) IsConcrete() bool {
	for _, t := range d.SymbolSlice {
		if !t.IsConcrete() {
			return false
		}
	}

	return true
}

func (d *Decl) Super() (ast.DeclType, bool) {
	cd, ok := d.Declaration.(ast.ChildDeclaration)
	if !ok {
		return nil, ok
	}
	return cd.Super()
}

func (d *Decl) IsTailed() bool {
	return d.Declaration.IsTailed()
}

func (d *Decl) rootDeclaration() ast.Declaration {
	decl := d.Declaration
	for {
		g, ok := decl.(*Decl)
		if ok {
			return decl
		}
		decl = g.Declaration
	}
}
