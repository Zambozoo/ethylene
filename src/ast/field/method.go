package field

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"

	"go.uber.org/zap"
)

type Method struct {
	Modifiers
	StartToken token.Token
	EndToken   token.Token

	Name_      token.Token
	Parameters []*token.Token
	Type       ast.FunType
	Stmt       ast.Statement
}

func (m *Method) Name() *token.Token {
	return &m.Name_
}
func (m *Method) Location() token.Location {
	var locatable token.Locatable = m.Stmt
	if m.Stmt != nil {
		locatable = &m.EndToken
	}

	return token.LocationBetween(&m.StartToken, locatable)
}

func (m *Method) String() string {
	return fmt.Sprintf("Member{Name:%s, Modifiers:%s, Type:%s, Parameters:%s, Stmt:%s}",
		m.Name(),
		m.Modifiers.String(),
		m.Type.String(),
		strs.Strings(m.Parameters),
		m.Stmt.String(),
	)
}

func (m *Method) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	if m.StartToken, err = p.Consume(token.TOK_FUN); err != nil {
		return err
	}

	var ok bool
	t, err := p.ParseType()
	if err != nil {
		return err
	}

	if m.Type, ok = t.(ast.FunType); !ok {
		return io.NewError("expected a function type for method", zap.Any("location", m.StartToken.Location()))
	}

	m.Name_, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	if !p.Match(token.TOK_ASSIGN) {
		m.EndToken, err = p.Consume(token.TOK_SEMICOLON)
		return err

	}

	if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
		return err
	}

	if !p.Match(token.TOK_RIGHTPAREN) {
		for {
			tok, err := p.Consume(token.TOK_IDENTIFIER)
			if err != nil {
				return err
			}

			m.Parameters = append(m.Parameters, &tok)
			if p.Match(token.TOK_RIGHTPAREN) {
				break
			}

			if _, err := p.Consume(token.TOK_COMMA); err != nil {
				return err
			}
		}
	}

	m.Stmt, err = p.ParseStmt()
	return err
}

type methodVariable struct {
	name  *token.Token
	type_ ast.Type
}

func (m *methodVariable) Name() *token.Token {
	return m.name
}
func (m *methodVariable) Type() ast.Type {
	return m.type_
}

func (m *Method) Semantic(p ast.SemanticParser) io.Error {
	// TODO: Bytecode
	if m.Stmt != nil {
		p.Scope().Wrap()
		defer p.Scope().Unwrap()
		for i, t := range m.Type.ParameterTypes() {
			name := m.Parameters[i]
			p.Scope().AddVariable(&methodVariable{name: name, type_: t})
		}

		t, err := m.Stmt.Semantic(p)
		if err != nil {
			return err
		}

		if t == nil {
			io.NewError("method missing return a value",
				zap.Any("name", m.Name()),
				zap.Any("location", m.Location()),
			).Log(io.Warnf)
		}

		if _, err := type_.MustExtend(t, m.Type.ReturnType()); err != nil {
			return err
		}
	}

	return nil
}
