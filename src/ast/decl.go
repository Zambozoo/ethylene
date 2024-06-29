package ast

import (
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Declaration represents a composite object declaration node in the AST.
type Declaration interface {
	Type
	// Syntax performs syntax analysis on the declaration.
	Syntax(p SyntaxParser) (Declaration, io.Error)
	// LinkParents links parent nodes within the AST for this declaration.
	LinkParents(p SemanticParser, visitedDecls *data.AsyncSet[Declaration], cycleMap map[string]struct{}) (data.Set[DeclType], io.Error)
	// LinkFields links parent and child method nodes within the AST for this declaration.
	LinkFields(p SemanticParser, visitedDecls *data.AsyncSet[Declaration]) io.Error

	// Semantic performs semantic analysis on the declaration.
	Semantic(p SemanticParser) io.Error

	// Name returns the token representing the declaration's name.
	Name() *token.Token
	// Members returns a map of member names to Member objects.
	Members() map[string]Member
	// Methods returns a map of method names to Method objects.
	Methods() map[string]Method
	//Declarations returns a map of inner declarations.
	Declarations() map[string]DeclField

	// IsInterface returns true if the declaration is an interface
	IsInterface() bool
	// IsAbstract returns true if the declaration is an abstract
	IsAbstract() bool
	// IsClass returns true if the declaration is a class
	IsClass() bool

	// GenericParamIndex returns the index of a given genericParam for the declaration and an existence flag
	GenericParamIndex(name string) (int, bool)
	// Generics returns the declaration's generic mapping
	Generics() []Type
}

type ChildDeclaration interface {
	Declaration
	// Parents returns the child declaration's parents
	Parents() data.Set[DeclType]
}
