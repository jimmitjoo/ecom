package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewDevelopmentLogger() (*Logger, error) {
	config := zap.NewDevelopmentConfig()

	// Set development-specific settings
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	// Enable all fields for development mode
	config.Development = true
	config.DisableCaller = false
	config.DisableStacktrace = false

	// Set minimum log level to Debug in development mode
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	zapLogger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}
