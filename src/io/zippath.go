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
