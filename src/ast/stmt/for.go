package stmt

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// For1 represents a 'for(condition){}' loop
type For1 struct {
	StartToken token.Token
	Condition  ast.Expression
	Then       ast.Statement
	Else       ast.Statement
}

func (f *For1) Location() token.Location {
	endLocatable := f.Then
	if f.Else != nil {
		endLocatable = f.Else
	}
	return token.LocationBetween(&f.StartToken, endLocatable)
}

func (f *For1) String() string {
	var elseString string
	if f.Else != nil {
		elseString = fmt.Sprintf("\nelse%s", f.Else.String())
	}
	return fmt.Sprintf("for(%s)\n%s%s",
		f.Condition.String(),
		f.Then.String(),
		elseString,
	)
}

func (f *For0) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	p.Scope().Wrap()
	defer p.Scope().Unwrap()
	if err := p.Scope().AddLabel(&emptyLabel); err != nil {
		return nil, err
	}

	return f.Stmt.Semantic(p)
}

// For1 represents a 'for{}' loop
type For0 struct {
	StartToken token.Token
	Stmt       ast.Statement
}

func (f *For0) Location() token.Location {
	return token.LocationBetween(&f.StartToken, f.Stmt)
}

func (f *For0) String() string {
	return fmt.Sprintf("for\n%s", f.Stmt.String())
}

func (f *For1) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	p.Scope().Wrap()
	defer p.Scope().Unwrap()
	if err := p.Scope().AddLabel(&emptyLabel); err != nil {
		return nil, err
	}

	t, err := f.Condition.Semantic(p)
	if err != nil {
		return nil, err
	} else if _, err := type_.MustExtend(p, t, &type_.Boolean{}); err != nil {
		return nil, err
	}

	returnType, err := f.Then.Semantic(p)
	if err != nil {
		return nil, err
	}

	if f.Else != nil {
		elseReturnType, err := f.Else.Semantic(p)
		if err != nil {
			return nil, err
		}
		returnType = type_.Join(returnType, elseReturnType)
	}

	return returnType, nil
}

func parseFor(p ast.SyntaxParser) (ast.Statement, io.Error) {
	startToken, err := p.Consume(token.TOK_FOR)
	if err != nil {
		return nil, err
	}

	if p.Peek().Type != token.TOK_LEFTPAREN {
		stmt, err := p.ParseStmt()
		if err != nil {
			return nil, err
		}

		return &For0{
			StartToken: startToken,
			Stmt:       stmt,
		}, nil
	}

	if _, err := p.Consume(token.TOK_LEFTPAREN); err != nil {
		return nil, err
	}

	var initialization ast.Statement
	if t := p.Peek(); t.Type != token.TOK_SEMICOLON && t.Type != token.TOK_VAR {
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}

		if !p.Match(token.TOK_SEMICOLON) {
			if _, err := p.Consume(token.TOK_RIGHTPAREN); err != nil {
				return nil, err
			}

			then, err := p.ParseStmt()
			if err != nil {
				return nil, err
			}

			var else_ ast.Statement
			if p.Match(token.TOK_ELSE) {
				else_, err = p.ParseStmt()
				if err != nil {
					return nil, err
				}
			}

			return &For1{
				StartToken: startToken,
				Condition:  expr,
				Then:       then,
				Else:       else_,
			}, nil
		} else {
			initialization = &Expr{
				Expr:     expr,
				EndToken: p.Prev(),
			}
		}
	} else if !p.Match(token.TOK_SEMICOLON) {
		initialization, err = p.ParseStmt()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expression
	if p.Peek().Type != token.TOK_SEMICOLON {
		condition, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.Consume(token.TOK_SEMICOLON); err != nil {
		return nil, err
	}

	var increment ast.Expression
	if p.Peek().Type != token.TOK_RIGHTPAREN {
		increment, err = p.ParseExpr()
		if err != nil {
			return nil, err
		}
	}
	incrementEndToken := p.Prev()

	if _, err := p.Consume(token.TOK_RIGHTPAREN); err != nil {
		return nil, err
	}

	thenStartToken := p.Peek()
	then, err := p.ParseStmt()
	if err != nil {
		return nil, err
	}

	var else_ ast.Statement
	if condition != nil && p.Match(token.TOK_ELSE) {
		else_, err = p.ParseStmt()
		if err != nil {
			return nil, err
		}
	}

	var thenStmt ast.Statement
	if increment != nil {
		thenStmt = &Block{
			BoundedStmt: BoundedStmt{
				StartToken: thenStartToken,
				EndToken:   p.Prev(),
			},
			Stmts: []ast.Statement{
				then,
				&Expr{Expr: increment, EndToken: incrementEndToken},
			},
		}
	} else {
		thenStmt = then
	}

	var forStmts []ast.Statement
	if initialization != nil {
		forStmts = append(forStmts, initialization)
	}

	if condition != nil {
		forStmts = append(forStmts,
			&For1{
				StartToken: startToken,
				Condition:  condition,
				Then:       thenStmt,
				Else:       else_,
			},
		)
	} else {
		forStmts = append(forStmts,
			&For0{
				StartToken: startToken,
				Stmt:       thenStmt,
			},
		)
	}

	return &Block{
		BoundedStmt: BoundedStmt{
			StartToken: startToken,
			EndToken:   p.Prev(),
		},
		Stmts: forStmts,
	}, nil

}
