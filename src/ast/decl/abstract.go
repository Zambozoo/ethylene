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

	SuperClass ast.DeclType   // Optional
	Implements []ast.DeclType // Interfaces this class implements
}

func (a *Abstract) SetTailed() io.Error {
	a.IsTailed = true
	return nil
}

func newAbstract() *Abstract {
	return &Abstract{
		BaseDecl:    newDecl(),
		GenericDecl: newGenericDecl(),
	}
}

func (a *Abstract) String() string {
	return fmt.Sprintf("Abstract{Name: %s%s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		a.Name().Value,
		a.TypesMap,
		parentsString(a.SuperClass, a.Implements),
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
		a.AddField(f)
	}
	a.BaseDecl.EndToken = p.Prev()

	return nil
}

func (a *Abstract) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) io.Error {
	if _, exists := visitedDecls.Get(a); exists {
		return nil
	}
	defer visitedDecls.Set(a)

	l := a.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return io.NewError("cyclical inheritance",
			zap.Any("abstract", a.Name()),
			zap.Any("location", l),
		)
	}
	cycleMap[l.String()] = struct{}{}
	defer delete(cycleMap, l.String())

	for _, parent := range a.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return err
		}
		_, isAbstract := parentDecl.(*Abstract)
		_, isClass := parentDecl.(*Class)
		_, isStruct := parentDecl.(*Struct)
		_, isEnum := parentDecl.(*Enum)
		if isAbstract {
			if a.SuperClass != nil {
				return io.NewError("abstracts cannot implement multiple concrete parents",
					zap.Any("abstract", a.Name()),
					zap.Any("location", a.Location()),
					zap.Any("parent", parentDecl.Name()),
				)
			}
			a.SuperClass = parent
		} else if isClass {
			return io.NewError("abstracts cannot implement classes",
				zap.Any("class", parentDecl.Name()),
				zap.Any("abstract", a.Name()),
				zap.Any("location", a.Location()),
			)
		} else if isStruct || isEnum {
			return io.NewError("abstract cannot implement struct or enum parents",
				zap.Any("class", a.Name()),
				zap.Any("location", a.Location()),
				zap.Any("parent", parentDecl.Name()),
			)
		}

		if err := parentDecl.LinkParents(p, visitedDecls, cycleMap); err != nil {
			return err
		}
	}

	return a.BaseDecl.LinkParents(p, visitedDecls, cycleMap)
}

func (a *Abstract) LinkMethods(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	for _, parent := range a.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return err
		}
		if err := a.BaseDecl.Extends(p, parentDecl, visitedDecls); err != nil {
			return err
		}
	}

	return nil
}

func (a *Abstract) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return a.BaseDecl.Semantic(p)
}
