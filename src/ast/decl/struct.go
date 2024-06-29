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

type Struct struct {
	BaseDecl

	IsTailed bool
}

func (*Struct) IsInterface() bool {
	return false
}
func (*Struct) IsAbstract() bool {
	return false
}
func (*Struct) IsClass() bool {
	return false
}

func (*Struct) IsConstant() bool {
	return false
}

func (*Struct) SetConstant() {}

func newStruct() *Struct {
	return &Struct{
		BaseDecl: newDecl(),
	}
}

func (s *Struct) String() string {
	return fmt.Sprintf("struct %s {\n%s\n%s\n%s\n%s}",
		s.Name().Value,
		strings.Join(maps.Keys(s.Methods_), "\n"),
		strings.Join(maps.Keys(s.Members_), "\n"),
		strings.Join(maps.Keys(s.StaticMembers_), "\n"),
		strings.Join(maps.Keys(s.StaticMembers_), "\n"),
	)
}

func (s *Struct) Syntax(p ast.SyntaxParser) (ast.Declaration, io.Error) {
	var err io.Error
	if s.BaseDecl.StartToken, err = p.Consume(token.TOK_STRUCT); err != nil {
		return nil, err
	}

	if s.Name_, err = p.Consume(token.TOK_IDENTIFIER); err != nil {
		return nil, err
	}
	genericDecl, err := generics.Syntax(s, p)
	if err != nil {
		return nil, err
	} else if genericDecl != nil {
		p.UnwrapScope()
		p.WrapScope(genericDecl)
	}

	if p.Match(token.TOK_TILDE) {
		s.IsTailed = true
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
			return nil, io.NewError("virtual fields are not allowed in structs",
				zap.Stringer("field", f.Name()),
				zap.Stringer("location", f.Location()),
			)
		}
		if err := s.AddField(f); err != nil {
			return nil, err
		}
	}
	s.BaseDecl.EndToken = p.Prev()

	if genericDecl != nil {
		return genericDecl, nil
	}
	return s, nil
}

func (s *Struct) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], _ map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	if _, exists := visitedDecls.Get(s); exists {
		return nil, nil
	}
	visitedDecls.Set(s)

	return nil, s.BaseDecl.LinkParents(p, visitedDecls)
}

func (s *Struct) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	if _, exists := visitedDecls.Get(s); exists {
		return nil
	}
	defer visitedDecls.Set(s)

	return s.BaseDecl.LinkFields(p, visitedDecls)
}

func (s *Struct) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return s.BaseDecl.Semantic(p)
}

func (s *Struct) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return s.Equals(p, parent)
}

func (s *Struct) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return s.Equals(p, parent)
}

func (s *Struct) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherStruct, ok := other.(*Struct); ok {
		return s == otherStruct, nil
	} else if otherDeclType, ok := other.(ast.DeclType); ok {
		otherDeclaration, err := otherDeclType.Declaration(p)
		if err != nil {
			return false, err
		} else if otherStruct, ok := otherDeclaration.(*Struct); ok {
			return s == otherStruct, nil
		}
	}

	return false, nil
}

func (s *Struct) Concretize(mapping []ast.Type) ast.Type {
	return s
}
