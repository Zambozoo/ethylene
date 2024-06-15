package io

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	filepath "path/filepath"
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

type ZipPath struct {
	zipPath string
	path    string
}

func (z *ZipPath) Name() string {
	path := strings.Split(z.path, string(os.PathSeparator))
	return strings.TrimSuffix(path[len(path)-1], fileExtension)
}

func (z *ZipPath) Read() (string, Error) {
	r, err := zip.OpenReader(z.zipPath)
	if err != nil {
		return "", &ZapError{
			Message: "couldn't open zip file",
			Fields: []zapcore.Field{
				zap.String("zip path", z.zipPath),
				zap.Error(err),
			},
		}
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != z.path {
			continue
		}

		rc, err := f.Open()
		buf := new(strings.Builder)
		if err != nil {
			return "", &ZapError{
				Message: "couldn't open file in zip",
				Fields: []zapcore.Field{
					zap.Any("file path", z),
					zap.Error(err),
				},
			}
		}

		defer rc.Close()
		if _, err = io.Copy(buf, rc); err != nil {
			return "", &ZapError{
				Message: "couldn't read file from zip",
				Fields: []zapcore.Field{
					zap.Any("file path", z),
					zap.Error(err),
				},
			}
		}

		return buf.String(), nil
	}

	return "", &ZapError{
		Message: "file not found in zip",
		Fields: []zapcore.Field{
			zap.Any("file path", z),
		},
	}
}

func (z *ZipPath) Stat() Error {
	r, err := zip.OpenReader(z.zipPath)
	if err != nil {
		return &ZapError{
			Message: "couldn't open zip file",
			Fields: []zapcore.Field{
				zap.String("zip path", z.zipPath),
				zap.Error(err),
			},
		}
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == z.path {
			return nil
		}
	}

	return &ZapError{
		Message: "file not found in zip",
		Fields: []zapcore.Field{
			zap.Any("file path", z),
		},
	}
}

func (z *ZipPath) Dir() Path {
	d := filepath.Dir(z.path)
	if d == "." {
		return nil
	}

	return &ZipPath{zipPath: z.zipPath, path: d}
}

func (z *ZipPath) Join(path string) Path {
	return &ZipPath{zipPath: z.zipPath, path: filepath.Join(z.path, path)}
}

func (z *ZipPath) String() string {
	return fmt.Sprintf("%s:%s", z.zipPath, z.path)
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
