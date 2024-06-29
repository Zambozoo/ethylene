package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// TODO: REFERENCE ACTUAL TYPE IN STD LIB

type Thread struct {
	Constant bool
	Type     ast.Type
}

func (t *Thread) String() string {
	return fmt.Sprintf("Thread{Type:%s}", t.Type.String())
}

func (t *Thread) Key(p ast.SemanticParser) (string, io.Error) {
	k, err := t.Type.Key(p)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("thread[%s]", k), nil
}

func (t *Thread) Location() token.Location {
	return t.Type.Location()
}

func (t *Thread) ExtendsAsPointer(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (t *Thread) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return t.Equals(p, parent)
}

func (t *Thread) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherThread, ok := other.(*Thread); ok {
		return t.Type.Equals(p, otherThread.Type)
	}

	return false, nil
}

func (t *Thread) Concretize(mapping []ast.Type) ast.Type {
	return &Thread{
		Type: t.Type.Concretize(mapping),
	}
}

func (t *Thread) IsConstant() bool {
	return t.Constant
}
func (t *Thread) SetConstant() {
	t.Constant = true
}
