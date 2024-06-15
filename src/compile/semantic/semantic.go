package semantic

import (
	"geth-cody/ast"
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/compile/syntax"
	"geth-cody/io"
)

type Parser struct {
	typeContext *TypeContext
	scope       *ast.Scope
	bytecodes   *bytecode.Bytecodes
}

func NewParser(fileEntry syntax.FileEntry, projectFiles map[string]*io.Project, fileEntries map[string]syntax.FileEntry) *Parser {
	return &Parser{
		typeContext: &TypeContext{
			fileEntry:    fileEntry,
			projectFiles: projectFiles,
			fileEntries:  fileEntries,
		},
	}
}

func (p *Parser) TypeContext() ast.TypeContext {
	return p.typeContext
}

func (p *Parser) WrapTypeContext(decl ast.Declaration) {
	p.typeContext.scope = append(p.typeContext.scope, decl)
}

func (p *Parser) UnwrapTypeContext() {
	p.typeContext.scope = p.typeContext.scope[:len(p.typeContext.scope)-1]
}

func (p *Parser) Scope() *ast.Scope {
	return p.scope
}

func (p *Parser) Parse() (*bytecode.Bytecodes, io.Error) {
	if err := p.typeContext.fileEntry.File.Semantic(p); err != nil {
		return nil, err
	}

	return p.bytecodes, nil
}
