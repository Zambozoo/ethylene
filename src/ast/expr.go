package ast

import "geth-cody/io"

// Expression represents an expression in the AST.
type Expression interface {
	// Inherits Node interface.
	Node
	// Semantic performs semantic analysis on the expression, returning its type.
	Semantic(p SemanticParser) (Type, io.Error)
}
