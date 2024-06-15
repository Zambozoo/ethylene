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

func (s Super) ExtendsAsPointer(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	return s.Equals(ctx, parent)
}

func (s Super) Extends(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	for _, t := range s {
		if extends, err := ctx.Extends(t, parent); err != nil || !extends {
			return false, err
		}
	}

	return true, nil
}

func (s Super) containsType(ctx ast.TypeContext, t ast.Type) (bool, io.Error) {
	for _, t2 := range s {
		if ok, err := t2.Equals(ctx, t); err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (s Super) Equals(ctx ast.TypeContext, parent ast.Type) (bool, io.Error) {
	if parentSuper, ok := parent.(Super); ok {
		if len(s) != len(parentSuper) {
			return false, nil
		}

		for _, t := range s {
			if ok, err := parentSuper.containsType(ctx, t); err != nil || !ok {
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
