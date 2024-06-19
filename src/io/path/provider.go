package path

import (
	"geth-cody/io"
	"os"
	filepath "path/filepath"
	"strings"

	"go.uber.org/zap"
)

// Provider creates new paths and writes out files.
type Provider interface {
	NewPath(path string) (Path, io.Error)
	WriteOutFile(outputFilePath string, bytes []byte) io.Error
}

// DefaultProvider provides a default implementation for creating new paths and writing out files.
type DefaultProvider struct{}

// NewPath creates a new path or returns an error if the path is invalid.
func (*DefaultProvider) NewPath(path string) (Path, io.Error) {
	splitPath := strings.Split(path, ":")
	if path == "" || len(splitPath) > 2 {
		return nil, io.NewError("invalid file path", zap.String("path", path))
	}

	absPath, err := filepath.Abs(splitPath[0])
	if err != nil {
		return nil, io.NewError("couldn't get absolute path",
			zap.String("path", splitPath[0]),
			zap.Error(err),
		)
	}

	if len(splitPath) == 2 {
		return &Zip{zipPath: absPath, path: splitPath[1]}, nil
	}

	f := File(absPath)
	return &f, nil
}

// WriteOutFile writes out the provided bytes to the provided output file path overwriting any existing file.
func (*DefaultProvider) WriteOutFile(outputFilePath string, bytes []byte) io.Error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return io.NewError("couldn't create out file",
			zap.String("out file path", outputFilePath),
			zap.Error(err),
		)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return io.NewError("couldn't write to out file",
			zap.String("out file path", outputFilePath),
			zap.Error(err),
		)
	}

	return nil
}
