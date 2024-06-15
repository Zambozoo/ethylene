package field

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Member struct {
	Modifiers
	StartToken token.Token
	EndToken   token.Token

	Type_ ast.Type
	Name_ token.Token
	Expr  ast.Expression // Optional initial value
}

func (m *Member) Name() *token.Token {
	return &m.Name_
}

func (m *Member) Type() ast.Type {
	return m.Type_
}

func (m *Member) Location() token.Location {
	return token.LocationBetween(&m.StartToken, &m.EndToken)
}
func (m *Member) String() string {
	return fmt.Sprintf("Member{Name:%s, Modifiers:%s, Type:%s, Expr:%s}",
		m.Name(),
		m.Modifiers.String(),
		m.Type_.String(),
		m.Expr.String(),
	)
}

func (m *Member) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if m.StartToken, err = p.Consume(token.TOK_VAR); err != nil {
		return err
	}

	m.Type_, err = p.ParseType()
	if err != nil {
		return err
	}

	m.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	if p.Match(token.TOK_ASSIGN) {
		m.Expr, err = p.ParseExpr()
		if err != nil {
			return err
		}
	}

	m.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	if err != nil {
		return err
	}

	return nil
}

func (m *Member) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Bytecode
	if m.Expr != nil {
		t, err := m.Expr.Semantic(p)
		if err != nil {
			return err
		} else if _, err := p.TypeContext().MustExtend(t, m.Type_); err != nil {
			return err
		}
	}

	return nil
}
