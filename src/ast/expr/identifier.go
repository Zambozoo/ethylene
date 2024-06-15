package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Identifier struct {
	token.Token
}

func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier{Value:%s}", i.Value)
}

func (i *Identifier) Syntax(p ast.SyntaxParser) (ast.Expression, io.Error) {
	var err io.Error
	i.Token, err = p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Identifier) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	v, _, err := p.Scope().GetVariable(i.Token)
	if err != nil {
		return nil, err
	}

	return v.Type(), nil
}
