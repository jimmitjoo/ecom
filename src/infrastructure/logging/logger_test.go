package logging

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger(t *testing.T) {
	// Skapa en ny logger
	logger, err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.Logger)
}

func TestLoggerWithContext(t *testing.T) {
	// Skapa en observerbar logger för att kunna inspektera loggade meddelanden
	core, recorded := observer.New(zap.InfoLevel)
	testLogger := &Logger{
		Logger: zap.New(core),
	}

	// Skapa context med logger
	ctx := WithContext(context.Background(), testLogger)
	assert.NotNil(t, ctx)

	// Hämta logger från context
	loggerFromCtx := FromContext(ctx)
	assert.NotNil(t, loggerFromCtx)
	assert.Equal(t, testLogger, loggerFromCtx)

	// Testa att logga något
	loggerFromCtx.Info("test message",
		zap.String("key", "value"),
	)

	// Verifiera att meddelandet loggades
	logs := recorded.All()
	assert.Len(t, logs, 1)
	assert.Equal(t, "test message", logs[0].Message)
	assert.Equal(t, "value", logs[0].ContextMap()["key"])
}

func TestFromContextWithNoLogger(t *testing.T) {
	// Testa att FromContext returnerar en nop-logger när ingen logger finns i context
	ctx := context.Background()
	logger := FromContext(ctx)
	assert.NotNil(t, logger)

	// Verifiera att det är en nop-logger genom att logga något
	// Detta ska inte orsaka några fel
	logger.Info("this should not cause any errors")
}

func TestLogLevels(t *testing.T) {
	// Skapa en observerbar logger
	core, recorded := observer.New(zap.DebugLevel)
	testLogger := &Logger{
		Logger: zap.New(core),
	}

	// Testa olika log levels
	testCases := []struct {
		level   string
		logFunc func(string, ...zap.Field)
	}{
		{"debug", testLogger.Debug},
		{"info", testLogger.Info},
		{"warn", testLogger.Warn},
		{"error", testLogger.Error},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			recorded.TakeAll() // Rensa tidigare logs
			message := tc.level + " message"

			tc.logFunc(message, zap.String("level", tc.level))

			logs := recorded.All()
			assert.Len(t, logs, 1)
			assert.Equal(t, message, logs[0].Message)
			assert.Equal(t, tc.level, logs[0].ContextMap()["level"])
		})
	}
}

func TestLoggerFields(t *testing.T) {
	// Skapa en observerbar logger
	core, recorded := observer.New(zap.InfoLevel)
	testLogger := &Logger{
		Logger: zap.New(core),
	}

	// Testa olika typer av fields
	testLogger.Info("test message",
		zap.String("string", "value"),
		zap.Int("int", 123),
		zap.Bool("bool", true),
		zap.Float64("float", 123.456),
	)

	logs := recorded.All()
	assert.Len(t, logs, 1)
	fields := logs[0].ContextMap()

	assert.Equal(t, "value", fields["string"])
	assert.Equal(t, int64(123), fields["int"])
	assert.Equal(t, true, fields["bool"])
	assert.Equal(t, 123.456, fields["float"])
}

func TestLoggerConcurrency(t *testing.T) {
	// Skapa en observerbar logger
	core, recorded := observer.New(zap.InfoLevel)
	testLogger := &Logger{
		Logger: zap.New(core),
	}

	// Logga från flera goroutines samtidigt
	const numGoroutines = 100
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			testLogger.Info("concurrent message",
				zap.Int("goroutine_id", id),
			)
			done <- true
		}(i)
	}

	// Vänta på att alla goroutines är klara
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verifiera att alla meddelanden loggades
	logs := recorded.All()
	assert.Len(t, logs, numGoroutines)

	// Verifiera att alla goroutine IDs finns med
	ids := make(map[int]bool)
	for _, log := range logs {
		id := int(log.ContextMap()["goroutine_id"].(int64))
		ids[id] = true
	}
	assert.Len(t, ids, numGoroutines)
}
