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
	GenericDecl

	IsTailed bool

	SuperClass ast.DeclType           // Optional
	Implements data.Set[ast.DeclType] // Interfaces this decl implements
}

func newClass() *Class {
	decl := newDecl()
	decl.IsClass = true
	return &Class{
		BaseDecl:    decl,
		GenericDecl: NewGenericDecl(),
	}
}

func (c *Class) SetTailed() io.Error {
	c.IsTailed = true
	return nil
}

func (c *Class) String() string {
	return fmt.Sprintf("Class{Name: %s%s, Parents: %s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		c.Name().Value,
		c.TypesMap,
		maps.Keys(c.Implements),
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

	if _, err := p.ParseDeclType(c); err != nil {
		return err
	}

	if p.Match(token.TOK_SUBTYPE) {
		if c.Implements, err = p.ParseParentTypes(); err != nil {
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
			if _, exists := c.TypesMap[f.Name().Value]; exists {
				return io.NewError("inner decl name duplicates generic type",
					zap.Any("decl", f.Name()),
					zap.Any("location", f.Location()),
				)
			}
		} else if f.HasModifier(ast.MOD_VIRTUAL) {
			return io.NewError("virtual fields are not allowed in classes", zap.Any("field", f.Name()))
		}

		if err := c.AddField(f); err != nil {
			return err
		}
	}
	c.BaseDecl.EndToken = p.Prev()

	return nil
}

func (c *Class) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	if _, exists := visitedDecls.Get(c); exists {
		return c.Implements, nil
	}

	l := c.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return nil, io.NewError("cyclical inheritance",
			zap.Any("class", c.Name()),
			zap.Any("location", l),
		)
	}
	cycleMap[l.String()] = struct{}{}
	defer delete(cycleMap, l.String())

	for _, parent := range c.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return nil, err
		}
		_, isAbstract := parentDecl.(*Abstract)
		_, isClass := parentDecl.(*Class)
		_, isStruct := parentDecl.(*Struct)
		_, isEnum := parentDecl.(*Enum)
		if isAbstract || isClass {
			if c.SuperClass != nil {
				return nil, io.NewError("class cannot implement multiple concrete parents",
					zap.Any("class", c.Name()),
					zap.Any("location", c.Location()),
					zap.Any("parent", parentDecl.Name()),
				)
			}
			c.SuperClass = parent
		} else if isStruct || isEnum {
			return nil, io.NewError("class cannot implement struct or enum parents",
				zap.Any("class", c.Name()),
				zap.Any("location", c.Location()),
				zap.Any("parent", parentDecl.Name()),
			)
		}

		parents, err := parentDecl.LinkParents(p, visitedDecls, cycleMap)
		if err != nil {
			return nil, err
		}
		for _, parent := range parents {
			c.Implements.Set(parent)
		}
	}
	visitedDecls.Set(c)

	return c.Implements, c.BaseDecl.LinkParents(p, visitedDecls)
}

func (c *Class) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(c); exists {
		return nil
	}

	for _, parent := range c.Implements {
		parentDecl, err := parent.Declaration()
		if err != nil {
			return err
		}

		if err := parentDecl.LinkFields(p, visitedDecls); err != nil {
			return err
		}
	}

	if c.SuperClass != nil {
		for name, m := range c.Members() {
			if _, ok := c.Members_[name]; !ok && !m.HasModifier(ast.MOD_PRIVATE) {
				c.Members_[name] = m
			}
		}
	}
	visitedDecls.Set(c)

	return c.BaseDecl.LinkFields(p, visitedDecls)
}

func (c *Class) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return c.BaseDecl.Semantic(p)
}
