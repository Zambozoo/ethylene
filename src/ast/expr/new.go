package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type New struct {
	StartToken token.Token
	EndToken   token.Token
	Type       ast.Type

	ArrayLengthExpr ast.Expression
	TailLengthExpr  ast.Expression
}

func (n *New) Location() token.Location {
	return token.LocationBetween(&n.StartToken, &n.EndToken)
}

func (n *New) String() string {
	return fmt.Sprintf("New{Type:%s, ArrayLengthExpr:%s, TailLengthExpr:%s}", n.Type.String(), n.ArrayLengthExpr.String(), n.TailLengthExpr.String())
}

func (n *New) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	var err io.Error
	n.StartToken, err = p.Consume(token.TOK_NEW)
	if err != nil {
		return nil, err
	}

	if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
		return nil, err
	}

	n.Type, err = p.ParseType()
	if err != nil {
		return nil, err
	}

	if p.Match(token.TOK_COMMA) {
		n.ArrayLengthExpr, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
		if p.Match(token.TOK_COMMA) {
			n.TailLengthExpr = n.ArrayLengthExpr
			n.ArrayLengthExpr, err = p.ParseExpr()
			if err != nil {
				return nil, err
			}
		}
	}

	n.EndToken, err = p.Consume(token.TOK_RIGHTPAREN)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *New) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	panic("implement me")
}
