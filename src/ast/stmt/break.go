package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Break represents a break statement
type Break struct {
	BoundedStmt
	Label token.Token
}

func (b *Break) String() string {
	var labelString string
	if b.Label.Value != "" {
		labelString = " " + b.Label.Value
	}
	return fmt.Sprintf("break%s;", labelString)
}

func (b *Break) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	b.BoundedStmt.StartToken, err = p.Consume(token.TOK_BREAK)
	if err != nil {
		return err
	}

	if p.Match(token.TOK_IDENTIFIER) {
		b.Label = p.Prev()
	} else {
		b.Label = emptyLabel
	}

	b.BoundedStmt.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	return err
}

func (b *Break) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	_, err := p.Scope().GetLabel(b.Label)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
