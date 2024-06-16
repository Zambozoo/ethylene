package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

// FunctionType represents a function signature
type Function struct {
	ReturnType_     ast.Type
	ParameterTypes_ []ast.Type
	EndToken        token.Token
}

func (f *Function) Location() token.Location {
	return token.LocationBetween(f.ReturnType_, &f.EndToken)
}

func (f *Function) String() string {
	return fmt.Sprintf("Function{ReturnType:%s, ParameterTypes:%s}", f.ReturnType_, strs.Strings(f.ParameterTypes_))
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

func (f *Function) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	return f.Equals(parent)
}

func (f *Function) Extends(parent ast.Type) (bool, io.Error) {
	if fun, ok := parent.(*Function); ok {
		if f.Arity() != fun.Arity() {
			return false, nil
		}

		if ok, err := f.ReturnType_.Extends(fun.ReturnType_); err != nil || !ok {
			return false, err
		}

		for i, childArgType := range f.ParameterTypes_ {
			parentArgType := fun.ParameterTypes_[i]
			if ok, err := parentArgType.Extends(childArgType); err != nil || !ok {
				return false, err
			}
		}
	}

	return false, nil
}

func (f *Function) Equals(other ast.Type) (bool, io.Error) {
	fOther, ok := other.(*Function)
	if !ok || f.Arity() != fOther.Arity() {
		return false, nil
	} else if ok, err := f.ReturnType_.Equals(fOther.ReturnType_); err != nil || !ok {
		return ok, err
	}

	for i, childArgType := range f.ParameterTypes_ {
		parentArgType := fOther.ParameterTypes_[i]
		if ok, err := childArgType.Equals(parentArgType); err != nil || !ok {
			return ok, err
		}
	}

	return true, nil
}
