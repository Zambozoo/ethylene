package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// TODO: REFERENCE ACTUAL TYPE IN STD LIB

type Iter struct {
	Constant bool
	Type     ast.Type
}

func (i *Iter) String() string {
	return fmt.Sprintf("Iter{Type:%s}", i.Type.String())
}

func (i *Iter) Key(p ast.SemanticParser) (string, io.Error) {
	k, err := i.Type.Key(p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Iter[%s]", k), nil
}

func (i *Iter) Location() token.Location {
	return i.Type.Location()
}

func (i *Iter) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (i *Iter) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return i.Equals(p, parent)
}

func (i *Iter) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherThread, ok := other.(*Iter); ok {
		return i.Type.Equals(p, otherThread.Type)
	}

	return false, nil
}

func (i *Iter) IsConstant() bool {
	return i.Constant
}
func (i *Iter) SetConstant() {
	i.Constant = true
}

func (i *Iter) Concretize(mapping []ast.Type) ast.Type {
	return &Iter{
		Type: i.Type.Concretize(mapping),
	}
}
