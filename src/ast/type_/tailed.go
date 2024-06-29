package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/compile/syntax/typeid"
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
	if otherTailed, ok := other.(*Tailed); ok && otherTailed.Size != -1 && t.Size != otherTailed.Size {
		return false, nil
	}

	d, err := t.Declaration(p)
	if err != nil {
		return false, err
	}

	return d.ExtendsAsPointer(p, other)
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

func (t *Tailed) TypeID(parser ast.SemanticParser) (ast.TypeID, io.Error) {
	tid, err := t.Type.TypeID(parser)
	if err != nil {
		return nil, err
	}

	if t.Size == -1 {
		return tid, nil
	}

	lid, err := parser.Types().NextListIndex([]uint64{uint64(tid.ListIndex()), uint64(t.Size)})
	if err != nil {
		return nil, err
	}

	index := typeid.ID_Thread.Index()
	if t.Constant {
		index |= 1 << 31
	}
	return typeid.NewTypeID(index, lid), nil
}

func (t *Tailed) IsConcrete() bool {
	return t.Type.IsConcrete()
}

func (t *Tailed) IsFieldable() bool {
	return t.Size == -1 && t.Type.IsFieldable()
}
