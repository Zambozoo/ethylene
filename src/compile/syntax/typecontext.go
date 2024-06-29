package syntax

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/io/path"

	"go.uber.org/zap"
)

type TypeContext struct {
	File      ast.File
	Project   *path.Project
	Scope     []ast.Declaration
	SymbolMap SymbolMap
}

func (tc *TypeContext) Declaration(tokens []token.Token) (ast.Declaration, io.Error) {
	if len(tokens) == 1 {
		if _, ok := tc.TopScope().GenericParamIndex(tokens[0].Value); ok {
			return nil, io.NewError("generic argument doesn't have declaration", zap.Stringer("token", &tokens[0]))
		}
	}

	if i, ok := tc.File.GetImport(tokens[0].Value); ok {
		file := tc.SymbolMap.Files[i.Path().String()]
		d := file.Declaration()
		for i := 1; i < len(tokens); i++ {
			decl, ok := d.Declarations()[tokens[i].Value]
			if !ok {
				return nil, io.NewError("missing declaration", zap.Stringer("location", tokens[0].Location()))
			}

			if !decl.HasModifier(ast.MOD_PUBLIC) {
				return nil, io.NewError("inaccessible declaration", zap.Stringer("location", tokens[0].Location()))
			}

			d = decl.Declaration()
		}

		return d, nil
	}

scope:
	for i := len(tc.Scope) - len(tokens); i >= 0; i-- {
		decl := tc.Scope[i]
		if d, ok := decl.Declarations()[tokens[0].Value]; ok {
			decl = d.Declaration()
		}

		d := decl
		for j, token := range tokens {
			if d.Name().Value != token.Value {
				continue scope
			}

			d = tc.Scope[i+j]
		}

		return decl, nil
	}

	return nil, io.NewError("missing declaration", zap.Stringer("location", tokens[0].Location()))
}

func (tc *TypeContext) TopScope() ast.Declaration {
	return tc.Scope[len(tc.Scope)-1]
}
