package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Abstract struct {
	BaseDecl
	GenericDecl

	IsTailed bool

	SuperClass ast.DeclType           // Optional
	Implements data.Set[ast.DeclType] // Interfaces this decl implements
}

func (a *Abstract) SetTailed() io.Error {
	a.IsTailed = true
	return nil
}

func newAbstract() *Abstract {
	return &Abstract{
		BaseDecl:    newDecl(),
		GenericDecl: NewGenericDecl(),
	}
}

func (a *Abstract) String() string {
	return fmt.Sprintf("Abstract{Name: %s%s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		a.Name().Value,
		a.TypesMap,
		maps.Keys(a.Implements),
		strings.Join(maps.Keys(a.Methods_), ","),
		strings.Join(maps.Keys(a.Members_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
	)
}

func (a *Abstract) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if a.BaseDecl.StartToken, err = p.Consume(token.TOK_ABSTRACT); err != nil {
		return err
	}

	if _, err := p.ParseDeclType(a); err != nil {
		return err
	}

	if p.Match(token.TOK_SUBTYPE) {
		if a.Implements, err = p.ParseParentTypes(); err != nil {
			return err
		}
	}

	if _, err := p.Consume(token.TOK_LEFTBRACE); err != nil {
		return err
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		f, err := p.ParseField()
		if err != nil {
			return err
		} else if _, ok := f.(ast.DeclField); ok {
			if _, exists := a.TypesMap[f.Name().Value]; exists {
				return io.NewError("inner decl name duplicates generic type",
					zap.Any("decl", f.Name()),
					zap.Any("location", f.Location()),
				)
			}
		}
		if err := a.AddField(f); err != nil {
			return err
		}
	}
	a.BaseDecl.EndToken = p.Prev()

	return nil
}

func (a *Abstract) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
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

	for _, parent := range a.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return nil, err
		}
		_, isAbstract := parentDecl.(*Abstract)
		_, isClass := parentDecl.(*Class)
		_, isStruct := parentDecl.(*Struct)
		_, isEnum := parentDecl.(*Enum)
		if isAbstract {
			if a.SuperClass != nil {
				return nil, io.NewError("abstracts cannot implement multiple concrete parents",
					zap.Any("abstract", a.Name()),
					zap.Any("location", a.Location()),
					zap.Any("parent", parentDecl.Name()),
				)
			}
			a.SuperClass = parent
		} else if isClass {
			return nil, io.NewError("abstracts cannot implement classes",
				zap.Any("class", parentDecl.Name()),
				zap.Any("abstract", a.Name()),
				zap.Any("location", a.Location()),
			)
		} else if isStruct || isEnum {
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
			a.Implements.Set(parent)
		}
	}
	visitedDecls.Set(a)

	return a.Implements, a.BaseDecl.LinkParents(p, visitedDecls)
}

func (a *Abstract) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(a); exists {
		return nil
	}

	for _, parent := range a.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return err
		}

		if err := parentDecl.LinkFields(p, visitedDecls); err != nil {
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
