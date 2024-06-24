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

func (t *Token) Location() Location {
	return t.Loc
}

func (t *Token) String() string {
	var value string
	switch t.Type {
	case TOK_CHARACTER:
		value = string(t.Rune)
	case TOK_INTEGER:
		value = fmt.Sprintf("%d", t.Integer)
	case TOK_FLOAT:
		value = fmt.Sprintf("%f", t.Float)
	}

	return fmt.Sprintf("Token{Type:%s, Value:%s, Location:%s}", t.Type.String(), value, &t.Loc)
}
