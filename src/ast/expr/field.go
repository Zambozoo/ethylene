package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Field represents expressions of the form
//
//	EXPR '.' IDENTIFIER
type Field struct {
	SuffixedToken
}

func (f *Field) String() string {
	return fmt.Sprintf("Field{Expr:%s,Value:%s}", f.Expr.String(), f.Token.Value)
}

func (f *Field) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}

// TypeField represents expressions of the form
//
//	'type' '(' TYPE ')' '.' IDENTIFIER
type TypeField struct {
	StartToken token.Token
	Type       ast.Type
	FieldName  token.Token
}

func (t *TypeField) Location() token.Location {
	return token.LocationBetween(&t.StartToken, &t.FieldName)
}

func (t *TypeField) String() string {
	return fmt.Sprintf("TypeField{Type:%s, FieldName:%s}", t.Type.String(), t.FieldName.String())
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

	if _, err := p.Consume(token.TOK_PERIOD); err != nil {
		return nil, err
	}

	t.FieldName, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TypeField) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}
