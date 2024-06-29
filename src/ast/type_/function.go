package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/stringers"
)

// FunctionType represents a function signature
type Function struct {
	Constant        bool
	ReturnType_     ast.Type
	ParameterTypes_ []ast.Type
	EndToken        token.Token
}

func (f *Function) Location() *token.Location {
	return token.LocationBetween(f.ReturnType_, &f.EndToken)
}

func (f *Function) String() string {
	return fmt.Sprintf("%s(%s)", f.ReturnType_, stringers.Join(f.ParameterTypes_, ","))
}

func (f *Function) ReturnType() ast.Type {
	return f.ReturnType_
}
func (f *Function) ParameterTypes() []ast.Type {
	return f.ParameterTypes_
}
func (f *Function) Arity() int {
	return len(f.ParameterTypes_)
}

func (f *Function) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return f.Equals(p, parent)
}

func (f *Function) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	if fun, ok := parent.(*Function); ok {
		if f.Arity() != fun.Arity() {
			return false, nil
		}

		if ok, err := f.ReturnType_.Extends(p, fun.ReturnType_); err != nil || !ok {
			return false, err
		}

		for i, childArgType := range f.ParameterTypes_ {
			parentArgType := fun.ParameterTypes_[i]
			if ok, err := parentArgType.Extends(p, childArgType); err != nil || !ok {
				return false, err
			}
		}
	}

	return true, nil
}

func (f *Function) Equals(p ast.SemanticParser, other ast.Type) (bool, io.Error) {
	fOther, ok := other.(*Function)
	if !ok || f.Arity() != fOther.Arity() {
		return false, nil
	} else if ok, err := f.ReturnType_.Equals(p, fOther.ReturnType_); err != nil || !ok {
		return ok, err
	}

	for i, childArgType := range f.ParameterTypes_ {
		parentArgType := fOther.ParameterTypes_[i]
		if ok, err := childArgType.Equals(p, parentArgType); err != nil || !ok {
			return ok, err
		}
	}

	return true, nil
}

func (f *Function) Concretize(mapping []ast.Type) ast.Type {
	concreteParamTypes := make([]ast.Type, len(f.ParameterTypes_))
	for i, t := range f.ParameterTypes_ {
		concreteParamTypes[i] = t.Concretize(mapping)
	}

	return &Function{
		Constant:        f.Constant,
		ReturnType_:     f.ReturnType_.Concretize(mapping),
		ParameterTypes_: concreteParamTypes,
		EndToken:        f.EndToken,
	}
}

func (f *Function) IsConstant() bool {
	return f.Constant
}
func (f *Function) SetConstant() {
	f.Constant = true
}
