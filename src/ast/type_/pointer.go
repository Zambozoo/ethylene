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

func (p *Pointer) Key(parser ast.SemanticParser) (string, io.Error) {
	k, err := p.Type.Key(parser)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s*", k), nil
}

func (p *Pointer) ExtendsAsPointer(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return p.Equals(parser, parent)
}

func (p *Pointer) Extends(parser ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	if parentPtr, ok := parent.(*Pointer); ok {
		return p.Type.Extends(parser, parentPtr.Type)
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
