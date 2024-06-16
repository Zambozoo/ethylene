package semantic

import (
	"geth-cody/ast"
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/compile/syntax"
	"geth-cody/io"
)

type Parser struct {
	scope     *ast.Scope
	bytecodes *bytecode.Bytecodes
	File_     ast.File
}

func NewParser(file ast.File, symbolMap syntax.SymbolMap) *Parser {
	panic("")
}

func (p *Parser) File() ast.File {
	return p.File_
}

func (p *Parser) Scope() *ast.Scope {
	return p.scope
}

func (p *Parser) Parse() (*bytecode.Bytecodes, io.Error) {
	if err := p.File_.Semantic(p); err != nil {
		return nil, err
	}

	return p.bytecodes, nil
}
