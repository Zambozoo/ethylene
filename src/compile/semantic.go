package compile

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/semantic"
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/compile/syntax"
	"geth-cody/io"
)

func filesChan(m map[string]ast.File) *data.Chan[ast.File] {
	c := data.NewChan[ast.File](io.Env.BufferSize)

	go func() {
		for _, file := range m {
			c.Send(file)
		}
		c.Close()
	}()

	return c
}

func semanticHelper(symbolMap syntax.SymbolMap, f func(ast.File, ast.SemanticParser, *data.AsyncSet[ast.Declaration]) io.Error) io.Error {
	parentsChan := filesChan(symbolMap.Files)
	visitedDecls := data.NewAsyncSet[ast.Declaration]()
	return parentsChan.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount,
		func(file ast.File) io.Error {
			p := semantic.NewParser(file, symbolMap)
			return f(file, p, visitedDecls)
		},
	)
}
func Semantic(symbolMap syntax.SymbolMap) (*bytecode.Bytecodes, io.Error) {
	if err := semanticHelper(symbolMap, ast.File.LinkParents); err != nil {
		return nil, err
	}

	if err := semanticHelper(symbolMap, ast.File.LinkFields); err != nil {
		return nil, err
	}

	var bytecodes *bytecode.Bytecodes
	bytecodesChan, closeBytecodes := data.RunUntilClosed(io.Env.BufferSize,
		func(b *bytecode.Bytecodes) {
			bytecodes.Add(b)
		},
	)

	semanticChan := filesChan(symbolMap.Files)
	semanticChan.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount,
		func(file ast.File) io.Error {
			p := semantic.NewParser(file, symbolMap)
			fileBytecodes, err := p.Parse()
			if err != nil {
				return err
			}
			bytecodesChan <- fileBytecodes
			return err
		},
	)

	closeBytecodes()

	return nil, nil
}
