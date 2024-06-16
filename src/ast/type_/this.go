package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

type Sub []ast.Type

func (s Sub) Location() token.Location {
	if len(s) == 0 {
		panic("empty sub")
	}

	return token.LocationBetween(s[0], s[len(s)-1])
}

func (s Sub) String() string {
	return fmt.Sprintf("Sub{Types:%s}", strs.Strings(s))
}

func (s Sub) ExtendsAsPointer(parent ast.Type) (bool, io.Error) {
	return s.Equals(parent)
}

func (s Sub) Extends(parent ast.Type) (bool, io.Error) {
	for _, t := range s {
		if extends, err := t.Extends(parent); err != nil || !extends {
			return false, err
		}
	}

	return true, nil
}

func (s Sub) containsType(t ast.Type) (bool, io.Error) {
	for _, t2 := range s {
		if ok, err := t2.Equals(t); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (s Sub) Equals(parent ast.Type) (bool, io.Error) {
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
