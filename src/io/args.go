package io

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultOutFilePath = "eth.out.zip"
	projectFileName    = "eth.json"
)

type Args struct {
	InputFilePath   *FilePath
	MainProjectPath *FilePath
	OutputFilePath  string
}

func NewArgs(args []string) (*Args, Error) {
	if len(args) == 2 {
		args = append(args, defaultOutFilePath)
	} else if len(args) != 3 {
		return nil, &ZapError{
			Message: "usage: ./geth [MAIN_FILE_PATH] [OUT_FILE_PATH?]",
			Fields: []zapcore.Field{
				zap.Strings("args", args),
			},
		}
	}

	var result Args
	if outputFilePath, err := NewFilePath(args[2]); err != nil {
		return nil, err
	} else {
		result.OutputFilePath = *(*string)(outputFilePath)
	}

	var err Error
	if result.InputFilePath, err = NewFilePath(args[1]); err != nil {
		return nil, err
	}

	if mainProjectPath, err := MainProjectPath(result.InputFilePath); err != nil {
		return nil, err
	} else {
		// Main project path is input file path ancestor directory and guaranteed to be a file path.
		result.MainProjectPath, _ = mainProjectPath.(*FilePath)
	}

	return &result, nil
}

func MainProjectPath(fp Path) (Path, Error) {
	for path := fp.Dir(); path != nil; path = path.Dir() {
		if path.Join(projectFileName).Stat() == nil {
			return path, nil
		}
	}

	return nil, &ZapError{
		Message: "main project path not found",
		Fields: []zapcore.Field{
			zap.Any("path", fp),
		},
	}
}
