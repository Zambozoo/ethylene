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

type Struct struct {
	BaseDecl

	IsTailed        bool
	GenericTypeArgs map[string]ast.GenericTypeArg // Generic type parameters
}

func newStruct() *Struct {
	return &Struct{
		BaseDecl:        newDecl(),
		GenericTypeArgs: map[string]ast.GenericTypeArg{},
	}
}

func (s *Struct) SetGenerics(genericTypeArgs map[string]ast.GenericTypeArg) io.Error {
	s.GenericTypeArgs = genericTypeArgs
	return nil
}

func (s *Struct) SetTailed() io.Error {
	s.IsTailed = true
	return nil
}

func (s *Struct) Generics() map[string]ast.GenericTypeArg {
	return s.GenericTypeArgs
}

func (s *Struct) String() string {
	return fmt.Sprintf("Struct{Name: %s%s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		s.Name().Value,
		s.GenericTypeArgs,
		strings.Join(maps.Keys(s.Methods_), ","),
		strings.Join(maps.Keys(s.Members_), ","),
		strings.Join(maps.Keys(s.StaticMembers_), ","),
		strings.Join(maps.Keys(s.StaticMembers_), ","),
	)
}

func (s *Struct) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if s.BaseDecl.StartToken, err = p.Consume(token.TOK_STRUCT); err != nil {
		return err
	}

	if _, err := p.ParseDeclType(s); err != nil {
		return err
	}

	if _, err := p.Consume(token.TOK_LEFTBRACE); err != nil {
		return err
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		f, err := p.ParseField()
		if err != nil {
			return err
		} else if _, ok := f.(ast.DeclField); ok {
			if _, exists := s.GenericTypeArgs[f.Name().Value]; exists {
				return io.NewError("inner decl name duplicates generic type",
					zap.Any("decl", f.Name()),
					zap.Any("location", f.Location()),
				)
			}
		} else if f.HasModifier(ast.MOD_VIRTUAL) {
			return io.NewError("virtual fields are not allowed in structs",
				zap.Any("field", f.Name()),
				zap.Any("location", f.Location()),
			)
		}
		s.AddField(f)
	}
	s.BaseDecl.EndToken = p.Prev()

	return nil
}

func (s *Struct) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) io.Error {
	if _, exists := visitedDecls.Get(s); exists {
		return nil
	}
	defer visitedDecls.Set(s)

	return s.BaseDecl.LinkParents(p, visitedDecls, cycleMap)
}
func (s *Struct) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Handle generic constraints
	return s.BaseDecl.Semantic(p)
}
