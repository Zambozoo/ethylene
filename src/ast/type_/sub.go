package type_

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
	"geth-cody/strs"
)

// Sub is a type constraint for a generic type.
//
//	A <: [B, C, ...]
type Sub struct {
	Type  ast.Type
	Types []ast.Type
}

func (s Sub) Location() token.Location {
	return token.LocationBetween(s.Type, s.Types[len(s.Types)-1])
}

func (s Sub) String() string {
	return fmt.Sprintf("Sub{Type:%s,Types:%s}", s.Type.String(), strs.Strings(s.Types))
}

func (s Sub) Equals(parent ast.GenericTypeArg) (bool, io.Error) {
	if parentSub, ok := parent.(Sub); ok {
		if len(s.Types) != len(parentSub.Types) {
			return false, nil
		}

		for _, t := range s.Types {
			if ok, err := containsType(parentSub.Types, t); err != nil || !ok {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}
