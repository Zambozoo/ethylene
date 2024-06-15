package io

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger() func() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.CallerKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableStacktrace = true

	logger = zap.Must(config.Build())
	if logger == nil {
		panic("")
	}

	return func() {
		logger.Sync()
	}
}

// Info logs at the info level with the common logger.
func Infof(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Debug logs at the debug level with the common logger.
func Debugf(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warn logs at the warn level with the common logger.
func Warnf(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs at the panic level with the common logger.
func Errorf(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
