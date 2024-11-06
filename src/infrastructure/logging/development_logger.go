package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewDevelopmentLogger() (*Logger, error) {
	config := zap.NewDevelopmentConfig()

	// Sätt development-specifika inställningar
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	// Aktivera alla fält för utvecklingsläge
	config.Development = true
	config.DisableCaller = false
	config.DisableStacktrace = false

	// Sätt lägsta lognivå till Debug i utvecklingsläge
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
