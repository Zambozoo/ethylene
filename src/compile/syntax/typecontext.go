package syntax

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/io/path"

	"go.uber.org/zap"
)

type TypeContext struct {
	file      ast.File
	project   *path.Project
	scope     []ast.Declaration
	symbolMap SymbolMap

	generics map[string]ast.DeclType
}

func (tc *TypeContext) Dependency(pkg string) (string, bool) {
	if version, ok := tc.project.Packages[pkg]; ok {
		return version, ok
	}

	return "", false
}

func (tc *TypeContext) Declaration(tokens []token.Token) (ast.Declaration, io.Error) {
	if i, ok := tc.file.GetImport(tokens[0].Value); ok {
		file := tc.symbolMap.Files[i.Path().String()]
		d := file.Declaration()
		for i := 1; i < len(tokens); i++ {
			decl, ok := d.Declarations()[tokens[i].Value]
			if !ok {
				return nil, io.NewError("missing declaration", zap.Any("location", tokens[0].Location()))
			}

			if !decl.HasModifier(ast.MOD_PUBLIC) {
				return nil, io.NewError("inaccessible declaration", zap.Any("location", tokens[0].Location()))
			}

			d = decl.Declaration()
		}

		return d, nil
	}

scope:
	for i := len(tc.scope) - len(tokens); i >= 0; i-- {
		decl := tc.scope[i]
		if d, ok := decl.Declarations()[tokens[0].Value]; ok {
			decl = d.Declaration()
		}

		for j, token := range tokens {
			if decl.Name().Value != token.Value {
				continue scope
			}

			decl = tc.scope[i+j]
		}

		return decl, nil
	}

	return nil, io.NewError("missing declaration", zap.Any("location", tokens[0].Location()))
}

func (tc *TypeContext) Generics() map[string]ast.DeclType {
	return tc.generics
}
