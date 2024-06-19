package path

import (
	"geth-cody/io"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Project is a project file.
type Project struct {
	path     Path
	Name     string
	Version  string
	Packages map[string]string
}

// NewProject returns a new project whose contents are read from a given path.
func NewProject(fp Path) (*Project, io.Error) {
	str, err := fp.Read()
	if err != nil {
		return nil, err
	}

	var p Project
	if err := yaml.Unmarshal([]byte(str), &p); err != nil {
		return nil, io.NewError("couldn't unmarshal project file",
			zap.Any("path", fp),
			zap.Error(err),
		)
	}

	p.path = fp
	return &p, nil
}

// Path returns the path to the project file.
func (p *Project) Path() Path {
	return p.path
}

// String returns the string representation of the project file path.
func (p *Project) String() string {
	return p.path.String()
}
