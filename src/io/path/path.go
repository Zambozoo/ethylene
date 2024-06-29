package path

import (
	"fmt"
	"geth-cody/io"

	"go.uber.org/zap"
)

const fileExtension = ".eth"

// Path is a path to a zipped or unzipped file or directory.
type Path interface {
	fmt.Stringer
	// Name returns the name of the file or directory without any leading path.
	Name() string
	// Read returns the contents of the file as a string.
	Read() (string, io.Error)
	// Stat returns an error if the file or directory does not exist.
	Stat() io.Error
	// Dir returns the parent directory of the file or directory.
	Dir() Path
	// Join returns a new path with the provided path appended to the end of the current path.
	Join(string) Path
	// MainProjectPath returns the path to the main project file.
	MainProjectPath() (Path, io.Error)
}

func mainProjectPath(fp Path) (Path, io.Error) {
	for path := fp.Dir(); path != nil; path = path.Dir() {
		if path.Join(projectFileName).Stat() == nil {
			return path, nil
		}
	}

	return nil, io.NewError("main project path not found", zap.Stringer("path", fp))
}
