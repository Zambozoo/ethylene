package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Tailed struct {
	Type     ast.Type
	Size     int64
	EndToken token.Token
}

func (t *Tailed) Location() token.Location {
	return token.LocationBetween(t.Type, &t.EndToken)
}

func (t *Tailed) String() string {
	return fmt.Sprintf("Tailed{Type:%s,Size:%d}", t.Type.String(), t.Size)
}

func (t *Tailed) ExtendsAsPointer(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (t *Tailed) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return t.Equals(ctx, parent)
}

func (t *Tailed) Equals(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	if otherTailed, ok := other.(*Tailed); ok {
		if t.Size != otherTailed.Size {
			return false, nil
		}

		panic("not implemented")
	}

	return false, nil
}

func (t *Tailed) Declaration(ctx ast.TypeContext) (ast.Declaration, io.Error) {
	panic("not implemented")
}
