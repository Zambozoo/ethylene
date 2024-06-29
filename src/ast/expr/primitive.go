package expr

import (
	"fmt"
	"geth-cody/ast"
	"geth-cody/ast/type_"
	"geth-cody/compile/lexer/token"
	"geth-cody/io"
)

// Integer represents valid integer literals
type Integer struct{ token.Token }

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Integer)
}
func (i *Integer) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.Integer{}, nil
}

// Float represents  valid floating point literals
type Float struct{ token.Token }

func (f *Float) String() string {
	return fmt.Sprintf("%f", f.Float)
}
func (f *Float) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.Float{}, nil
}

// Character represents valid character literals
type Character struct{ token.Token }

func (c *Character) String() string {
	return fmt.Sprintf("%c", c.Rune)
}
func (c *Character) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.Character{}, nil
}

// String represents valid string literals
type String struct{ token.Token }

func (s *String) String() string {
	return s.Value
}
func (s *String) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.String{}, nil
}

// True represents the `true` boolean
type True struct{ token.Token }

func (t *True) String() string {
	return "true"
}
func (t *True) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.Boolean{}, nil
}

// False represents the `false` boolean
type False struct{ token.Token }

func (f *False) String() string {
	return "false"
}
func (f *False) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.Boolean{}, nil
}

// Null represents the `null` reference
type Null struct{ token.Token }

func (n *Null) String() string {
	return "null"
}
func (n *Null) Semantic(p ast.SemanticParser) (ast.Type, io.Error) {
	// TODO: Scope and bytecode
	return &type_.Null{}, nil
}
