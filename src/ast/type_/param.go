package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Param is the variable name: `E` in `class List[E]`
type Param struct {
	token.Token
	Constant bool
	Context_ ast.TypeContext
	Decl     ast.Declaration
	Index    int
}

func (p *Param) String() string {
	return p.Token.Value
}

func (p *Param) Declarations() map[string]ast.DeclField {
	return nil
}
func (p *Param) GenericMap(ast.SemanticParser) (map[string]ast.Type, io.Error) {
	return nil, nil
}
func (p *Param) IsAbstract() bool  { return false }
func (p *Param) IsClass() bool     { return false }
func (p *Param) IsInterface() bool { return false }
func (p *Param) LinkFields(_ ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return nil
}
func (p *Param) LinkParents(_ ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration], cycleMap map[string]struct{}) (data.Set[ast.DeclType], io.Error) {
	return nil, nil
}

func (p *Param) Name() *token.Token {
	return nil
}
func (p *Param) Members() map[string]ast.Member {
	return nil
}
func (p *Param) StaticMembers() map[string]ast.Member {
	return nil
}
func (p *Param) Methods() map[string]ast.Method {
	return nil
}
func (p *Param) StaticMethods() map[string]ast.Method {
	return nil
}
func (p *Param) Semantic(_ ast.SemanticParser) io.Error {
	return nil
}
func (p *Param) Syntax(_ ast.SyntaxParser) (ast.Declaration, io.Error) {
	return nil, nil
}

// Equals returns true if the type is the same as other.
func (p *Param) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherParam, ok := other.(*Param); ok {
		return p.Token.Value == otherParam.Token.Value && p.Decl == otherParam.Decl, nil
	}
	return false, nil
}
func (p *Param) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Param) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}
func (p *Param) Concretize(types []ast.Type) ast.Type {
	return types[p.Index]
}
func (p *Param) IsConstant() bool {
	return p.Constant
}
func (p *Param) SetConstant() {
	p.Constant = true
}
