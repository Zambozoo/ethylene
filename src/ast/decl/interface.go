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

type Interface struct {
	BaseDecl
	GenericDecl

	Implements data.Set[ast.DeclType] // Interfaces this decl implements
}

func newInterface() *Interface {
	return &Interface{
		BaseDecl:    newDecl(),
		GenericDecl: NewGenericDecl(),
	}
}

func (i *Interface) SetTailed() io.Error {
	return io.NewError("interfaces cannot be tailed", zap.Any("location", i.Name_.Location()))
}

func (a *Interface) String() string {
	return fmt.Sprintf("Interface{Name: %s%s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		a.Name().Value,
		a.TypesMap,
		maps.Keys(a.Implements),
		strings.Join(maps.Keys(a.Methods_), ","),
		strings.Join(maps.Keys(a.Members_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
		strings.Join(maps.Keys(a.StaticMembers_), ","),
	)
}

func (i *Interface) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if i.BaseDecl.StartToken, err = p.Consume(token.TOK_INTERFACE); err != nil {
		return err
	}

	if _, err := p.ParseDeclType(i); err != nil {
		return err
	}

	if p.Match(token.TOK_SUBTYPE) {
		if i.Implements, err = p.ParseParentTypes(); err != nil {
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
			if _, exists := i.TypesMap[f.Name().Value]; exists {
				return io.NewError("inner decl name duplicates generic type",
					zap.Any("decl", f.Name()),
					zap.Any("location", f.Location()),
				)
			}
		} else if !f.HasModifier(ast.MOD_STATIC) && !f.HasModifier(ast.MOD_VIRTUAL) {
			return io.NewError("only static and virual fields are allowed in interfaces",
				zap.Any("field", f.Name()),
				zap.Any("location", f.Location()),
			)
		}
		if err := i.AddField(f); err != nil {
			return err
		}
	}
	i.BaseDecl.EndToken = p.Prev()

	return nil
}

func (i *Interface) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
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

	for _, parent := range i.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return nil, err
		}

		if _, isInterface := parentDecl.(*Interface); !isInterface {
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
			i.Implements.Set(parent)
		}

		if err := i.BaseDecl.Extends(p, parentDecl, visitedDecls); err != nil {
			return nil, err
		}
	}
	visitedDecls.Set(i)

	return i.Implements, i.BaseDecl.LinkParents(p, visitedDecls)
}

func (i *Interface) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(i); exists {
		return nil
	}

	for _, parent := range i.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return err
		}

		if err := parentDecl.LinkFields(p, visitedDecls); err != nil {
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
