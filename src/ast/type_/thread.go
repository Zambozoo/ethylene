package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Thread struct {
	Type ast.Type
}

func (t *Thread) String() string {
	return fmt.Sprintf("Thread{Type:%s}", t.Type.String())
}

func (t *Thread) Location() token.Location {
	return t.Type.Location()
}

func (t *Thread) ExtendsAsPointer(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (t *Thread) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return t.Equals(ctx, parent)
}

func (t *Thread) Equals(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	if otherThread, ok := other.(*Thread); ok {
		return t.Type.Equals(ctx, otherThread.Type)
	}

	return false, nil
}
