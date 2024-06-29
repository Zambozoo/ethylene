package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/io"
)

type SubType struct {
	Binary
}

func (s *SubType) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := s.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := s.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := s.MustBothExtendOne(p, left, right, &type_.TypeID{}); err != nil {
		return nil, err
	}

	return &type_.TypeID{}, nil
}
func (s *SubType) String() string {
	return fmt.Sprintf("%s <: %s", s.Left.String(), s.Right.String())
}

type SuperType struct {
	Binary
}

func (s *SuperType) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	left, err := s.Left.Semantic(p)
	if err != nil {
		return nil, err
	}

	right, err := s.Right.Semantic(p)
	if err != nil {
		return nil, err
	}

	if _, err := s.MustBothExtendOne(p, left, right, &type_.TypeID{}); err != nil {
		return nil, err
	}

	return &type_.TypeID{}, nil
}
func (s *SuperType) String() string {
	return fmt.Sprintf("%s :> %s", s.Left.String(), s.Right.String())
}
