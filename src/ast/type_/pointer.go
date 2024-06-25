package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// PointerType represents a pointer to another type
type Pointer struct {
	Constant bool
	Type     ast.Type
	EndToken token.Token
}

func (p *Pointer) Location() token.Location {
	return token.LocationBetween(p.Type, &p.EndToken)
}

func (p *Pointer) String() string {
	return fmt.Sprintf("Pointer{Type:%s}", p.Type.String())
}

func (p *Pointer) Key() string {
	return fmt.Sprintf("%s*", p.Type.Key())
}

func (p *Pointer) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	return p.Equals(parent)
}

func (p *Pointer) Extends(parent ast.Type) (bool, io.Error) {
	if parentPtr, ok := parent.(*Pointer); ok {
		return p.Type.Extends(parentPtr.Type)
	}

	return false, nil
}

func (p *Pointer) Equals(other ast.Type) (bool, io.Error) {
	if otherPtr, ok := other.(*Pointer); ok {
		return p.Type.Equals(otherPtr.Type)
	}

	return false, nil
}

func (p *Pointer) Concretize(mapping map[string]ast.Type) ast.Type {
	return &Array{
		Constant: p.Constant,
		Type:     p.Type.Concretize(mapping),
		EndToken: p.EndToken,
	}
}

func (p *Pointer) IsConstant() bool {
	return p.Constant
}
func (p *Pointer) SetConstant() {
	p.Constant = true
}
