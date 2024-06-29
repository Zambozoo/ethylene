package expr

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// This represents expressions of the form
//
//	'this'
type This struct{ token.Token }

func (t *This) String() string {
	return "this"
}

func (t *This) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// Super represents expressions of the form
//
//	'super'
type Super struct{ token.Token }

func (s *Super) String() string {
	return "super"
}

func (s *Super) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}
