package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

// Super is a type constraint for a generic type.
//
//	A :> [B, C, ...]
type Super struct {
	Type  ast.Type
	Types []ast.Type
}

func (s Super) Location() token.Location {
	return token.LocationBetween(s.Type, s.Types[len(s.Types)-1])
}

func (s Super) String() string {
	return fmt.Sprintf("Super{Type:%s,Types:%s}", s.Type.String(), strs.Strings(s.Types))
}

func (s Super) Equals(parent ast.GenericTypeArg) (bool, io.Error) {
	if parentSuper, ok := parent.(Super); ok {
		if len(s.Types) != len(parentSuper.Types) {
			return false, nil
		}

		for _, t := range s.Types {
			if ok, err := containsType(parentSuper.Types, t); err != nil || !ok {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}
