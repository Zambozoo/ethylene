package file

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/io/path"
)

type Import struct {
	StartToken token.Token
	Dependency string
	Path_      string
	EndToken   token.Token
	FilePath   path.Path
}

func (i *Import) Path() path.Path {
	return i.FilePath
}

func (i *Import) Location() *token.Location {
	return token.LocationBetween(&i.StartToken, &i.EndToken)
}

func (i Import) String() string {
	if i.Dependency == "" {
		return fmt.Sprintf("import %q", i.Path_)
	}
	return fmt.Sprintf("import %s(%q)", i.Dependency, i.Path())
}

func (i *Import) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if i.StartToken, err = p.Consume(token.TOK_IMPORT); err != nil {
		return err
	}

	if p.Match(token.TOK_STRING) {
		i.Path_ = p.Prev().Value
	} else {
		tok, err := p.Consume(token.TOK_IDENTIFIER)
		if err != nil {
			return err
		}
		i.Dependency = tok.Value

		if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
			return err
		}

		path, err := p.Consume(token.TOK_STRING)
		if err != nil {
			return err
		}
		i.Path_ = path.Value

		if _, err := p.Consume(token.TOK_RIGHTPAREN); err != nil {
			return err
		}
	}

	i.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	if err != nil {
		return err
	}

	i.FilePath, err = p.AddPath(i.Dependency, i.Path_)
	if err != nil {
		return err
	}

	return nil
}
