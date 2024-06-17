package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Composite struct {
	Context_ ast.TypeContext
	Tokens   []token.Token
}

func (c *Composite) Name() token.Token {
	return c.Tokens[0]
}

func (c *Composite) Context() ast.TypeContext {
	return c.Context_
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

func (c *Composite) Syntax(p ast.SyntaxParser) (ast.DeclType, io.Error) {
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

	return c, nil
}

func (c *Composite) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	if ok, err := c.Equals(parent); err != nil || !ok {
		return ok, err
	}
	pComposite, ok := parent.(*Composite)
	if ok {
		return false, nil
	}

	cDecl, err := c.Declaration()
	if err != nil {
		return false, err
	}
	cChildDecl, ok := cDecl.(ast.ChildDeclaration)
	if !ok {
		return false, nil
	}

	for _, parentType := range cChildDecl.Parents() {
		if ok, err := parentType.ExtendsAsPointer(pComposite); err != nil || ok {
			return ok, err
		}
	}

	return false, nil
}

func (c *Composite) Extends(parent ast.Type) (bool, io.Error) {
	return c.Equals(parent)
}

func (c *Composite) Equals(other ast.GenericTypeArg) (bool, io.Error) {
	if otherComposite, ok := other.(*Composite); ok {
		var cDeclaration, otherDeclaration ast.Declaration
		var err io.Error
		if cDeclaration, err = c.Declaration(); err != nil {
			return false, err
		} else if otherDeclaration, err = otherComposite.Declaration(); err != nil {
			return false, err
		}

		return cDeclaration == otherDeclaration, nil
	}

	return false, nil
}

func (c *Composite) Declaration() (ast.Declaration, io.Error) {
	return c.Context_.Declaration(c.Tokens)
}
