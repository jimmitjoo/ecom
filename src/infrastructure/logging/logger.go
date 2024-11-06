package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
)

type contextKey string

const loggerKey = contextKey("logger")

// Logger wraps zap logger with additional context
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new structured logger
func NewLogger() (*Logger, error) {
	env := os.Getenv("GO_ENV")

	if env == "development" {
		return NewDevelopmentLogger()
	}

	return NewProductionLogger()
}

// NewProductionLogger creates a new production logger
func NewProductionLogger() (*Logger, error) {
	config := zap.NewProductionConfig()
	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: zapLogger}, nil
}

// WithContext adds logger to context
func WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	return &Logger{Logger: zap.NewNop()}
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithRequestID adds request ID to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

// WithTraceID adds trace ID to the logger
func (l *Logger) WithTraceID(traceID string) *Logger {
	return l.WithFields(zap.String("trace_id", traceID))
}

// WithUserID adds user ID to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithFields(zap.String("user_id", userID))
}

// WithError adds error to the logger
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}
