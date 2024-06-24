package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// TODO: REFERENCE ACTUAL TYPE IN STD LIB

type Thread struct {
	Type ast.Type
}

func (t *Thread) String() string {
	return fmt.Sprintf("Thread{Type:%s}", t.Type.String())
}

func (t *Thread) Key() string {
	return fmt.Sprintf("Thread[%s]", t.Type.Key())
}

func (t *Thread) Location() token.Location {
	return t.Type.Location()
}

func (t *Thread) ExtendsAsPointer(other ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (t *Thread) Extends(parent ast.Type) (bool, io.Error) {
	return t.Equals(parent)
}

func (t *Thread) Equals(other ast.Type) (bool, io.Error) {
	if otherThread, ok := other.(*Thread); ok {
		return t.Type.Equals(otherThread.Type)
	}

	return false, nil
}
