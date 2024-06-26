package type_

import (
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"

	"go.uber.org/zap"
)

// Lookup represents a declaration lookup
type Lookup struct {
	Constant bool
	Context_ ast.TypeContext
	Tokens   []token.Token
}

func (l *Lookup) IsConstant() bool {
	return l.Constant
}
func (l *Lookup) SetConstant() {
	l.Constant = true
}

func (l *Lookup) Name() token.Token {
	return l.Tokens[0]
}

func (l *Lookup) Context() ast.TypeContext {
	return l.Context_
}

func (l *Lookup) Location() *token.Location {
	return token.LocationBetween(&l.Tokens[0], &l.Tokens[len(l.Tokens)-1])
}

func (l *Lookup) String() string {
	var tokensString, spacer string
	for _, t := range l.Tokens {
		tokensString += spacer + t.Value
		spacer = "."
	}
	return tokensString
}

func (l *Lookup) Declaration(_ ast.SemanticParser) (ast.Declaration, io.Error) {
	d, err := l.Context_.Declaration(l.Tokens)
	if err != nil {
		return nil, err
	} else if len(d.Generics()) != 0 {
		return nil, io.NewError("expected generic declaration for lookup type",
			zap.Stringer("type", l),
			zap.Stringer("location", l.Location()),
		)
	}

	return d, nil
}

func (l *Lookup) Concretize(mapping []ast.Type) ast.Type {
	return l
}

func (l *Lookup) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	decl, err := l.Declaration(p)
	if err != nil {
		return false, err
	}
	return decl.ExtendsAsPointer(p, parent)
}

func (l *Lookup) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return l.Equals(p, parent)
}

func (l *Lookup) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherComposite, ok := other.(*Lookup); ok {
		cDeclaration, err := l.Declaration(p)
		if err != nil {
			return false, err
		}

		otherDeclaration, err := otherComposite.Declaration(p)
		if err != nil {
			return false, err
		}

		return cDeclaration == otherDeclaration, nil
	} else if otherDeclaration, ok := other.(ast.Declaration); ok {
		cDeclaration, err := l.Declaration(p)
		if err != nil {
			return false, err
		}

		return cDeclaration == otherDeclaration, nil
	}

	return false, nil
}

func (l *Lookup) Syntax(p ast.SyntaxParser) (ast.DeclType, io.Error) {
	tok, err := p.Consume(token.TOK_IDENTIFIER)
	if err != nil {
		return nil, err
	}

	l.Tokens = append(l.Tokens, tok)
	for p.Match(token.TOK_PERIOD) {
		tok, err := p.Consume(token.TOK_IDENTIFIER)
		if err != nil {
			return nil, err
		}

		l.Tokens = append(l.Tokens, tok)
	}

	return l, nil
}

func (l *Lookup) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	d, _ := l.Declaration(nil)

	tid, err := d.TypeID(parser)
	if err != nil {
		return nil, err
	}

	index := tid.Index()
	if l.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, tid.ListIndex()), nil
}

func (l *Lookup) IsConcrete() bool {
	d, _ := l.Declaration(nil)
	return d.IsConcrete()
}

func (*Lookup) IsFieldable() bool {
	return true
}
