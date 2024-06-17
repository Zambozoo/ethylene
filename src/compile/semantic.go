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

func Semantic(symbolMap syntax.SymbolMap) (*bytecode.Bytecodes, io.Error) {
	parentsChan := filesChan(symbolMap.Files)
	visitedDecls := data.NewAsyncSet[ast.Declaration]()
	parentsChan.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount,
		func(file ast.File) io.Error {
			p := semantic.NewParser(file, symbolMap)
			return file.LinkParents(p, visitedDecls)
		},
	)

	methodsChan := filesChan(symbolMap.Files)
	methodsChan.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount,
		func(file ast.File) io.Error {
			p := semantic.NewParser(file, symbolMap)
			return file.LinkMethods(p, visitedDecls)
		},
	)

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
