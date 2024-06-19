package compile

import (
	"geth-cody/io"
	"geth-cody/io/path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPath struct {
	stringFunc          func() string
	nameFunc            func() string
	readFunc            func() (string, io.Error)
	statFunc            func() io.Error
	dirFunc             func() path.Path
	joinFunc            func(string) path.Path
	mainProjectPathFunc func() (path.Path, io.Error)
}

func (m *mockPath) String() string                         { return m.stringFunc() }
func (m *mockPath) Name() string                           { return m.nameFunc() }
func (m *mockPath) Read() (string, io.Error)               { return m.readFunc() }
func (m *mockPath) Stat() io.Error                         { return m.statFunc() }
func (m *mockPath) Dir() path.Path                         { return m.dirFunc() }
func (m *mockPath) Join(s string) path.Path                { return m.joinFunc(s) }
func (m *mockPath) MainProjectPath() (path.Path, io.Error) { return m.mainProjectPathFunc() }

func mockDirPath(paths map[string]path.Path) path.Path {
	return &mockPath{
		stringFunc: func() string {
			return ""
		},
		joinFunc: func(s string) path.Path {
			if p, ok := paths[s]; ok {
				return p
			}
			return &mockPath{}
		},
	}
}

func mockProjectPath(paths map[string]path.Path) path.Path {
	return &mockPath{
		stringFunc: func() string {
			return "eth.yaml"
		},
		readFunc: func() (string, io.Error) {
			return `name: "example"
version: "0.0.0"
packages:
  example_package: "0.0.0"`, nil
		},
		dirFunc: func() path.Path {
			return mockDirPath(paths)
		},
	}
}

func mockPaths(codeContents map[string]string) map[string]path.Path {
	m := map[string]path.Path{}
	m["eth.yaml"] = mockProjectPath(m)
	m[""] = mockDirPath(m)
	for filePath, content := range codeContents {
		filePath, content := filePath, content
		m[filePath] = &mockPath{
			stringFunc:          func() string { return filePath },
			nameFunc:            func() string { return strings.TrimSuffix(filePath, ".eth") },
			readFunc:            func() (string, io.Error) { return content, nil },
			statFunc:            func() io.Error { return nil },
			dirFunc:             func() path.Path { return m[""] },
			mainProjectPathFunc: func() (path.Path, io.Error) { return m["eth.yaml"], nil },
		}
	}

	return m
}

type mockPathProvider struct {
	newPathFunc      func(path string) (path.Path, io.Error)
	writeOutFileFunc func(outputFilePath string, bytes []byte) io.Error
}

func (m *mockPathProvider) NewPath(path string) (path.Path, io.Error) {
	return m.newPathFunc(path)
}
func (m *mockPathProvider) WriteOutFile(outputFilePath string, bytes []byte) io.Error {
	return m.writeOutFileFunc(outputFilePath, bytes)
}

func newMockPathProvider(paths map[string]path.Path) path.Provider {
	return &mockPathProvider{
		newPathFunc: func(s string) (path.Path, io.Error) {
			if p, ok := paths[s]; ok {
				return p, nil
			}
			return nil, io.NewError("test error")
		},
	}
}

func TestSyntax(t *testing.T) {
	t.Parallel()

	type test struct {
		name         string
		path         path.Path
		pathProvider path.Provider
		errFunc      assert.ErrorAssertionFunc
	}
	tests := []test{
		{
			name:    "invalid file",
			path:    mockPaths(map[string]string{"test.eth": ""})["test.eth"],
			errFunc: assert.Error,
		},
		{
			name:    "valid file",
			path:    mockPaths(map[string]string{"test.eth": "class Class {}"})["test.eth"],
			errFunc: assert.NoError,
		},
		{
			name: "valid files",
			path: mockPaths(map[string]string{
				"test.eth":  `import "test2.eth"; class Class <: Interface {}`,
				"test2.eth": "interface Interface {}",
			})["test.eth"],
			errFunc: assert.NoError,
		},
		{
			name: "invalid second file",
			path: mockPaths(map[string]string{
				"test.eth":  `import "test2.eth"; class Class <: Interface {}`,
				"test2.eth": "",
			})["test.eth"],
			errFunc: assert.Error,
		},
		func() test {
			paths := mockPaths(map[string]string{
				"test.eth":                                 `import example_package("test2.eth"); class Class <: Interface {}`,
				"pkgs/example_package~0.0.0.zip":           "",
				"pkgs/example_package~0.0.0.zip:test2.eth": "interface Interface {}",
			})
			return test{
				name:         "valid package file",
				path:         paths["test.eth"],
				pathProvider: newMockPathProvider(paths),
				errFunc:      assert.NoError,
			}
		}(),
		func() test {
			paths := mockPaths(map[string]string{
				"test.eth":                                 `import example_package("test2.eth"); class Class <: Interface {}`,
				"pkgs/example_package~0.0.0.zip":           "",
				"pkgs/example_package~0.0.0.zip:test2.eth": "",
			})
			return test{
				name:         "invalid package file",
				path:         paths["test.eth"],
				pathProvider: newMockPathProvider(paths),
				errFunc:      assert.Error,
			}
		}(),
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := Syntax(test.pathProvider, test.path)
			test.errFunc(t, err)
		})
	}
}
