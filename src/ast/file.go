package ast

import (
	"geth-cody/compile/data"
	"geth-cody/io"
	"geth-cody/io/path"
)

// Import represents a file import.
type Import interface {
	// Path returns the Import's associated path
	Path() path.Path
}

// File represents a file node in the AST.
type File interface {
	// Inherits Node interface.
	Node

	// GetImport returns the associated import for a dependency name
	GetImport(name string) (Import, bool)
	// Declaration returns the file's main declaration
	Declaration() Declaration

	// Syntax parses the file syntactically
	Syntax(p SyntaxParser) io.Error

	// LinkParents links parent nodes within the AST, facilitating inheritance.
	LinkParents(p SemanticParser, visitedDecls *data.AsyncSet[Declaration]) io.Error
	// LinkParents links parent and child field nodes within the AST, facilitating inheritance.
	LinkFields(p SemanticParser, visitedDecls *data.AsyncSet[Declaration]) io.Error
	// Semantic performs semantic analysis on the file node.
	Semantic(p SemanticParser) io.Error
}
