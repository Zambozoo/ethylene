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

	IsTailed bool

	SuperClass ast.DeclType           // Optional
	Parents_   data.Set[ast.DeclType] // Interfaces this decl implements
}

func (c *Class) Parents() data.Set[ast.DeclType] {
	return c.Parents_
}

func (*Class) IsInterface() bool {
	return false
}
func (*Class) IsAbstract() bool {
	return false
}
func (*Class) IsClass() bool {
	return true
}

func (c *Class) IsConstant() bool {
	return false
}

func (c *Class) SetConstant() {}

func newClass() *Class {
	decl := newDecl()
	decl.IsClass = true
	return &Class{
		BaseDecl: decl,
	}
}

func (c *Class) String() string {
	var parentsString string
	if len(c.Parents_) > 0 {
		parentsString = "<: [" + strings.Join(maps.Keys(c.Parents_), ",") + "]"
	}
	return fmt.Sprintf("class %s%s {\n%s\n%s\n%s\n%s}",
		c.Name().Value,
		parentsString,
		strings.Join(maps.Keys(c.Methods_), "\n"),
		strings.Join(maps.Keys(c.Members_), "\n"),
		strings.Join(maps.Keys(c.StaticMembers_), "\n"),
		strings.Join(maps.Keys(c.StaticMembers_), "\n"),
	)
}

func (c *Class) Syntax(p ast.SyntaxParser) (ast.Declaration, io.Error) {
	var err io.Error
	if c.BaseDecl.StartToken, err = p.Consume(token.TOK_CLASS); err != nil {
		return nil, err
	}

	if c.Name_, err = p.Consume(token.TOK_IDENTIFIER); err != nil {
		return nil, err
	}

	genericDecl, err := generics.Syntax(c, p)
	if err != nil {
		return nil, err
	} else if genericDecl != nil {
		p.UnwrapScope()
		p.WrapScope(genericDecl)
	}

	if p.Match(token.TOK_TILDE) {
		c.IsTailed = true
	}

	if p.Match(token.TOK_SUBTYPE) {
		if c.Parents_, err = p.ParseParentTypes(); err != nil {
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
						zap.Stringer("decl", f.Name()),
						zap.Stringer("location", f.Location()),
					)
				}
			}
		} else if f.HasModifier(ast.MOD_VIRTUAL) {
			return nil, io.NewError("virtual fields are not allowed in classes", zap.Stringer("field", f.Name()))
		}

		if err := c.AddField(f); err != nil {
			return nil, err
		}
	}
	c.BaseDecl.EndToken = p.Prev()

	c.BaseDecl.Index, err = p.Types().NextClassIndex(c)
	if err != nil {
		return nil, err
	}

	if genericDecl != nil {
		return genericDecl, nil
	}
	return c, nil
}

func (c *Class) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	return c.LinkParentsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return c.Parents_, nil
		},
		func(parent ast.DeclType) io.Error {
			c.Parents_.Set(parent)
			return nil
		},
		visitedDecls, cycleMap)
}

func (c *Class) LinkParentsWithProvider(
	p ast.SemanticParser,
	parentProviderFunc generics.ParentProviderFunc,
	parentSetterFunc generics.ParentSetterFunc,
	visitedDecls *data.AsyncSet[ast.Declaration],
	cycleMap map[string]struct{},
) (data.Set[ast.DeclType], io.Error) {
	if _, exists := visitedDecls.Get(c); exists {
		return c.Parents_, nil
	}

	l := c.Location()
	if _, isCyclical := cycleMap[l.String()]; isCyclical {
		return nil, io.NewError("cyclical inheritance",
			zap.Stringer("class", c.Name()),
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

		if parentDecl.IsClass() || parentDecl.IsAbstract() {
			if c.SuperClass != nil {
				return nil, io.NewError("class cannot implement multiple concrete parents",
					zap.Stringer("class", c.Name()),
					zap.Stringer("location", c.Location()),
					zap.Stringer("parent", parentDecl.Name()),
				)
			}
			c.SuperClass = parent
		} else if !parentDecl.IsInterface() {
			return nil, io.NewError("class cannot implement struct or enum parents",
				zap.Stringer("class", c.Name()),
				zap.Stringer("location", c.Location()),
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
	visitedDecls.Set(c)
	ps, err = parentProviderFunc()
	if err != nil {
		return nil, err
	}
	return ps, c.BaseDecl.LinkParents(p, visitedDecls)
}

func (c *Class) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return c.LinkFieldsWithProvider(p,
		func() (data.Set[ast.DeclType], io.Error) {
			return c.Parents_, nil
		},
		func() map[string]ast.Method {
			return c.Methods_
		},
		visitedDecls,
	)
}

func (c *Class) LinkFieldsWithProvider(p ast.SemanticParser, parentProviderFunc generics.ParentProviderFunc, methodProviderFunc generics.MethodProviderFunc, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(c); exists {
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
		if err := c.BaseDecl.ExtendsParent(p, methodProviderFunc, parentDecl, visitedDecls); err != nil {
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

func (c *Class) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return c.Equals(p, parent)
}

func (c *Class) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	for _, p := range c.Parents_ {
		if equals, err := p.Equals(parser, parent); equals || err != nil {
			return equals, err
		}
	}

	return c.Equals(parser, parent)
}

func (c *Class) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherClass, ok := other.(*Class); ok {
		return c == otherClass, nil
	}
	if otherDeclType, ok := other.(ast.DeclType); ok {
		otherDeclaration, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		}
		if otherClass, ok := otherDeclaration.(*Class); ok {
			return c == otherClass, nil
		}
	}

	return false, nil
}

func (c *Class) Concretize(mapping []ast.Type) ast.Type {
	return c
}

func (c *Class) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	return typeid.NewTypeID(parser.Types().ClassIndex(c.Index), 0), nil
}
