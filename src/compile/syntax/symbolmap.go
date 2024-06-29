package syntax

import (
	"geth-cody/ast"
	"geth-cody/io/path"
)

type SymbolMap struct {
	Projects map[string]*path.Project
	Files    map[string]ast.File
	Types    ast.Types
}
