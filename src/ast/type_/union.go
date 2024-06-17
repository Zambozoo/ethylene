package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

// Union represents a type that could be any of the given types.
type Union []ast.Type

func (u Union) Location() token.Location {
	return token.LocationBetween(u[0], u[len(u)-1])
}

func (u Union) String() string {
	return fmt.Sprintf("Union{Types:%s}", strs.Strings(u))
}

func (u Union) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	return u.Equals(parent)
}

// Extends returns true if all
func (u Union) Extends(parent ast.Type) (bool, io.Error) {
	for _, t := range u {
		if extends, err := t.Extends(parent); err != nil || !extends {
			return false, err
		}
	}

	return true, nil
}

func containsType(ts []ast.Type, t ast.Type) (bool, io.Error) {
	for _, t2 := range ts {
		if ok, err := t2.Equals(t); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (u Union) Equals(parent ast.Type) (bool, io.Error) {
	if parentUnion, ok := parent.(Union); ok {
		if len(u) != len(parentUnion) {
			return false, nil
		}

		for _, t := range u {
			if ok, err := containsType(parentUnion, t); err != nil || !ok {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}

func Join(ts ...ast.Type) ast.Type {
	if len(ts) == 0 {
		return nil
	}

	var returnType Union
	for _, t := range ts {
		if t == nil {
			continue
		}

		if s, ok := t.(Union); ok {
			t = Join(s...)
		}

		returnType = append(returnType, t)
	}

	return returnType
}

func (u Union) Concretize(mapping map[string]ast.Type) ast.Type {
	concreteTypes := make(Union, len(u))
	for i, t := range u {
		concreteTypes[i] = t.Concretize(mapping)
	}

	return concreteTypes
}
