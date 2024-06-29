package ast

import (
	"geth-cody/compile/lexer/token"
	"geth-cody/io"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type Modifier string

const (
	MOD_PUBLIC    Modifier = "public"
	MOD_PRIVATE   Modifier = "private"
	MOD_PROTECTED Modifier = "protected"

	MOD_STATIC  Modifier = "static"
	MOD_VIRTUAL Modifier = "virtual"
	MOD_NATIVE  Modifier = "native"
)

func (m Modifier) String() string {
	return string(m)
}

func SyntaxModifiers(p SyntaxParser) (map[Modifier]struct{}, io.Error) {
	startT := p.Peek()
	modifiers := map[Modifier]struct{}{}

	for {
		switch t := p.Peek(); t.Type {
		case token.TOK_PUBLIC, token.TOK_PRIVATE, token.TOK_PROTECTED, token.TOK_STATIC, token.TOK_NATIVE, token.TOK_VIRTUAL:
			p.Next()
			m := Modifier(t.Type)
			if _, ok := modifiers[m]; ok {
				return nil, io.NewError("duplicate modifier", zap.String("modifier", string(m)))
			}
			modifiers[Modifier(t.Type)] = struct{}{}
		default:
			accessModiferCount := 0
			for _, m := range []Modifier{MOD_PUBLIC, MOD_PRIVATE, MOD_PROTECTED} {
				if _, ok := modifiers[m]; ok {
					accessModiferCount += 1
				}
			}
			if accessModiferCount > 1 {
				endT := p.Prev()
				return nil, io.NewError("multiple access modifiers",
					zap.Stringers("modifiers", maps.Keys(modifiers)),
					zap.Stringer("location", token.LocationBetween(&startT, &endT)),
				)
			}

			_, isStatic := modifiers[MOD_STATIC]
			_, isNative := modifiers[MOD_NATIVE]
			if _, isVirtual := modifiers[MOD_VIRTUAL]; isVirtual && (isStatic || isNative) {
				endT := p.Prev()

				return nil, io.NewError("virtual methods cannot be static or native",
					zap.Stringers("modifiers", maps.Keys(modifiers)),
					zap.Stringer("location", token.LocationBetween(&startT, &endT)),
				)
			}

			return modifiers, nil
		}
	}
}
