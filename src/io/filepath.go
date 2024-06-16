package io

import (
	"os"
	filepath "path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FilePath string

func NewFilePath(path string) (*FilePath, Error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, &ZapError{
			Message: "couldn't get absolute path",
			Fields: []zapcore.Field{
				zap.String("path", path),
				zap.Error(err),
			},
		}
	}

	f := FilePath(absPath)
	return &f, nil
}

func (f *FilePath) Name() string {
	path := strings.Split(string(*f), string(os.PathSeparator))
	return strings.TrimSuffix(path[len(path)-1], fileExtension)
}

func (f *FilePath) Read() (string, Error) {
	fileBytes, err := os.ReadFile(string(*f))
	if err != nil {
		return "", &ZapError{
			Message: "couldn't read file",
			Fields: []zapcore.Field{
				zap.Any("file path", f),
				zap.Error(err),
			},
		}
	}

	return string(fileBytes), nil
}

func (f *FilePath) Stat() Error {
	_, err := os.Stat(string(*f))
	if err != nil {
		return &ZapError{
			Message: "couldn't stat file",
			Fields: []zapcore.Field{
				zap.Any("file path", f),
				zap.Error(err),
			},
		}
	}

	return nil
}

func (f *FilePath) Dir() Path {
	d := FilePath(filepath.Dir(string(*f)))
	if d == "." {
		return nil
	}

	return &d
}

func (f *FilePath) Join(path string) Path {
	j := FilePath(filepath.Join(string(*f), path))
	return &j
}

func (f *FilePath) String() string {
	return string(*f)
}
