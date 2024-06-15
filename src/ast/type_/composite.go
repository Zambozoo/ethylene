package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"math"

	"go.uber.org/zap"
)

type Composite struct {
	Tokens []token.Token
}

func (c *Composite) Location() token.Location {
	return token.LocationBetween(&c.Tokens[0], &c.Tokens[len(c.Tokens)-1])
}

func (c *Composite) String() string {
	var tokensString, spacer string
	for _, t := range c.Tokens {
		tokensString += spacer + t.Value
		spacer = "."
	}
	return fmt.Sprintf("Composite{Tokens:%s}", tokensString)
}

func (c *Composite) Syntax(p ast.SyntaxParser) (ast.Type, io.Error) {
	tok, err := p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	c.Tokens = append(c.Tokens, tok)
	for p.Match(token.TOK_PERIOD) {
		tok, err := p.Consume(token.TOK_IDENTIFIER)
		if err != nil {
			return nil, err
		}

		c.Tokens = append(c.Tokens, tok)
	}

	var t ast.Type = c
	if p.Match(token.TOK_TILDE) {
		size := int64(-1)
		if p.Match(token.TOK_INTEGER) {
			tok := p.Prev()
			if tok.Integer > math.MaxInt {
				return nil, io.NewError("type tail size is larger than max signed integer limit", zap.Any("token", tok))
			}
			size = int64(tok.Integer)
		}
		t = &Tailed{Type: t, Size: size, EndToken: p.Prev()}
	}

	return t, nil
}

func (c *Composite) ExtendsAsPointer(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (c *Composite) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return c.Equals(ctx, parent)
}

func (c *Composite) Equals(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	if otherComposite, ok := other.(*Composite); ok {
		var cDeclaration, otherDeclaration ast.Declaration
		var err io.Error
		if cDeclaration, err = c.Declaration(ctx); err != nil {
			return false, err
		} else if otherDeclaration, err = otherComposite.Declaration(ctx); err != nil {
			return false, err
		}

		return cDeclaration == otherDeclaration, nil
	}

	return false, nil
}

func (c *Composite) Declaration(ctx ast.TypeContext) (ast.Declaration, io.Error) {
	return ctx.Declaration(c.Tokens)
}
