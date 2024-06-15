package io

import (
	"errors"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Error interface {
	error
	Log(f func(string, ...zapcore.Field))
}

type ZapErrors []Error

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

	e := make(ZapErrors, 0, n)
	for _, err := range errs {
		if err != nil {
			e = append(e, err)
		}
	}

	return e
}

func (e ZapErrors) Error() string {
	err := errors.Join(e)
	return err.Error()
}

func (es ZapErrors) Log(f func(string, ...zapcore.Field)) {
	for _, e := range es {
		e.Log(f)
	}
}

type ZapError struct {
	Message string
	Fields  []zapcore.Field
}

func (e *ZapError) Log(f func(string, ...zapcore.Field)) {
	f(e.Message, e.Fields...)
}

func (e *ZapError) Error() string {
	return fmt.Sprintf(e.Message+"%v", e.Fields)
}

func NewError(message string, fields ...zapcore.Field) Error {
	return &ZapError{
		Message: message,
		Fields:  fields,
	}
}
