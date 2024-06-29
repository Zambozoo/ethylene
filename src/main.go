package main

import (
	"geth-cody/compile"
	"geth-cody/io"
	"geth-cody/io/path"
	"os"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

func main() {
	closeLogger := io.InitLogger()
	defer closeLogger()

	args, err := path.NewArgs(os.Args)
	if err != nil {
		err.Log(io.Errorf)
		return
	}
	pathProvider := path.DefaultProvider{}
	io.Infof("running Ethylene bytecode transpiler", zap.Stringer("args", args))

	io.Infof("[Syntax] Start")
	symbolMap, err := compile.Syntax(&pathProvider, args.InputFilePath)
	if err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("[Syntax] End",
		zap.Strings("file entries", maps.Keys(symbolMap.Files)),
		zap.Strings("project files", maps.Keys(symbolMap.Projects)),
	)

	io.Infof("[Semantic] Start")
	bytecodes, err := compile.Semantic(symbolMap)
	if err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("[Semantic] End", zap.Int64("bytecodes length", bytecodes.Size()))

	io.Infof("[Export] Start")
	if err := compile.Export(&pathProvider, args.OutputFilePath, bytecodes); err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("[Export] End", zap.String("out path", args.OutputFilePath))
}
