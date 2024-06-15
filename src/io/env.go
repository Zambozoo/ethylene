package io

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type Environment struct {
	ThreadCount int    `env:"THREAD_COUNT,default=16"`
	BufferSize  int    `env:"BUFFER_SIZE,default=128"`
	StdLibPath  string `env:"ETHYLENE_HOME,default=std_0.0.0.zip"`
}

var Env = func() Environment {
	var env Environment
	if err := envconfig.Process(context.Background(), &env); err != nil {
		panic(fmt.Errorf("couldn't load environment: %w", err))
	}

	return env
}()
