package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"

	"go.uber.org/zap"
)

type GenericSubtype []ast.Type

func (g GenericSubtype) String() string {
	return fmt.Sprintf("GenericSubtype{%s}", strs.Strings(g))
}

type GenericSupertype []ast.Type

func (g GenericSupertype) String() string {
	return fmt.Sprintf("GenericSubtype{%s}", strs.Strings(g))
}

func syntaxDeclTypes(p ast.SyntaxParser) (data.Set[ast.DeclType], io.Error) {
	declTypes := data.Set[ast.DeclType]{}
	syntaxTypes, err := syntaxTypes(p)
	if err != nil {
		return nil, err
	}

	for _, t := range syntaxTypes {
		dt, isDeclType := t.(ast.DeclType)
		if !isDeclType {
			return nil, io.NewError("expected parent DeclType", zap.Any("type", t))
		} else if _, exists := declTypes.Get(dt); exists {
			return nil, io.NewError("duplicate parent type", zap.Any("type", t))
		}

		declTypes.Set(dt)
	}

	return declTypes, nil
}

func syntaxTypes(p ast.SyntaxParser) ([]ast.Type, io.Error) {
	var types []ast.Type
	if p.Match(token.TOK_LEFTBRACKET) {
		for {
			t, err := p.ParseType()
			if err != nil {
				return nil, err
			}

			types = append(types, t)
			if p.Match(token.TOK_RIGHTBRACKET) {
				break
			}

			if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return nil, err
			}
		}
	} else {
		t, err := p.ParseType()
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}

	return types, nil
}

func syntaxGenericConstraints(p ast.SyntaxParser) (map[string]ast.GenericConstraint, io.Error) {
	genericConstraints := map[string]ast.GenericConstraint{}
	if p.Match(token.TOK_LEFTBRACKET) {
		for {
			name, err := p.Consume(token.TOK_IDENTIFIER)
			if err != nil {
				return nil, err
			}
			if p.Match(token.TOK_SUBTYPE) {
				ts, err := syntaxTypes(p)
				if err != nil {
					return nil, err
				}
				genericConstraints[name.Value] = GenericSubtype(ts)
			} else if p.Match(token.TOK_SUPERTYPE) {
				ts, err := syntaxTypes(p)
				if err != nil {
					return nil, err
				}
				genericConstraints[name.Value] = GenericSupertype(ts)
			} else {
				genericConstraints[name.Value] = nil
			}

			if p.Match(token.TOK_RIGHTBRACKET) {
				break
			}

			if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return nil, err
			}
		}
	}

	return genericConstraints, nil
}
