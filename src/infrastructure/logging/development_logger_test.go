package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestDevelopmentLogger(t *testing.T) {
	// Skapa en observerbar development logger
	core, recorded := observer.New(zapcore.DebugLevel)
	logger := &Logger{
		Logger: zap.New(core,
			zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel)),
	}

	// Testa olika loggnivåer
	testCases := []struct {
		level   string
		logFunc func(string, ...zap.Field)
	}{
		{"debug", logger.Debug},
		{"info", logger.Info},
		{"warn", logger.Warn},
		{"error", logger.Error},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			recorded.TakeAll() // Rensa tidigare logs
			message := tc.level + " test message"

			// Logga med extra utvecklingsinformation
			tc.logFunc(message,
				zap.String("request_id", "test-123"),
				zap.String("function", "TestFunction"),
				zap.Int("line", 42),
			)

			logs := recorded.All()
			assert.Len(t, logs, 1)
			assert.Equal(t, message, logs[0].Message)
			assert.Contains(t, logs[0].ContextMap(), "request_id")
			assert.Contains(t, logs[0].ContextMap(), "function")
			assert.Contains(t, logs[0].ContextMap(), "line")
			assert.Contains(t, logs[0].Caller.String(), "development_logger_test.go")
		})
	}
}
