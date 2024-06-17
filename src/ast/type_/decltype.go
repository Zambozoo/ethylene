package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

func singleTokenCompositeSyntax(p ast.SyntaxParser) (*Composite, io.Error) {
	tok, err := p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	return &Composite{
		Context_: p.TypeContext(),
		Tokens:   []token.Token{tok},
	}, nil
}

func typeList(p ast.SyntaxParser) ([]ast.Type, io.Error) {
	if !p.Match(token.TOK_LEFTBRACKET) {
		t, err := p.ParseType()
		if err != nil {
			return nil, err
		}

		return []ast.Type{t}, nil
	}

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

	return types, nil
}

func genericSyntax(p ast.SyntaxParser, dt *Composite, d ast.Declaration) (*Generic, io.Error) {
	generics := map[string]ast.GenericTypeArg{}
	var types []ast.GenericTypeArg
	for {
		t, err := singleTokenCompositeSyntax(p)
		if err != nil {
			return nil, err
		}

		if _, exists := generics[t.Tokens[0].Value]; exists {
			return nil, io.NewError("duplicate generic type",
				zap.Any("type", t.Tokens[0].Value),
				zap.Any("location", t.Tokens[0].Location),
			)
		}
		generics[t.Tokens[0].Value] = t

		var gta ast.GenericTypeArg = t
		if p.Match(token.TOK_SUBTYPE) {
			types, err := typeList(p)
			if err != nil {
				return nil, err
			}

			gta = &Sub{
				Type:  t,
				Types: types,
			}
		} else if p.Match(token.TOK_SUPERTYPE) {
			types, err := typeList(p)
			if err != nil {
				return nil, err
			}

			gta = &Super{
				Type:  t,
				Types: types,
			}
		}

		types = append(types, gta)

		if p.Match(token.TOK_RIGHTBRACKET) {
			break
		} else if _, err := p.Consume(token.TOK_COMMA); err != nil {
			return nil, err
		}
	}

	if err := d.SetGenerics(generics); err != nil {
		return nil, err
	}

	return &Generic{
		Context_:     p.TypeContext(),
		Type:         dt,
		GenericTypes: types,
		EndToken:     p.Prev(),
	}, nil
}

// DeclTypeSyntax parses a declaration type, either a Composite, Tailed, or Generic.
func SyntaxDecl(p ast.SyntaxParser, d ast.Declaration) (ast.DeclType, io.Error) {
	t, err := singleTokenCompositeSyntax(p)
	if err != nil {
		return nil, err
	}
	d.SetName(t.Tokens[0])

	var dt ast.DeclType = t
	if p.Match(token.TOK_LEFTBRACKET) {
		dt, err = genericSyntax(p, t, d)
		if err != nil {
			return nil, err
		}
	}

	if p.Match(token.TOK_TILDE) {
		if err := d.SetTailed(); err != nil {
			return nil, err
		}

		size := int64(-1)
		return &Tailed{
			Type:     dt,
			Size:     size,
			EndToken: p.Prev(),
		}, nil
	}

	return dt, nil
}
