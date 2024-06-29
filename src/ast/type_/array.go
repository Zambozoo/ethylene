package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Array represents an array of elements of another type
type Array struct {
	Constant bool
	Type     ast.Type
	Size     int64
	EndToken token.Token
}

func (a *Array) Location() *token.Location {
	return token.LocationBetween(a.Type, &a.EndToken)
}

func (a *Array) String() string {
	return fmt.Sprintf("%s[%d]", a.Type.String(), a.Size)
}

func (a *Array) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return a.Equals(p, parent)
}

func (a *Array) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	if ptr, ok := parent.(*Pointer); ok {
		return a.Type.Extends(p, ptr.Type)
	}

	return false, nil
}

func (a *Array) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if arr, ok := other.(*Array); ok {
		if a.Size != arr.Size {
			return false, nil
		}

		return a.Type.Equals(p, arr.Type)
	}

	return false, nil
}

func (a *Array) Concretize(mapping []ast.Type) ast.Type {
	return &Array{
		Constant: a.Constant,
		Type:     a.Type.Concretize(mapping),
		Size:     a.Size,
		EndToken: a.EndToken,
	}
}

func (a *Array) IsConstant() bool {
	return a.Constant
}
func (a *Array) SetConstant() {
	a.Constant = true
}

func (a *Array) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	// TODO: FIGURE OUT SIZE
	panic("")
}

func (a *Array) IsConcrete() bool {
	return a.Type.IsConcrete()
}
