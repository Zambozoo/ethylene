package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

// Block represents a block of statements enclosed in braces
type Block struct {
	BoundedStmt
	Stmts []ast.Statement
}

func (b *Block) String() string {
	return fmt.Sprintf("Block{Stmts:%s}", strs.Strings(b.Stmts))
}

func (b *Block) Syntax(p ast.SyntaxParser) io.Error {
	var err io.Error
	b.BoundedStmt.StartToken, err = p.Consume(token.TOK_LEFTBRACE)
	if err != nil {
		return err
	}

	for !p.Match(token.TOK_RIGHTBRACE) {
		stmt, err := p.ParseStmt()
		if err != nil {
			return err
		}

		b.Stmts = append(b.Stmts, stmt)
	}
	b.BoundedStmt.EndToken = p.Prev()

	return nil
}

func (b *Block) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	p.Scope().Wrap()
	defer p.Scope().Unwrap()

	var returnType ast.Type
	for _, stmt := range b.Stmts {
		t, err := stmt.Semantic(p)
		if err != nil {
			return nil, err
		}
		returnType = type_.Join(returnType, t)
	}

	return returnType, nil
}
