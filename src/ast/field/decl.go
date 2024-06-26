package field

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Decl struct {
	Modifiers
	StartToken   token.Token
	Declaration_ ast.Declaration
}

func (d *Decl) Name() *token.Token {
	return d.Declaration_.Name()
}

func (d *Decl) Declaration() ast.Declaration {
	return d.Declaration_
}

func (d *Decl) Location() *token.Location {
	return token.LocationBetween(&d.StartToken, d.Declaration_)
}
func (d *Decl) String() string {
	return d.Declaration_.String()
}

func (d *Decl) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	d.Declaration_, err = p.ParseDecl()
	return err
}

func (d *Decl) Semantic(p ast.SemanticParser) io.Error {
	return d.Declaration_.Semantic(p)
}

func (d *Decl) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	_, err := d.Declaration_.LinkParents(p, visitedDecls, map[string]struct{}{})
	return err
}

func (d *Decl) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return d.Declaration_.LinkFields(p, visitedDecls)
}
