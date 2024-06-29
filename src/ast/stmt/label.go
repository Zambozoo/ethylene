package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
)

var emptyLabel = token.Token{}

// Statement represents a label statement
type Label struct {
	StartToken token.Token
	Label      token.Token
	Stmt       ast.Statement
}

func (l *Label) Location() token.Location {
	return token.LocationBetween(&l.StartToken, l.Stmt)
}

func (l *Label) String() string {
	return fmt.Sprintf("label %s: %s", l.Label.Value, l.Stmt.String())
}

func (l *Label) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	l.StartToken, err = p.Consume(token.TOK_LABEL)
	if err != nil {
		return err
	}

	l.Label, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return err
	}

	if _, err := p.Consume(token.TOK_COLON); err != nil {
		return err
	}

	l.Stmt, err = p.ParseStmt()
	if err != nil {
		return err
	}

	_, ifFor0 := l.Stmt.(*For0)
	_, ifFor1 := l.Stmt.(*For1)
	if !ifFor0 && !ifFor1 {
		return io.NewError("label statement must be followed by a for0 statement",
			zap.Any("label", l.Label.Value),
			zap.Any("location", l.Location()),
		)
	}

	return err
}

func (l *Label) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	p.Scope().Wrap()
	defer p.Scope().Unwrap()
	p.Scope().AddLabel(&l.Label)

	return l.Stmt.Semantic(p)
}
