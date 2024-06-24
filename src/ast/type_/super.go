package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

type Super []ast.Type

func (s Super) Location() token.Location {
	if len(s) == 0 {
		panic("empty super")
	}

	return token.LocationBetween(s[0], s[len(s)-1])
}

func (s Super) String() string {
	return fmt.Sprintf("Super{Types:%s}", strs.Strings(s))
}

func (s Super) Key() string {
	var str string
	var spacer string
	for _, t := range s {
		str += spacer + t.Key()
		spacer = ", "
	}
	return fmt.Sprintf("<:[%s]", str)
}

func (s Super) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	return s.Equals(parent)
}

func (s Super) Extends(parent ast.Type) (bool, io.Error) {
	for _, t := range s {
		if extends, err := t.Extends(parent); err != nil || !extends {
			return false, err
		}
	}

	return true, nil
}

func (s Super) containsType(t ast.Type) (bool, io.Error) {
	for _, t2 := range s {
		if ok, err := t2.Equals(t); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (s Super) Equals(parent ast.Type) (bool, io.Error) {
	if parentSuper, ok := parent.(Super); ok {
		if len(s) != len(parentSuper) {
			return false, nil
		}

		for _, t := range s {
			if ok, err := parentSuper.containsType(t); err != nil || !ok {
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

	var returnType Super
	for _, t := range ts {
		if t == nil {
			continue
		}

		if s, ok := t.(Super); ok {
			t = Join(s...)
		}

		returnType = append(returnType, t)
	}

	return returnType
}
