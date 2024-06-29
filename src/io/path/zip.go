package path

import (
	"archive/zip"
	"fmt"
	"geth-cody/io"
	goio "io"
	"os"
	filepath "path/filepath"
	"strings"

	"go.uber.org/zap"
)

// Zip is a path to a zipped file.
type Zip struct {
	zipPath string
	path    string
}

// Name returns the name of the file or directory without any leading path.
func (z *Zip) Name() string {
	path := strings.Split(z.path, string(os.PathSeparator))
	return strings.TrimSuffix(path[len(path)-1], fileExtension)
}

// Read returns the contents of the file as a string.
func (z *Zip) Read() (string, io.Error) {
	r, err := zip.OpenReader(z.zipPath)
	if err != nil {
		return "", io.NewError("couldn't open zip file",
			zap.String("zip path", z.zipPath),
			zap.Error(err),
		)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != z.path {
			continue
		}

		rc, err := f.Open()
		buf := new(strings.Builder)
		if err != nil {
			return "", io.NewError("couldn't open file in zip",
				zap.Stringer("file path", z),
				zap.Error(err),
			)
		}

		defer rc.Close()
		if _, err = goio.Copy(buf, rc); err != nil {
			return "", io.NewError("couldn't read file from zip",
				zap.Stringer("file path", z),
				zap.Error(err),
			)
		}

		return buf.String(), nil
	}

	return "", io.NewError("file not found in zip", zap.Stringer("file path", z))
}

// Stat returns an error if the file or directory does not exist.
func (z *Zip) Stat() io.Error {
	r, err := zip.OpenReader(z.zipPath)
	if err != nil {
		return io.NewError("couldn't open zip file",
			zap.String("zip path", z.zipPath),
			zap.Error(err),
		)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == z.path {
			return nil
		}
	}

	return io.NewError("file not found in zip", zap.Stringer("file path", z))
}

// Dir returns the parent directory of the file or directory.
func (z *Zip) Dir() Path {
	d := filepath.Dir(z.path)
	if d == "." {
		return nil
	}

	return &Zip{zipPath: z.zipPath, path: d}
}

// Join returns a new path with the provided path appended to the end of the current path.
func (z *Zip) Join(path string) Path {
	return &Zip{zipPath: z.zipPath, path: filepath.Join(z.path, path)}
}

// MainProjectPath returns the path to the main project file.
func (z *Zip) MainProjectPath() (Path, io.Error) {
	return mainProjectPath(z)
}

// String returns a string representation of the path.
func (z *Zip) String() string {
	return fmt.Sprintf("%s:%s", z.zipPath, z.path)
}
