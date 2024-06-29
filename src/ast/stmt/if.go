package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// If represents an 'if' statement
type If struct {
	StartToken token.Token
	Condition  ast.Expression
	Then       ast.Statement
	Else       ast.Statement
}

func (i *If) Location() token.Location {
	endLocatable := i.Then
	if i.Else != nil {
		endLocatable = i.Else
	}
	return token.LocationBetween(&i.StartToken, endLocatable)
}

func (i *If) String() string {
	return fmt.Sprintf("If{Condition:%s,Then:%s,Else:%s}", i.Condition.String(), i.Then.String(), i.Else.String())
}

func (i *If) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	i.StartToken, err = p.Consume(token.TOK_IF)
	if err != nil {
		return err
	}

	if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
		return err
	}

	i.Condition, err = p.ParseExpr()
	if err != nil {
		return err
	}

	if _, err := p.Consume(token.TOK_RIGHTPAREN); err != nil {
		return err
	}

	i.Then, err = p.ParseStmt()
	if err != nil {
		return err
	}

	if p.Match(token.TOK_ELSE) {
		i.Else, err = p.ParseStmt()
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *If) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	t, err := i.Condition.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(p, t, &type_.Boolean{}); err != nil {
		return nil, err
	}

	returnType, err := i.Then.Semantic(p)
	if err != nil {
		return nil, err
	}

	if i.Else != nil {
		elseReturnType, err := i.Else.Semantic(p)
		if err != nil {
			return nil, err
		}
		returnType = type_.Join(returnType, elseReturnType)
	}

	return returnType, nil
}
