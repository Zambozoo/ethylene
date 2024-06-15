package decl

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"

	"go.uber.org/zap"
)

type GenericConstraint interface {
	fmt.Stringer
}
type GenericSubtype []ast.Type

func (g GenericSubtype) String() string {
	return fmt.Sprintf("GenericSubtype{%s}", strs.Strings(g))
}

type GenericSupertype []ast.Type

func (g GenericSupertype) String() string {
	return fmt.Sprintf("GenericSubtype{%s}", strs.Strings(g))
}

func syntaxDeclTypes(p ast.SyntaxParser) ([]ast.DeclType, io.Error) {
	var declTypes []ast.DeclType
	syntaxTypes, err := syntaxTypes(p)
	if err != nil {
		return nil, err
	}

	for _, t := range syntaxTypes {
		if dt, isDeclType := t.(ast.DeclType); isDeclType {
			declTypes = append(declTypes, dt)
		} else {
			return nil, io.NewError("expected DeclType", zap.Any("type", t))
		}
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

func syntaxGenericConstraints(p ast.SyntaxParser) (map[string]GenericConstraint, io.Error) {
	genericConstraints := map[string]GenericConstraint{}
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
