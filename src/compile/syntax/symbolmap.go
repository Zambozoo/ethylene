package syntax

import (
	"geth-cody/ast"
	"geth-cody/io"
)

type SymbolMap struct {
	Projects map[string]*io.Project
	Files    map[string]ast.File
}
