package ast

import (
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// TypeContext represents the context in which types are resolved.
type TypeContext interface {
	// Declaration returns the declaration for a given set of tokens
	Declaration(tokens []token.Token) (Declaration, io.Error)
	// TopScope returns the top scope declaration
	TopScope() Declaration
}

// Type represents a type in the AST, like primitives or composite types.
type Type interface {
	// Inherits from Type interface.
	Node

	// Extends returns true if the type extends parent.
	Extends(p SemanticParser, parent Type) (bool, io.Error)
	// ExtendsAsPointer returns true if the type extends parent as a pointer.
	ExtendsAsPointer(p SemanticParser, parent Type) (bool, io.Error)
	// Equals returns true if the type is the same as other.
	Equals(p SemanticParser, other Type) (bool, io.Error)

	// Concretize returns a concrete type, replacing any generics according to the mapping.
	Concretize(mapping []Type) Type

	// IsContant returns true if the type is constant.
	IsConstant() bool
	// SetConstant sets the type to be constant.
	SetConstant()
}

// FunType represents a function type, detailing its return and parameter types.
type FunType interface {
	// Inherits from Type interface.
	Type
	// ReturnType returns the return type of the function.
	ReturnType() Type
	// ParameterTypes returns a list of parameter types for the function.
	ParameterTypes() []Type
	// Arity returns the number of parameters the function takes.
	Arity() int
}

// DeclType represents a declaration type, linking a type to its declaration.
type DeclType interface {
	// Inherits from Type interface.
	Type

	// Name returns the name of the declaration
	Name() token.Token

	Context() TypeContext
	// Declaration returns the Declaration associated with the type within a given context.
	Declaration(p SemanticParser) (Declaration, io.Error)
}
