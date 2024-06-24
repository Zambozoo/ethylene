package compile

import (
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/io"
	"geth-cody/io/path"
)

func Export(pathProvider path.Provider, outputFilePath string, bytecodes *bytecode.Bytecodes) io.Error {
	var bytes []byte
	pathProvider.WriteOutFile(outputFilePath, bytes)
	return nil
}
