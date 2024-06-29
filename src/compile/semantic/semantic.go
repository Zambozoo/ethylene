package semantic

import (
	"geth-cody/ast"
	"geth-cody/ast/decl/generics"
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/compile/syntax"
	"geth-cody/io"
)

type Parser struct {
	scope     *ast.Scope
	bytecodes *bytecode.Bytecodes
	File_     ast.File
	symbolMap syntax.SymbolMap
}

func NewParser(file ast.File, symbolMap syntax.SymbolMap) *Parser {
	return &Parser{
		scope:     ast.NewScope(),
		bytecodes: &bytecode.Bytecodes{},
		File_:     file,
		symbolMap: symbolMap,
	}
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

func (p *Parser) WrapDeclWithGeneric(d ast.Declaration, symbolSlice []ast.Type) ast.Declaration {
	return generics.NewDecl(d, symbolSlice)
}

func (p *Parser) Types() ast.Types {
	return p.symbolMap.Types
}