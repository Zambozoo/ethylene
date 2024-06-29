package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/decl/generics"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Abstract struct {
	BaseDecl

	IsTailed bool

	SuperClass ast.DeclType           // Optional
	Implements data.Set[ast.DeclType] // Interfaces this decl implements
}

func (a *Abstract) Parents() data.Set[ast.DeclType] {
	return a.Implements
}

func (a *Abstract) SetTailed() io.Error {
	a.IsTailed = true
	return nil
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
	return fmt.Sprintf("Abstract{Name: %s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		a.Name().Value,
		maps.Keys(a.Implements),
		strings.Join(maps.Keys(a.Methods_), ","),
		strings.Join(maps.Keys(a.Members_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
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
		if a.Implements, err = p.ParseParentTypes(); err != nil {
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
					zap.Any("decl", f.Name()),
					zap.Any("location", f.Location()),
				)
			}
		}
		if err := a.AddField(f); err != nil {
			return nil, err
		}
	}
	a.BaseDecl.EndToken = p.Prev()

	if genericDecl != nil {
		return genericDecl, nil
	}
	return a, nil
}

func (a *Abstract) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	return a.LinkParentsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return a.Implements, nil
		},
		func(parent ast.DeclType) io.Error {
			a.Implements.Set(parent)
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
		return a.Implements, nil
	}

	l := a.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return nil, io.NewError("cyclical inheritance",
			zap.Any("abstract", a.Name()),
			zap.Any("location", l),
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
					zap.Any("abstract", a.Name()),
					zap.Any("location", a.Location()),
					zap.Any("parent", parentDecl.Name()),
				)
			}
			a.SuperClass = parent
		} else if parentDecl.IsClass() {
			return nil, io.NewError("abstracts cannot implement classes",
				zap.Any("class", parentDecl.Name()),
				zap.Any("abstract", a.Name()),
				zap.Any("location", a.Location()),
			)
		} else if !parentDecl.IsInterface() {
			return nil, io.NewError("abstract cannot implement struct or enum parents",
				zap.Any("class", a.Name()),
				zap.Any("location", a.Location()),
				zap.Any("parent", parentDecl.Name()),
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
			return a.Implements, nil
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
	for _, p := range a.Implements {
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

func (a *Abstract) Key(_ ast.SemanticParser) (string, io.Error) {
	l := a.Location()
	return fmt.Sprintf("%s:%s", l.String(), a.Name_.Value), nil
}

func (a *Abstract) Concretize(mapping []ast.Type) ast.Type {
	return a
}
