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

// Class represents a declaration of the form
//
//	IDENTIFIER_LIST := IDENTIFIER (`,` IDENTIFIER)*
//	GENERIC := IDENTIFIER (`<:` `[` IDENTIFIER_LIST `]`)?
//	GENERIC_LIST := GENERIC (`,` GENRERIC)*
//
//	PARENT := IDENTIFIER (`[` IDENTIFIER_LIST `]`)?
//	PARENT_LIST := PARENT (`,` PARENT)*
//
//	`class` IDENTIFIER `~`? (`[` GENERIC_LIST `]`)? `[` PARENT_LIST `]` `{` FIELD* `}`
type Class struct {
	BaseDecl

	IsTailed           bool
	GenericConstraints map[string]ast.GenericConstraint // Generic type parameters

	SuperClass ast.DeclType   // Optional
	Implements []ast.DeclType // Interfaces this class implements
}

func newClass() *Class {
	decl := newDecl()
	decl.IsClass = true
	return &Class{
		BaseDecl:           decl,
		GenericConstraints: map[string]ast.GenericConstraint{},
	}
}

func (c *Class) Generics() map[string]ast.GenericConstraint {
	return c.GenericConstraints
}

func (c *Class) String() string {
	return fmt.Sprintf("Class{Name: %s%s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		c.Name().Value,
		c.GenericConstraints,
		parentsString(c.SuperClass, c.Implements),
		strings.Join(maps.Keys(c.Methods_), ","),
		strings.Join(maps.Keys(c.Members_), ","),
		strings.Join(maps.Keys(c.StaticMembers_), ","),
		strings.Join(maps.Keys(c.StaticMembers_), ","),
	)
}

func (c *Class) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if c.BaseDecl.StartToken, err = p.Consume(token.TOK_CLASS); err != nil {
		return err
	}

	c.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	c.GenericConstraints, err = syntaxGenericConstraints(p)
	if err != nil {
		return err
	}

	if p.Match(token.TOK_TILDE) {
		c.IsTailed = true
	}

	if p.Match(token.TOK_SUBTYPE) {
		c.Implements, err = syntaxDeclTypes(p)
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
		} else if f.HasModifier(ast.MOD_VIRTUAL) {
			return io.NewError("virtual fields are not allowed in classes", zap.Any("field", f.Name()))
		}

		c.AddField(f)
	}
	c.BaseDecl.EndToken = p.Prev()

	return nil
}

func (c *Class) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) io.Error {
	if _, exists := visitedDecls.Get(c); exists {
		return nil
	}
	defer visitedDecls.Set(c)

	l := c.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return io.NewError("cyclical inheritance",
			zap.Any("class", c.Name()),
			zap.Any("location", l),
		)
	}
	cycleMap[l.String()] = struct{}{}
	defer delete(cycleMap, l.String())

	for _, parent := range c.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return err
		}
		_, isAbstract := parentDecl.(*Abstract)
		_, isClass := parentDecl.(*Class)
		_, isStruct := parentDecl.(*Struct)
		_, isEnum := parentDecl.(*Enum)
		if isAbstract || isClass {
			if c.SuperClass != nil {
				return io.NewError("class cannot implement multiple concrete parents",
					zap.Any("class", c.Name()),
					zap.Any("location", c.Location()),
					zap.Any("parent", parentDecl.Name()),
				)
			}
			c.SuperClass = parent
		} else if isStruct || isEnum {
			return io.NewError("class cannot implement struct or enum parents",
				zap.Any("class", c.Name()),
				zap.Any("location", c.Location()),
				zap.Any("parent", parentDecl.Name()),
			)
		}

		if err := parentDecl.LinkParents(p, visitedDecls, cycleMap); err != nil {
			return err
		}

		if err := c.BaseDecl.Extends(p, parentDecl, visitedDecls); err != nil {
			return err
		}
	}

	return c.BaseDecl.LinkParents(p, visitedDecls, cycleMap)
}

func (c *Class) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return c.BaseDecl.Semantic(p)
}
