package io

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const fileExtension = ".eth"

type Path interface {
	fmt.Stringer
	Name() string
	Read() (string, Error)
	Stat() Error
	Dir() Path
	Join(string) Path
}

func NewPath(path string) (Path, Error) {
	splitPath := strings.Split(path, ":")
	if path == "" || len(splitPath) > 2 {
		return nil, &ZapError{
			Message: "invalid file path",
			Fields: []zapcore.Field{
				zap.String("path", path),
			},
		}
	}

	absPath, err := filepath.Abs(splitPath[0])
	if err != nil {
		return nil, &ZapError{
			Message: "couldn't get absolute path",
			Fields: []zapcore.Field{
				zap.String("path", splitPath[0]),
				zap.Error(err),
			},
		}
	}

	if len(splitPath) == 2 {
		return &ZipPath{zipPath: absPath, path: splitPath[1]}, nil
	}

	f := FilePath(absPath)
	return &f, nil
}

func WriteOutFile(outputFilePath string, bytes []byte) Error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return &ZapError{
			Message: "couldn't create out file",
			Fields: []zapcore.Field{
				zap.String("out file path", outputFilePath),
				zap.Error(err),
			},
		}
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return &ZapError{
			Message: "couldn't write to out file",
			Fields: []zapcore.Field{
				zap.String("out file path", outputFilePath),
				zap.Error(err),
			},
		}
	}

	return nil
}
