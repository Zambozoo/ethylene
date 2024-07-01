package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Field represents expressions of the form
//
//	EXPR '.' IDENTIFIER
type Field struct {
	SuffixedToken
}

func (f *Field) String() string {
	return fmt.Sprintf("%s.%s", f.Expr.String(), f.Token.Value)
}

func (f *Field) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	t, err := f.Expr.Semantic(p)
	if err != nil {
		return nil, err
	}

	dt, ok := t.(ast.DeclType)
	if !ok {
		return nil, io.NewError("only declaration types can have fields",
			zap.Stringer("location", f.Location()),
			zap.Stringer("type", t),
		)
	}

	decl, err := dt.Declaration(p)
	if err != nil {
		return nil, err
	}

	field, ok := decl.Members()[f.Token.Value]
	if !ok {
		field, ok = decl.Methods()[f.Token.Value]
		if !ok {
			return nil, io.NewError("field doesn't exist for expression",
				zap.String("field", f.Token.Value),
				zap.Stringer("location", f.Location()),
				zap.Stringer("type", t),
			)
		}
	}

	return field.Type(), nil
}

// TypeField represents expressions of the form
//
//	'type' '(' TYPE ')' '.' IDENTIFIER
type TypeField struct {
	StartToken token.Token
	Type       ast.Type
	Token      token.Token
}

func (t *TypeField) Location() *token.Location {
	return token.LocationBetween(&t.StartToken, &t.Token)
}

func (t *TypeField) String() string {
	fieldString := t.Token.String()
	if fieldString != "" {
		fieldString = "." + fieldString
	}
	return fmt.Sprintf("type(%s)%s", t.Type.String(), fieldString)
}

func (t *TypeField) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	var err io.Error
	t.StartToken, err = p.Consume(token.TOK_TYPE)
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
		return nil, err
	}

	t.Type, err = p.ParseType()
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(token.TOK_RIGHTPAREN); err != nil {
		return nil, err
	}

	if p.Match(token.TOK_PERIOD) {
		t.Token, err = p.Consume(token.TOK_IDENTIFIER)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (t *TypeField) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	var emptyToken token.Token
	if t.Token == emptyToken {
		if _, err := t.Type.TypeID(p); err != nil {
			return nil, err
		}

		return type_.NewTypeID(), nil
	}

	dt, ok := t.Type.(ast.DeclType)

	if !ok || !dt.IsFieldable() {
		return nil, io.NewError("type field expressions cannot have generic or tailed types",
			zap.Stringer("location", t.Location()),
			zap.Stringer("type", t.Type),
		)
	}

	decl, err := dt.Declaration(p)
	if err != nil {
		return nil, err
	}

	field, ok := decl.StaticMembers()[t.Token.Value]
	if !ok {
		field, ok = decl.Methods()[t.Token.Value]
		if !ok {
			return nil, io.NewError("field doesn't exist for type",
				zap.String("field", t.Token.Value),
				zap.Stringer("location", t.Location()),
				zap.Stringer("type", t.Type),
			)
		}
	}

	return field.Type(), nil
}
