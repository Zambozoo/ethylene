package syntax

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/decl"
	"geth-cody/ast/expr"
	"geth-cody/ast/field"
	"geth-cody/ast/file"
	"geth-cody/ast/stmt"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"maps"
	"slices"

	"go.uber.org/zap"
)

type FileEntry struct {
	Path io.Path
	File ast.File
}

type Chan[T any] interface {
	Send(T)
}

type Parser struct {
	path          io.Path
	tokens        []token.Token
	curTokenIndex int
	scope         []ast.Declaration
	file          ast.File

	symbolMap      SymbolMap
	unvisitedPaths Chan[io.Path]
	project        *io.Project
	mainDirPath    *io.FilePath
}

func NewParser(tokens []token.Token, project *io.Project, path, mainDirPath io.Path, unvisitedPaths Chan[io.Path], symbolMap SymbolMap) *Parser {
	return &Parser{
		path:           path,
		tokens:         tokens,
		symbolMap:      symbolMap,
		unvisitedPaths: unvisitedPaths,
		project:        project,
	}
}

func (p *Parser) File() ast.File {
	return p.file
}

func (p *Parser) Path() io.Path {
	return p.path
}

func (p *Parser) AddPath(dependency, path string) (io.Path, io.Error) {
	var targetPath io.Path
	if dependency == "" {
		targetPath = p.mainDirPath.Join(path)
	} else {
		version, ok := p.project.Packages[dependency]
		if !ok {
			return nil, io.NewError("couldn't find dependency in project",
				zap.String("dependency", dependency),
				zap.Any("path", p.path),
				zap.Any("project", p.project),
			)
		}
		zipFileName := fmt.Sprintf("pkgs/%s~%s.zip", dependency, version)
		zipFilePath := fmt.Sprintf("%s:%s", p.mainDirPath.Join(zipFileName).String(), path)

		var err io.Error
		if targetPath, err = io.NewPath(zipFilePath); err != nil {
			return nil, err
		}
	}

	if err := targetPath.Stat(); err != nil {
		return nil, err
	}

	p.unvisitedPaths.Send(targetPath)
	return targetPath, nil
}

func (p *Parser) Peek() token.Token {
	return p.tokens[p.curTokenIndex]
}

func (p *Parser) Prev() token.Token {
	return p.tokens[p.curTokenIndex-1]
}

func (p *Parser) Next() token.Token {
	p.curTokenIndex++
	return p.tokens[p.curTokenIndex-1]
}

func (p *Parser) Match(ts ...token.Type) bool {
	for _, t := range ts {
		if p.tokens[p.curTokenIndex].Type == t {
			p.curTokenIndex++
			return true
		}
	}

	return false
}

func (p *Parser) Consume(t token.Type) (token.Token, io.Error) {
	if p.tokens[p.curTokenIndex].Type == t {
		p.curTokenIndex++
		return p.tokens[p.curTokenIndex-1], nil
	}

	return token.Token{}, io.NewError("expected token type",
		zap.Any("expected", t),
		zap.Any("actual", p.tokens[p.curTokenIndex]),
	)
}

func (p *Parser) WrapScope(decl ast.Declaration) {
	p.scope = append(p.scope, decl)
}

func (p *Parser) UnwrapScope() {
	p.scope = p.scope[:len(p.scope)-1]
}

func (p *Parser) Scope() []ast.Declaration {
	return p.scope
}

func (p *Parser) TypeContext() ast.TypeContext {
	return &TypeContext{
		file:      p.file,
		project:   p.project,
		scope:     slices.Clone(p.scope),
		symbolMap: p.symbolMap,
		generics:  maps.Clone(p.scope[len(p.scope)-1].GenericsMap()),
	}
}

func (p *Parser) Parse() (ast.File, io.Error) {
	p.file = file.New(p)
	if err := p.file.Syntax(p); err != nil {
		return nil, err
	}

	return p.file, nil
}

func (p *Parser) ParseType() (ast.Type, io.Error) {
	return type_.Syntax(p)
}

func (p *Parser) ParseDeclType(d ast.Declaration) (ast.DeclType, io.Error) {
	return type_.SyntaxDecl(p, d)
}

func (p *Parser) ParseParentTypes() ([]ast.DeclType, io.Error) {
	return type_.SyntaxParents(p)
}

func (p *Parser) ParseDecl() (ast.Declaration, io.Error) {
	return decl.Syntax(p)
}
func (p *Parser) ParseField() (ast.Field, io.Error) {
	return field.Syntax(p)
}
func (p *Parser) ParseStmt() (ast.Statement, io.Error) {
	return stmt.Syntax(p)
}
func (p *Parser) ParseExpr() (ast.Expression, io.Error) {
	return expr.Syntax(p)
}
