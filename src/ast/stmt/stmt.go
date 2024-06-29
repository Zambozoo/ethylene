package stmt

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type statement interface {
	ast.Statement
	Syntax(p ast.SyntaxParser) io.Error
}

type BoundedStmt struct {
	StartToken token.Token
	EndToken   token.Token
}

func (s *BoundedStmt) Location() *token.Location {
	return token.LocationBetween(&s.StartToken, &s.EndToken)
}

func Syntax(p ast.SyntaxParser) (ast.Statement, io.Error) {
	var s statement
	switch t := p.Peek(); t.Type {
	case token.TOK_LEFTBRACE:
		s = &Block{}
	case token.TOK_BREAK:
		s = &Break{}
	case token.TOK_CONTINUE:
		s = &Continue{}
	case token.TOK_DELETE:
		s = &Delete{}
	case token.TOK_FOR:
		return parseFor(p)
	case token.TOK_IF:
		s = &If{}
	case token.TOK_LABEL:
		s = &Label{}
	case token.TOK_PANIC:
		s = &Panic{}
	case token.TOK_PRINT:
		s = &Print{}
	case token.TOK_RETURN:
		s = &Return{}
	case token.TOK_VAR:
		s = &Var{}
	default:
		s = &Expr{}
	}

	if err := s.Syntax(p); err != nil {
		return nil, err
	}

	return s, nil
}
