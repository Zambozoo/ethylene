package token

import "fmt"

type Token struct {
	Type Type

	Rune    rune
	Value   string
	Integer uint64
	Float   float64

	Loc Location
}

func (t *Token) Location() *Location {
	return &t.Loc
}

func (t *Token) String() string {
	switch t.Type {
	case TOK_CHARACTER:
		return fmt.Sprintf("'%s'", string(t.Rune))
	case TOK_INTEGER:
		return fmt.Sprintf("%d", t.Integer)
	case TOK_FLOAT:
		return fmt.Sprintf("%f", t.Float)
	case TOK_IDENTIFIER:
		return t.Value
	case TOK_STRING:
		return fmt.Sprintf("%q", t.Value)
	}

	return t.Type.String()
}
