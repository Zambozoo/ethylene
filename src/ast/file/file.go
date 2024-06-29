package file

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/io/path"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type File struct {
	Imports      map[string]Import
	Declaration_ ast.Declaration

	Path_   path.Path
	Project *path.Project
}

func (f *File) Path() path.Path {
	return f.Path_
}

func (f *File) GetImport(name string) (ast.Import, bool) {
	if i, ok := f.Imports[name]; ok {
		return &i, ok
	}

	return nil, false
}

func (f *File) Declaration() ast.Declaration {
	return f.Declaration_
}

func (f *File) Location() *token.Location {
	var locatable token.Locatable = f.Declaration_
	if len(f.Imports) > 0 {
		for _, v := range f.Imports {
			locatable = &v
			break
		}
	}

	return token.LocationBetween(locatable, f.Declaration_)
}

func (f *File) String() string {
	return fmt.Sprintf("%s\n%s\n", strings.Join(maps.Keys(f.Imports), "\n"), f.Declaration_.String())
}

func (f *File) Syntax(p ast.SyntaxParser) io.Error {
	for p.Peek().Type == token.TOK_IMPORT {
		i := Import{}
		if err := i.Syntax(p); err != nil {
			return err
		}
		name := i.FilePath.Name()
		if _, ok := f.Imports[name]; ok {
			return io.NewError("duplicate import",
				zap.Stringer("path", i.FilePath),
				zap.Stringer("location", i.Location()),
			)
		}
		f.Imports[name] = i
	}

	var err io.Error
	f.Declaration_, err = p.ParseDecl()
	return err
}

func New(p ast.SyntaxParser) ast.File {
	return &File{
		Imports: map[string]Import{},
		Path_:   p.Path(),
	}
}

func (f *File) Semantic(p ast.SemanticParser) io.Error {
	return f.Declaration_.Semantic(p)
}

func (f *File) LinkParents(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	_, err := f.Declaration_.LinkParents(p, visitedDecls, map[string]struct{}{})
	return err
}

func (f *File) LinkFields(p ast.SemanticParser, visitedDecls *data.AsyncSet[ast.Declaration]) io.Error {
	return f.Declaration_.LinkFields(p, visitedDecls)
}
