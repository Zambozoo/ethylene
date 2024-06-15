package main

import (
	"geth-cody/compile"
	"geth-cody/io"
	"os"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

func main() {
	closeLogger := io.InitLogger()
	defer closeLogger()

	args, err := io.NewArgs(os.Args)
	if err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("running Ethylene bytecode transpiler", zap.Any("args", args))

	io.Infof("[Syntax] Start")
	symbolMap, err := compile.Syntax(args.InputFilePath)
	if err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("[Syntax] End",
		zap.Strings("file entries", maps.Keys(symbolMap.FileEntries)),
		zap.Strings("project files", maps.Keys(symbolMap.ProjectFiles)),
	)

	io.Infof("[Semantic] Start")
	bytecodes, err := compile.Semantic(symbolMap)
	if err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("[Semantic] End", zap.Int64("bytecodes length", bytecodes.Size()))

	io.Infof("[Export] Start")
	if err := compile.Export(args.OutputFilePath, bytecodes); err != nil {
		err.Log(io.Errorf)
		return
	}
	io.Infof("[Export] End", zap.String("out path", args.OutputFilePath))
}
