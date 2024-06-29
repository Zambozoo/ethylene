package ast

import "geth-cody/io"

// Statement represents a statement in the AST.
type Statement interface {
	// Inherits Node interface.
	Node
	// Semantic performs semantic analysis on the statement, returning its return type.
	Semantic(p SemanticParser) (Type, io.Error)
}
