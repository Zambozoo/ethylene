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

	IsTailed           bool
	GenericConstraints map[string]GenericConstraint // Generic type parameters
}

func newStruct() *Struct {
	return &Struct{
		BaseDecl:           newDecl(),
		GenericConstraints: map[string]GenericConstraint{},
	}
}

func (s *Struct) String() string {
	return fmt.Sprintf("Interface{Name: %s%s, Members: %s, Methods: %s, StaticMembers: %s, StaticMethods: %s}",
		s.Name().Value,
		s.GenericConstraints,
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

	s.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	s.GenericConstraints, err = syntaxGenericConstraints(p)
	if err != nil {
		return err
	}

	if p.Match(token.TOK_TILDE) {
		s.IsTailed = true
	}

	if _, err := p.Consume(token.TOK_LEFTBRACE); err != nil {
		return err
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		f, err := p.ParseField()
		if err != nil {
			return err
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
