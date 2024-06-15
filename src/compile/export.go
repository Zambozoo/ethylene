package compile

import (
	"geth-cody/compile/semantic/bytecode"
	"geth-cody/io"
)

func Export(outputFilePath string, bytecodes *bytecode.Bytecodes) io.Error {
	var bytes []byte
	io.WriteOutFile(outputFilePath, bytes)
	return nil
}
