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

	TypeID(SemanticParser) (TypeID, io.Error)
	IsConcrete() bool
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
	IsFieldable() bool
}

type Types interface {
	NextEnumIndex(d Declaration) (uint32, io.Error)
	NextStructIndex(d Declaration) (uint32, io.Error)
	NextClassIndex(d Declaration) (uint32, io.Error)
	NextAbstractIndex(d Declaration) (uint32, io.Error)
	NextInterfaceIndex(d Declaration) (uint32, io.Error)
	NextListIndex(ids []uint64) (uint32, io.Error)

	EnumIndex(index uint32) uint32
	StructIndex(index uint32) uint32
	ClassIndex(index uint32) uint32
	AbstractIndex(index uint32) uint32
	InterfaceIndex(index uint32) uint32
	ListIndex(ids []uint64) uint32

	MaxIndex() uint64
}

type TypeID interface {
	ID() uint64
	Index() uint32
	ListIndex() uint32
}
