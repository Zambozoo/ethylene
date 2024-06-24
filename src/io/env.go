package io

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

// Environment holds the environment variables for the compiler.
type Environment struct {
	ThreadCount int    `env:"THREAD_COUNT,default=16"`
	BufferSize  int    `env:"BUFFER_SIZE,default=128"`
	StdLibPath  string `env:"ETHYLENE_HOME,default=std_0.0.0.zip"`
}

// Env is the environment variables for the compiler.
var Env = func() Environment {
	var env Environment
	if err := envconfig.Process(context.Background(), &env); err != nil {
		panic(fmt.Errorf("couldn't load environment: %w", err))
	}

	return env
}()
