package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Tailed struct {
	Type     ast.DeclType
	Size     int64
	EndToken token.Token
}

func (t *Tailed) Location() token.Location {
	return token.LocationBetween(t.Type, &t.EndToken)
}

func (t *Tailed) String() string {
	return fmt.Sprintf("Tailed{Type:%s,Size:%d}", t.Type.String(), t.Size)
}

func (t *Tailed) ExtendsAsPointer(other ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (t *Tailed) Extends(parent ast.Type) (bool, io.Error) {
	return t.Equals(parent)
}

func (t *Tailed) Equals(other ast.Type) (bool, io.Error) {
	if otherTailed, ok := other.(*Tailed); ok {
		if t.Size != otherTailed.Size {
			return false, nil
		}

		return t.Type.Equals(other)
	}

	return false, nil
}

func (t *Tailed) Declaration() (ast.Declaration, io.Error) {
	return t.Type.Declaration()
}
