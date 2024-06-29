package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

func genericParent(p ast.SyntaxParser, dt *Lookup) (*Generic, io.Error) {
	var types []ast.Type
	for {
		t, err := p.ParseType()
		if err != nil {
			return nil, err
		}
		types = append(types, t)

		if p.Match(token.TOK_RIGHTBRACKET) {
			break
		} else if _, err := p.Consume(token.TOK_COMMA); err != nil {
			return nil, err
		}
	}

	return &Generic{
		Context_:     p.TypeContext(),
		Type:         dt,
		GenericTypes: types,
		EndToken:     p.Prev(),
	}, nil

}

func syntaxParent(p ast.SyntaxParser) (ast.DeclType, io.Error) {
	t, err := (&Lookup{Context_: p.TypeContext()}).Syntax(p)
	if err != nil {
		return nil, err
	}

	var dt ast.DeclType = t
	if p.Match(token.TOK_LEFTBRACKET) {
		dt, err = genericParent(p, t.(*Lookup))
		if err != nil {
			return nil, err
		}
	}

	if p.Match(token.TOK_TILDE) {
		size := int64(-1)
		return &Tailed{
			Type:     dt,
			Size:     size,
			EndToken: p.Prev(),
		}, nil
	}

	return dt, nil
}

func SyntaxParents(p ast.SyntaxParser) (data.Set[ast.DeclType], io.Error) {
	parents := data.Set[ast.DeclType]{}
	if !p.Match(token.TOK_LEFTBRACKET) {
		t, err := syntaxParent(p)
		if err != nil {
			return nil, err
		}
		parents.Set(t)

		return parents, nil
	}

	for {
		t, err := syntaxParent(p)
		if err != nil {
			return nil, err
		}

		if !parents.Set(t) {
			return nil, io.NewError("duplicate parent type", zap.Stringer("type", t))
		}

		if p.Match(token.TOK_RIGHTBRACKET) {
			break
		} else if _, err := p.Consume(token.TOK_COMMA); err != nil {
			return nil, err
		}
	}

	return parents, nil
}
