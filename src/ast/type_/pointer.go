package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
	"geth-cody/io"
	"math"
)

// PointerType represents a pointer to another type
type Pointer struct {
	Constant bool
	Type     ast.Type
	EndToken token.Token
}

func (p *Pointer) Location() *token.Location {
	return token.LocationBetween(p.Type, &p.EndToken)
}

func (p *Pointer) String() string {
	return fmt.Sprintf("%s*", p.Type.String())
}

func (p *Pointer) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}

func (p *Pointer) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	if parentPtr, ok := parent.(*Pointer); ok {
		return p.Type.ExtendsAsPointer(parser, parentPtr.Type)
	}

	return false, nil
}

func (p *Pointer) Equals(parser ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherPtr, ok := other.(*Pointer); ok {
		return p.Type.Equals(parser, otherPtr.Type)
	}

	return false, nil
}

func (p *Pointer) Concretize(mapping []ast.Type) ast.Type {
	return &Pointer{
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

func (p *Pointer) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	var t ast.Type = p
	for depth := 0; depth < math.MaxInt8; depth++ {
		if p, ok := t.(*Pointer); ok {
			t = p.Type
			continue
		}
		typeID, err := t.TypeID(parser)
		if err != nil {
			return nil, err
		}

		index := uint32(depth<<28) | typeID.Index()
		if p.Constant {
			index |= 1 << 31
		}

		return typeid.NewTypeID(index, typeID.ListIndex()), nil
	}

	return nil, io.NewError("pointers have a depth limit of 2^28-1")
}

func (p *Pointer) IsConcrete() bool {
	return p.Type.IsConcrete()
}
