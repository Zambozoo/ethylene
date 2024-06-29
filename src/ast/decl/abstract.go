package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/decl/generics"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Abstract struct {
	BaseDecl

	IsTailed bool

	SuperClass ast.DeclType           // Optional
	Parents_   data.Set[ast.DeclType] // Interfaces this decl implements
}

func (a *Abstract) Parents() data.Set[ast.DeclType] {
	return a.Parents_
}

func (*Abstract) IsInterface() bool {
	return false
}
func (*Abstract) IsAbstract() bool {
	return true
}
func (*Abstract) IsClass() bool {
	return false
}

func (a *Abstract) IsConstant() bool {
	return false
}

func (a *Abstract) SetConstant() {}

func newAbstract() *Abstract {
	return &Abstract{
		BaseDecl: newDecl(),
	}
}

func (a *Abstract) String() string {
	var parentsString string
	if len(a.Parents_) > 0 {
		parentsString = "<: [" + strings.Join(maps.Keys(a.Parents_), ",") + "]"
	}
	return fmt.Sprintf("abstract %s%s {\n%s\n%s\n%s\n%s}",
		a.Name().Value,
		parentsString,
		strings.Join(maps.Keys(a.Methods_), "\n"),
		strings.Join(maps.Keys(a.Members_), "\n"),
		strings.Join(maps.Keys(a.StaticMembers_), "\n"),
		strings.Join(maps.Keys(a.StaticMembers_), "\n"),
	)
}

func (a *Abstract) Syntax(p ast.SyntaxParser) (ast.Declaration, io.Error) {
	var err io.Error
	if a.BaseDecl.StartToken, err = p.Consume(token.TOK_ABSTRACT); err != nil {
		return nil, err
	} else if a.Name_, err = p.Consume(token.TOK_IDENTIFIER); err != nil {
		return nil, err
	}
	genericDecl, err := generics.Syntax(a, p)
	if err != nil {
		return nil, err
	} else if genericDecl != nil {
		p.UnwrapScope()
		p.WrapScope(genericDecl)
	}

	if p.Match(token.TOK_TILDE) {
		a.IsTailed = true
	}

	if p.Match(token.TOK_SUBTYPE) {
		if a.Parents_, err = p.ParseParentTypes(); err != nil {
			return nil, err
		}
	}

	if _, err := p.Consume(token.TOK_LEFTBRACE); err != nil {
		return nil, err
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		f, err := p.ParseField()
		if err != nil {
			return nil, err
		} else if _, ok := f.(ast.DeclField); ok {
			if _, ok := genericDecl.GenericParamIndex(f.Name().Value); ok {
				return nil, io.NewError("inner decl name duplicates generic type",
					zap.Stringer("decl", f.Name()),
					zap.Stringer("location", f.Location()),
				)
			}
		}
		if err := a.AddField(f); err != nil {
			return nil, err
		}
	}
	a.BaseDecl.EndToken = p.Prev()

	a.BaseDecl.Index, err = p.Types().NextAbstractIndex(a)
	if err != nil {
		return nil, err
	}

	if genericDecl != nil {
		return genericDecl, nil
	}
	return a, nil
}

func (a *Abstract) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	return a.LinkParentsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return a.Parents_, nil
		},
		func(parent ast.DeclType) io.Error {
			a.Parents_.Set(parent)
			return nil
		},
		visitedDecls, cycleMap)
}

func (a *Abstract) LinkParentsWithProvider(
	p ast.SemanticParser,
	parentProviderFunc generics.ParentProviderFunc,
	parentSetterFunc generics.ParentSetterFunc,
	visitedDecls *data.AsyncSet[ast.Declaration],
	cycleMap map[string]struct{},
) (data.Set[ast.DeclType], io.Error) {
	if _, exists := visitedDecls.Get(a); exists {
		return a.Parents_, nil
	}

	l := a.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return nil, io.NewError("cyclical inheritance",
			zap.Stringer("abstract", a.Name()),
			zap.Stringer("location", l),
		)
	}
	cycleMap[l.String()] = struct{}{}
	defer delete(cycleMap, l.String())

	ps, err := parentProviderFunc()
	if err != nil {
		return nil, err
	}
	for _, parent := range ps {
		parentDecl, err := parent.Declaration(p)
		if err != nil {
			return nil, err
		}

		if parentDecl.IsAbstract() {
			if a.SuperClass != nil {
				return nil, io.NewError("abstracts cannot implement multiple concrete parents",
					zap.Stringer("abstract", a.Name()),
					zap.Stringer("location", a.Location()),
					zap.Stringer("parent", parentDecl.Name()),
				)
			}
			a.SuperClass = parent
		} else if parentDecl.IsClass() {
			return nil, io.NewError("abstracts cannot implement classes",
				zap.Stringer("class", parentDecl.Name()),
				zap.Stringer("abstract", a.Name()),
				zap.Stringer("location", a.Location()),
			)
		} else if !parentDecl.IsInterface() {
			return nil, io.NewError("abstract cannot implement struct or enum parents",
				zap.Stringer("class", a.Name()),
				zap.Stringer("location", a.Location()),
				zap.Stringer("parent", parentDecl.Name()),
			)
		}

		parents, err := parentDecl.LinkParents(p, visitedDecls, cycleMap)
		if err != nil {
			return nil, err
		}
		for _, parent := range parents {
			parentSetterFunc(parent)
		}
	}
	visitedDecls.Set(a)

	ps, err = parentProviderFunc()
	if err != nil {
		return nil, err
	}
	return ps, a.BaseDecl.LinkParents(p, visitedDecls)
}

func (a *Abstract) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return a.LinkFieldsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return a.Parents_, nil
		},
		func() map[string]ast.Method {
			return a.Methods_
		},
		visitedDecls,
	)
}

func (a *Abstract) LinkFieldsWithProvider(p ast.SemanticParser, parentProviderFunc generics.ParentProviderFunc, methodProviderFunc generics.MethodProviderFunc, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(a); exists {
		return nil
	}
	ps, err := parentProviderFunc()
	if err != nil {
		return err
	}
	for _, parent := range ps {
		parentDecl, err := parent.Declaration(p)
		if err != nil {
			return err
		}

		if err := parentDecl.LinkFields(p, visitedDecls); err != nil {
			return err
		}
		if err := a.BaseDecl.ExtendsParent(p, methodProviderFunc, parentDecl, visitedDecls); err != nil {
			return err
		}
	}

	if a.SuperClass != nil {
		for name, m := range a.Members() {
			if _, ok := a.Members_[name]; !ok && !m.HasModifier(ast.MOD_PRIVATE) {
				a.Members_[name] = m
			}
		}
	}
	visitedDecls.Set(a)

	return a.BaseDecl.LinkFields(p, visitedDecls)
}

func (a *Abstract) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return a.BaseDecl.Semantic(p)
}

func (a *Abstract) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return a.Equals(p, parent)
}

func (a *Abstract) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	for _, p := range a.Parents_ {
		if equals, err := p.Equals(parser, parent); equals || err != nil {
			return equals, err
		}
	}

	return a.Equals(parser, parent)
}

func (a *Abstract) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherAbstract, ok := other.(*Abstract); ok {
		return a == otherAbstract, nil
	} else if otherDeclType, ok := other.(ast.DeclType); ok {
		otherDeclaration, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		}
		if otherAbstract, ok := otherDeclaration.(*Abstract); ok {
			return a == otherAbstract, nil
		}
	}

	return false, nil
}

func (a *Abstract) Concretize(mapping []ast.Type) ast.Type {
	return a
}

func (a *Abstract) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	return typeid.NewTypeID(parser.Types().AbstractIndex(a.Index), 0), nil
}
