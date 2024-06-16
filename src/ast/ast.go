package ast

import (
	"fmt"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// SyntaxParser defines an interface for parsing syntax elements of a language.
type SyntaxParser interface {
	// AddPath ensures a file path dependency:path will be lexed.
	AddPath(dependency, path string) (io.Path, io.Error)

	Path() io.Path

	// Peek returns the next token without consuming it.
	Peek() token.Token
	// Prev returns the previously consumed token with unconsuming it.
	Prev() token.Token
	// Next consumes and returns the next token.
	Next() token.Token
	// Match checks if the next token matches any of the provided token types.
	Match(ts ...token.Type) bool
	// Consume checks if the next token is of the specified type and consumes it.
	Consume(t token.Type) (token.Token, io.Error)

	WrapScope(decl Declaration)
	UnwrapScope()
	TypeContext() TypeContext

	// ParseType parses and returns a Type.
	ParseType() (Type, io.Error)
	// ParseDecl parses and returns a Declaration.
	ParseDecl() (Declaration, io.Error)
	// ParseField parses and returns a Field.
	ParseField() (Field, io.Error)
	// ParseStmt parses and returns a Statement.
	ParseStmt() (Statement, io.Error)
	// ParseExpr parses and returns an Expression.
	ParseExpr() (Expression, io.Error)
}

// SemanticParser defines an interface for semantic analysis of parsed elements.
type SemanticParser interface {
	File() File
	// Scope returns the current scope for resolving identifiers.
	Scope() *Scope
}

// Node is an interface for all AST nodes, providing basic functionalities.
type Node interface {
	// Allows the node to be converted to a string.
	fmt.Stringer
	// Returns the location of the node in the source code.
	Location() token.Location
}

// Import represents a file import.
type Import interface {
	Path() io.Path
}

// File represents a file node in the AST.
type File interface {
	// Inherits Node interface.
	Node

	GetImport(name string) (Import, bool)
	Declaration() Declaration

	Syntax(p SyntaxParser) io.Error

	// LinkParents links parent nodes within the AST, facilitating inheritance.
	LinkParents(p SemanticParser, visitedDecls *data.AsyncSet[Declaration]) io.Error
	// Semantic performs semantic analysis on the file node.
	Semantic(p SemanticParser) io.Error
}

type GenericConstraint interface {
	fmt.Stringer
}

// Declaration represents a composite object declaration node in the AST.
type Declaration interface {
	// Inherits Node interface.
	Node
	// Syntax performs syntax analysis on the declaration.
	Syntax(p SyntaxParser) io.Error
	// LinkParents links parent nodes within the AST for this declaration.
	LinkParents(p SemanticParser, visitedDecls *data.AsyncSet[Declaration], cycleMap map[string]struct{}) io.Error
	// Semantic performs semantic analysis on the declaration.
	Semantic(p SemanticParser) io.Error

	Generics() map[string]GenericConstraint

	// Name returns the token representing the declaration's name.
	Name() *token.Token
	// Members returns a map of member names to Member objects.
	Members() map[string]Member
	// Methods returns a map of method names to Method objects.
	Methods() map[string]Method
	//Declarations returns a map of inner declarations.
	Declarations() map[string]DeclField
}

type ChildDeclaration interface {
	Declaration
	Parents() []Type
}

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

// Statement represents a statement in the AST.
type Statement interface {
	// Inherits Node interface.
	Node
	// Semantic performs semantic analysis on the statement, returning its return type.
	Semantic(p SemanticParser) (Type, io.Error)
}

// Expression represents an expression in the AST.
type Expression interface {
	// Inherits Node interface.
	Node
	// Semantic performs semantic analysis on the expression, returning its type.
	Semantic(p SemanticParser) (Type, io.Error)
}

// TypeContext represents the context in which types are resolved.
type TypeContext interface {
	Declaration(tokens []token.Token) (Declaration, io.Error)
	Dependency(pkg string) (string, bool)
}

// Type represents a type in the AST, like primitives or composite types.
type Type interface {
	// Inherits Node interface.
	Node
	// Extends returns true if the type extends parent.
	Extends(parent Type) (bool, io.Error)
	// ExtendsAsPointer returns true if the type extends parent as a pointer.
	ExtendsAsPointer(parent Type) (bool, io.Error)
	// Equals returns true if the type is the same as other.
	Equals(other Type) (bool, io.Error)
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

	Context() TypeContext
	// Declaration returns the Declaration associated with the type within a given context.
	Declaration() (Declaration, io.Error)
}

// Method represents a method, which is a function associated with a type.
type Method interface {
	// Inherits from Field interface.
	Field
	// Type returns the type of the method.
	Type() Type
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

	// Declaration returns the field's declaration.
	Declaration() Declaration
}
