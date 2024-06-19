package io

import (
	"errors"
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Error is an interface for errors that can be logged.
type Error interface {
	error
	// Log logs the error.
	Log(f func(string, ...zapcore.Field))
}

// JoinError joins multiple errors into a single error.
func JoinError(errs ...Error) Error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}

	e := make(zapErrors, 0, n)
	for _, err := range errs {
		if err != nil {
			e = append(e, err)
		}
	}

	return e
}

type zapErrors []Error

func (e zapErrors) Error() string {
	err := errors.Join(e)
	return err.Error()
}

func (es zapErrors) Log(f func(string, ...zapcore.Field)) {
	for _, e := range es {
		e.Log(f)
	}
}

type zapError struct {
	message string
	fields  []zapcore.Field
}

func (e *zapError) Log(f func(string, ...zapcore.Field)) {
	f(e.message, e.fields...)
}

func (e *zapError) Error() string {
	return fmt.Sprintf(e.message+"%v", e.fields)
}

func NewError(message string, fields ...zapcore.Field) Error {
	return &zapError{
		message: message,
		fields:  fields,
	}
}
