package path

import (
	"geth-cody/io"
	"os"
	filepath "path/filepath"
	"strings"

	"go.uber.org/zap"
)

// File is a path to an unzipped .eth file.
type File string

func (f *File) Name() string {
	path := strings.Split(string(*f), string(os.PathSeparator))
	return strings.TrimSuffix(path[len(path)-1], fileExtension)
}

// Read returns the contents of the file as a string.
func (f *File) Read() (string, io.Error) {
	fileBytes, err := os.ReadFile(string(*f))
	if err != nil {
		return "", io.NewError("couldn't read file",
			zap.Any("file path", f),
			zap.Error(err),
		)
	}

	return string(fileBytes), nil
}

// Stat returns an error if the file does not exist.
func (f *File) Stat() io.Error {
	_, err := os.Stat(string(*f))
	if err != nil {
		return io.NewError("couldn't stat file",
			zap.Any("file path", f),
			zap.Error(err),
		)
	}

	return nil
}

// Dir returns the parent directory of the file.
func (f *File) Dir() Path {
	d := File(filepath.Dir(string(*f)))
	if d == "." {
		return nil
	}

	return &d
}

// Join returns a new path with the provided path appended to the end of the current path.
func (f *File) Join(path string) Path {
	j := File(filepath.Join(string(*f), path))
	return &j
}

// MainProjectPath returns the path to the main project file.
func (f *File) MainProjectPath() (Path, io.Error) {
	return mainProjectPath(f)
}

// String returns the string representation of the file path.
func (f *File) String() string {
	return string(*f)
}
