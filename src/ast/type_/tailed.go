package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

type Tailed struct {
	Constant bool
	Type     ast.DeclType
	Size     int64
	EndToken token.Token
}

func (t *Tailed) Name() token.Token {
	return t.Type.Name()
}

func (t *Tailed) Context() ast.TypeContext {
	return t.Type.Context()
}

func (t *Tailed) Location() *token.Location {
	return token.LocationBetween(t.Type, &t.EndToken)
}

func (t *Tailed) String() string {
	var tail string
	if t.Size >= 0 {
		tail = fmt.Sprintf("%d", t.Size)
	}
	return fmt.Sprintf("%s~%s", t.Type.String(), tail)
}

func (t *Tailed) ExtendsAsPointer(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	panic("not implemented")
}

func (t *Tailed) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return t.Equals(p, parent)
}

func (t *Tailed) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	if otherTailed, ok := other.(*Tailed); ok {
		if t.Size != otherTailed.Size {
			return false, nil
		}

		return t.Type.Equals(p, other)
	}

	return false, nil
}

func (t *Tailed) Declaration(p ast.SemanticParser) (ast.Declaration, io.Error) {
	return t.Type.Declaration(p)
}

func (t *Tailed) Concretize(mapping []ast.Type) ast.Type {
	return &Tailed{
		Constant: t.Constant,
		Type:     t.Type.Concretize(mapping).(ast.DeclType),
		Size:     t.Size,
	}
}

func (t *Tailed) IsConstant() bool {
	return t.Constant
}
func (t *Tailed) SetConstant() {
	t.Constant = true
}
