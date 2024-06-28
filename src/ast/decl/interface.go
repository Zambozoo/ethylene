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

type Interface struct {
	BaseDecl

	Implements data.Set[ast.DeclType] // Interfaces this decl implements
}

func (i *Interface) Parents() data.Set[ast.DeclType] {
	return i.Implements
}

func (*Interface) IsInterface() bool {
	return true
}
func (*Interface) IsAbstract() bool {
	return false
}
func (*Interface) IsClass() bool {
	return false
}

func (i *Interface) IsConstant() bool {
	return false
}

func (i *Interface) SetConstant() {}

func newInterface() *Interface {
	return &Interface{
		BaseDecl: newDecl(),
	}
}

func (i *Interface) SetTailed() io.Error {
	return io.NewError("interfaces cannot be tailed", zap.Any("location", i.Name_.Location()))
}

func (a *Interface) String() string {
	return fmt.Sprintf("Interface{Name: %s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		a.Name().Value,
		maps.Keys(a.Implements),
		strings.Join(maps.Keys(a.Methods_), ","),
		strings.Join(maps.Keys(a.Members_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
	)
}

func (i *Interface) Syntax(p ast.SyntaxParser) (ast.Declaration, io.Error) {
	var err io.Error
	if i.BaseDecl.StartToken, err = p.Consume(token.TOK_INTERFACE); err != nil {
		return nil, err
	}

	if i.Name_, err = p.Consume(token.TOK_IDENTIFIER); err != nil {
		return nil, err
	}

	genericDecl, err := generics.Syntax(i, p)
	if err != nil {
		return nil, err
	} else if genericDecl != nil {
		p.UnwrapScope()
		p.WrapScope(genericDecl)
	}

	if p.Match(token.TOK_SUBTYPE) {
		if i.Implements, err = p.ParseParentTypes(); err != nil {
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
			if genericDecl != nil {
				if _, ok := genericDecl.GenericParamIndex(f.Name().Value); ok {
					return nil, io.NewError("inner decl name duplicates generic type",
						zap.Any("decl", f.Name()),
						zap.Any("location", f.Location()),
					)
				}
			}
		} else if !f.HasModifier(ast.MOD_STATIC) && !f.HasModifier(ast.MOD_VIRTUAL) {
			return nil, io.NewError("only static and virtual fields are allowed in interfaces",
				zap.Any("field", f.Name()),
				zap.Any("location", f.Location()),
			)
		}
		if err := i.AddField(f); err != nil {
			return nil, err
		}
	}
	i.BaseDecl.EndToken = p.Prev()

	if genericDecl != nil {
		return genericDecl, nil
	}
	return i, nil
}

func (i *Interface) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	return i.LinkParentsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return i.Implements, nil
		},
		func(parent ast.DeclType) io.Error {
			i.Implements.Set(parent)
			return nil
		},
		visitedDecls, cycleMap)
}

func (i *Interface) LinkParentsWithProvider(
	p ast.SemanticParser,
	parentProviderFunc generics.ParentProviderFunc,
	parentSetterFunc generics.ParentSetterFunc,
	visitedDecls *data.AsyncSet[ast.Declaration],
	cycleMap map[string]struct{},
) (data.Set[ast.DeclType], io.Error) {
	if _, exists := visitedDecls.Get(i); exists {
		return i.Implements, nil
	}

	l := i.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return nil, io.NewError("cyclical inheritance",
			zap.Any("interface", i.Name()),
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

		if !parentDecl.IsInterface() {
			return nil, io.NewError("interface can only implement interface parents",
				zap.Any("interface", i.Name()),
				zap.Any("location", i.Location()),
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
	visitedDecls.Set(i)
	ps, err = parentProviderFunc()
	if err != nil {
		return nil, err
	}
	return ps, i.BaseDecl.LinkParents(p, visitedDecls)
}

func (i *Interface) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return i.LinkFieldsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return i.Implements, nil
		},
		func() map[string]ast.Method {
			return i.Methods_
		},
		visitedDecls,
	)
}

func (i *Interface) LinkFieldsWithProvider(p ast.SemanticParser, parentProviderFunc generics.ParentProviderFunc, methodProviderFunc generics.MethodProviderFunc, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(i); exists {
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
		if err := i.BaseDecl.ExtendsParent(p, methodProviderFunc, parentDecl, visitedDecls); err != nil {
			return err
		}
	}
	visitedDecls.Set(i)

	return i.BaseDecl.LinkFields(p, visitedDecls)
}

func (i *Interface) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return i.BaseDecl.Semantic(p)
}

func (i *Interface) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return i.Equals(p, parent)
}

func (i *Interface) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	for _, p := range i.Implements {
		if equals, err := p.Equals(parser, parent); equals || err != nil {
			return equals, err
		}
	}

	return i.Equals(parser, parent)
}

func (i *Interface) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherInterace, ok := other.(*Interface); ok {
		return i == otherInterace, nil
	} else if otherDeclType, ok := other.(ast.DeclType); ok {
		otherDeclaration, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		}
		if otherInterace, ok := otherDeclaration.(*Interface); ok {
			return i == otherInterace, nil
		}
	}

	return false, nil
}

func (i *Interface) Key(p ast.SemanticParser) (string, io.Error) {
	l := i.Location()
	return fmt.Sprintf("%s:%s", l.String(), i.Name_.Value), nil
}

func (i *Interface) Concretize(mapping []ast.Type) ast.Type {
	return i
}
