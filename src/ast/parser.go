package ast

import (
	"fmt"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/io/path"
)

// Node is an interface for all AST nodes, providing basic functionalities.
type Node interface {
	// Allows the node to be converted to a string.
	fmt.Stringer
	// Returns the location of the node in the source code.
	Location() *token.Location
}

// SyntaxParser defines an interface for parsing syntax elements of a language.
type SyntaxParser interface {
	// AddPath ensures a file path dependency:path will be lexed.
	AddPath(dependency, path string) (path.Path, io.Error)

	// Path returns the path associated with the parser
	Path() path.Path
	// File returns the File associated with the parser
	File() File

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

	// WrapScope wraps the current typing scope
	WrapScope(decl Declaration)
	// UnwrapScope unwraps the current typing scope
	UnwrapScope()
	// TypeContext returns a copy of the current type context
	TypeContext() TypeContext

	// ParseType parses and returns a Type.
	ParseType() (Type, io.Error)
	// ParseParentTypes parses and returns parent DeclTypes.
	ParseParentTypes() (data.Set[DeclType], io.Error)
	// ParseDecl parses and returns a Declaration.
	ParseDecl() (Declaration, io.Error)
	// ParseField parses and returns a Field.
	ParseField() (Field, io.Error)
	// ParseStmt parses and returns a Statement.
	ParseStmt() (Statement, io.Error)
	// ParseExpr parses and returns an Expression.
	ParseExpr() (Expression, io.Error)

	Types() Types
}

// SemanticParser defines an interface for semantic analysis of parsed elements.
type SemanticParser interface {
	// File returns the file associated with the parser
	File() File
	// Scope returns the current scope for resolving identifiers.
	Scope() *Scope
	// WrapDeclWithGeneric wraps the declaration with a mapping.
	WrapDeclWithGeneric(d Declaration, slice []Type) Declaration
	
	Types() Types
}
