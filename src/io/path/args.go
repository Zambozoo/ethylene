package path

import (
	"fmt"
	"geth-cody/io"
	filepath "path/filepath"

	"go.uber.org/zap"
)

const (
	defaultOutFilePath = "eth.out.zip"
	projectFileName    = "eth.json"
)

// Args is a struct that holds all the arguments passed to the compiler.
type Args struct {
	InputFilePath   *File
	MainProjectPath *File
	OutputFilePath  string
}

func (a *Args) String() string {
	return fmt.Sprintf("Args{InputFilePath:%q,MainProjectPath:%q,OutputFilePath%q}",
		a.InputFilePath.String(),
		a.MainProjectPath.String(),
		a.OutputFilePath,
	)
}

// NewArgs instantiates a new Args struct from the provided arguments.
func NewArgs(args []string) (*Args, io.Error) {
	if len(args) == 2 {
		args = append(args, defaultOutFilePath)
	} else if len(args) != 3 {
		return nil, io.NewError("usage: ./geth [MAIN_FILE_PATH] [OUT_FILE_PATH?]", zap.Strings("args", args))
	}

	var result Args
	if outputFilePath, err := newFilePath(args[2]); err != nil {
		return nil, err
	} else {
		result.OutputFilePath = *(*string)(outputFilePath)
	}

	var err io.Error
	if result.InputFilePath, err = newFilePath(args[1]); err != nil {
		return nil, err
	}

	if mainProjectPath, err := result.InputFilePath.MainProjectPath(); err != nil {
		return nil, err
	} else {
		// Main project path is input file path ancestor directory and guaranteed to be a file path.
		result.MainProjectPath, _ = mainProjectPath.(*File)
	}

	return &result, nil
}

func newFilePath(path string) (*File, io.Error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, io.NewError("couldn't get absolute path",
			zap.String("path", path),
			zap.Error(err),
		)
	}

	f := File(absPath)
	return &f, nil
}
