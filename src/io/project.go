package io

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Project struct {
	path     Path
	Name     string
	Version  string
	Packages map[string]string
}

func NewProject(fp Path) (*Project, Error) {
	str, err := fp.Read()
	if err != nil {
		return nil, err
	}

	var p Project
	if err := yaml.Unmarshal([]byte(str), &p); err != nil {
		return nil, &ZapError{
			Message: "couldn't unmarshal project file",
			Fields: []zapcore.Field{
				zap.Any("path", fp),
				zap.Error(err),
			},
		}
	}

	p.path = fp
	return &p, nil
}

func (p *Project) Path() Path {
	return p.path
}

func (p *Project) String() string {
	return p.path.String()
}
