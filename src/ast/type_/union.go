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
	return strs.Strings(u, "`")
}

func (u Union) Key(p ast.SemanticParser) (string, io.Error) {
	var s string
	var spacer string
	for _, t := range u {
		k, err := t.Key(p)
		if err != nil {
			return "", err
		}
		s += spacer + k
		spacer = ","
	}
	return fmt.Sprintf("(%s)", s), nil
}

func (u Union) ExtendsAsPointer(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	return u.Equals(p, parent)
}

// Extends returns true if all
func (u Union) Extends(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	for _, t := range u {
		if extends, err := t.Extends(p, parent); err != nil || !extends {
			return false, err
		}
	}

	return true, nil
}

func containsType(p ast.SemanticParser, ts []ast.Type, t ast.Type) (bool, io.Error) {
	for _, t2 := range ts {
		if ok, err := t2.Equals(p, t); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (u Union) Equals(p ast.SemanticParser, parent ast.Type) (bool, io.Error) {
	if parentUnion, ok := parent.(Union); ok {
		if len(u) != len(parentUnion) {
			return false, nil
		}

		for _, t := range u {
			if ok, err := containsType(p, parentUnion, t); err != nil || !ok {
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

func (u Union) Concretize(mapping []ast.Type) ast.Type {
	concreteTypes := make(Union, len(u))
	for i, t := range u {
		concreteTypes[i] = t.Concretize(mapping)
	}

	return concreteTypes
}

func (Union) IsConstant() bool {
	return false
}
func (Union) SetConstant() {
	panic("unreachable ")
}
