package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Continue represents a break statement
type Continue struct {
	BoundedStmt
	Label token.Token
}

func (c *Continue) String() string {
	var labelString string
	if c.Label.Value != "" {
		labelString = " " + c.Label.Value
	}
	return fmt.Sprintf("continue%s;", labelString)
}

func (c *Continue) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	c.BoundedStmt.StartToken, err = p.Consume(token.TOK_CONTINUE)
	if err != nil {
		return err
	}

	if p.Match(token.TOK_IDENTIFIER) {
		c.Label = p.Prev()
	} else {
		c.Label = emptyLabel
	}

	c.BoundedStmt.EndToken, err = p.Consume(token.TOK_SEMICOLON)
	return err
}

func (c *Continue) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	_, err := p.Scope().GetLabel(c.Label)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
