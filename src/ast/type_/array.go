package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Array represents an array of elements of another type
type Array struct {
	Type     ast.Type
	Size     int64
	EndToken token.Token
}

func (a *Array) Location() token.Location {
	return token.LocationBetween(a.Type, &a.EndToken)
}

func (a *Array) String() string {
	return fmt.Sprintf("Array{Type:%s,Size:%d}", a.Type.String(), a.Size)
}

func (a *Array) ExtendsAsPointer(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return a.Equals(ctx, parent)
}

func (a *Array) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	if ptr, ok := parent.(*Pointer); ok {
		return a.Type.Extends(ctx, ptr.Type)
	}

	return false, nil
}

func (a *Array) Equals(ctx ast.TypeContext, other ast.Type) (bool, io.Error) {
	if arr, ok := other.(*Array); ok {
		if a.Size != arr.Size {
			return false, nil
		}

		return a.Type.Equals(ctx, arr.Type)
	}

	return false, nil
}
