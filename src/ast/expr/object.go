package expr

import (
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// This represents expressions of the form
//
//	'this'
type This struct{ token.Token }

func (t *This) String() string {
	return "this"
}

func (t *This) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	d := p.Scope().Declaration()
	if d == nil {
		return nil, io.NewError("this expression in static field",
			zap.Stringer("location", t.Location()),
		)
	}

	return &type_.Pointer{
		Type: d,
	}, nil
}

// Super represents expressions of the form
//
//	'super'
type Super struct{ token.Token }

func (s *Super) String() string {
	return "super"
}

func (s *Super) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	d := p.Scope().Declaration()
	if d == nil {
		return nil, io.NewError("super expression in static field",
			zap.Stringer("location", s.Location()),
		)
	}

	cd, ok := d.(ast.ChildDeclaration)
	if !ok {
		return nil, io.NewError("super expression in non-superable declaration",
			zap.Stringer("location", s.Location()),
		)
	}
	st, ok := cd.Super()
	if !ok {
		return nil, io.NewError("super expression in non-superable declaration",
			zap.Stringer("location", s.Location()),
		)
	}

	return &type_.Pointer{
		Type: st,
	}, nil
}
