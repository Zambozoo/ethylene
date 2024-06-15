package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Interface struct {
	BaseDecl

	GenericConstraints map[string]GenericConstraint // Generic type parameters

	Implements []ast.DeclType // Interfaces this class implements
}

func newInterface() *Interface {
	return &Interface{
		BaseDecl:           newDecl(),
		GenericConstraints: map[string]GenericConstraint{},
	}
}

func (a *Interface) String() string {
	return fmt.Sprintf("Interface{Name: %s%s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		a.Name().Value,
		a.GenericConstraints,
		strs.Strings(a.Implements),
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

	i.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	i.GenericConstraints, err = syntaxGenericConstraints(p)
	if err != nil {
		return err
	}

	if p.Match(token.TOK_SUBTYPE) {
		i.Implements, err = syntaxDeclTypes(p)
		if err != nil {
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
		} else if !f.HasModifier(ast.MOD_STATIC) && !f.HasModifier(ast.MOD_VIRTUAL) {
			return io.NewError("only static and virual fields are allowed in interfaces",
				zap.Any("field", f.Name()),
				zap.Any("location", f.Location()),
			)
		}
		i.AddField(f)
	}
	i.BaseDecl.EndToken = p.Prev()

	return nil
}

func (i *Interface) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) io.Error {
	if _, exists := visitedDecls.Get(i); exists {
		return nil
	}
	defer visitedDecls.Set(i)

	l := i.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return io.NewError("cyclical inheritance",
			zap.Any("interface", i.Name()),
			zap.Any("location", l),
		)
	}
	cycleMap[l.String()] = struct{}{}
	defer delete(cycleMap, l.String())

	for _, parent := range i.Implements {
		parentDecl, err := parent.Declaration(p.TypeContext())
		if err != nil {
			return err
		}

		if _, isInterface := parentDecl.(*Interface); !isInterface {
			return io.NewError("interface can only implement interface parents",
				zap.Any("interface", i.Name()),
				zap.Any("location", i.Location()),
				zap.Any("parent", parentDecl.Name()),
			)
		}

		if err := parentDecl.LinkParents(p, visitedDecls, cycleMap); err != nil {
			return err
		}

		if err := i.BaseDecl.Extends(p, parentDecl, visitedDecls); err != nil {
			return err
		}
	}

	return i.BaseDecl.LinkParents(p, visitedDecls, cycleMap)
}

func (i *Interface) Semantic(p ast.SemanticParser) io.Error {
	p.WrapTypeContext(i)
	defer p.UnwrapTypeContext()

	// TODO: Handle generic constraints
	return i.BaseDecl.Semantic(p)
}
