package compile

import (
	"geth-cody/ast"
	"geth-cody/compile/data"
	"geth-cody/compile/semantic"
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/compile/syntax"
	"geth-cody/io"
)

func fileEntriesChan(m map[string]syntax.FileEntry) *data.Chan[syntax.FileEntry] {
	c := data.NewChan[syntax.FileEntry](io.Env.BufferSize)

	go func() {
		for _, fileEntry := range m {
			c.Send(fileEntry)
		}
		c.Close()
	}()

	return c
}

func Semantic(symbolMap SymbolMap) (*bytecode.Bytecodes, io.Error) {
	parentsChan := fileEntriesChan(symbolMap.FileEntries)
	visitedDecls := data.NewAsyncSet[ast.Declaration]()
	parentsChan.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount,
		func(fe syntax.FileEntry) io.Error {
			p := semantic.NewParser(fe, symbolMap.ProjectFiles, symbolMap.FileEntries)
			return fe.File.LinkParents(p, visitedDecls)
		},
	)

	var bytecodes *bytecode.Bytecodes
	bytecodesChan, closeBytecodes := data.RunUntilClosed(io.Env.BufferSize,
		func(b *bytecode.Bytecodes) {
			bytecodes.Add(b)
		},
	)

	semanticChan := fileEntriesChan(symbolMap.FileEntries)
	semanticChan.AsyncForEach(io.Env.BufferSize, io.Env.ThreadCount,
		func(fe syntax.FileEntry) io.Error {
			p := semantic.NewParser(fe, symbolMap.ProjectFiles, symbolMap.FileEntries)
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
