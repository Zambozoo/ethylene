package semantic

import (
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax"
	"geth-cody/io"

	"go.uber.org/zap"
)

type TypeContext struct {
	scope        []ast.Declaration
	fileEntry    syntax.FileEntry
	projectFiles map[string]*io.Project
	fileEntries  map[string]syntax.FileEntry
}

func (tc *TypeContext) Project() *io.Project {
	return tc.fileEntry.Project
}

func (tc *TypeContext) Dependency(pkg string) (string, bool) {
	if version, ok := tc.fileEntry.Project.Packages[pkg]; ok {
		return version, ok
	}

	return "", false
}

func (tc *TypeContext) Declaration(tokens []token.Token) (ast.Declaration, io.Error) {
	if i, ok := tc.fileEntry.File.GetImport(tokens[0].Value); ok {
		fileEntry := tc.fileEntries[i.Path().String()]
		d := fileEntry.File.Declaration()
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

// Equals returns true if a and b are the same type given a context.
func (tc *TypeContext) Equals(a, b ast.Type) (bool, io.Error) {
	return a.Equals(tc, b)
}

// Extends returns true if child extends parent given a context.
func (tc *TypeContext) Extends(child, parent ast.Type) (bool, io.Error) {
	return child.Extends(tc, parent)
}

// MustExtend returns the parent type that extends child and parent given a context, or an erro if none exists.
func (tc *TypeContext) MustExtend(child, parent ast.Type, parents ...ast.Type) (ast.Type, io.Error) {
	return type_.MustExtend(tc, child, parent, parents...)
}
