package ast

import (
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Field represents a field within the AST.
type Field interface {
	// Inherits Node interface.
	Node
	// Semantic performs semantic analysis on the field.
	Semantic(p SemanticParser) io.Error

	// Name returns the token representing the field's name.
	Name() *token.Token
	// HasModifier checks if the field has a specific modifier.
	HasModifier(Modifier) bool
}

// Method represents a method, which is a function associated with a type.
type Method interface {
	// Inherits from Field interface.
	Field
	// Type returns the type of the method.
	Type() Type
	// ReturnType returns the return type of the method.
	ReturnType() Type
}

// Member represents a member variable within a type.
type Member interface {
	// Inherits from Field interface.
	Field
	// Type returns the type of the member.
	Type() Type
}

// DeclField represents a field declaration in the AST.
type DeclField interface {
	// Inherits from Field interface.
	Field
	// LinkParents links this declaration with its parent nodes within the AST.
	LinkParents(p SemanticParser, visitedDecls *data.AsyncSet[Declaration]) io.Error
	// LinkFields links this declaration's methods with its parent nodes' methods within the AST.
	LinkFields(p SemanticParser, visitedDecls *data.AsyncSet[Declaration]) io.Error

	// Declaration returns the field's declaration.
	Declaration() Declaration
}
