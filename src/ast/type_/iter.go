package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Iter struct {
	context ast.TypeContext
	Type    ast.Type
}

func (i *Iter) Context() ast.TypeContext {
	return i.context
}

func (i *Iter) String() string {
	return fmt.Sprintf("Iter{Type:%s}", i.Type.String())
}

func (i *Iter) Location() token.Location {
	return i.Type.Location()
}

func (i *Iter) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (i *Iter) Extends(parent ast.Type) (bool, io.Error) {
	return i.Equals(parent)
}

func (i *Iter) Equals(other ast.Type) (bool, io.Error) {
	if otherThread, ok := other.(*Thread); ok {
		return i.Type.Equals(otherThread.Type)
	}

	return false, nil
}
