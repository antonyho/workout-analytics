package logger_test

import (
	"bytes"
	"testing"

	"github.com/antonyho/workout-analytics/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.SetLevel(logger.InfoLevel)

	logger.Info("test info level message")
	logger.Error("test error level message")

	logMsg, err := buf.ReadString('\n')
	assert.NoError(t, err)
	assert.Contains(t, logMsg, "INFO:")
	assert.Contains(t, logMsg, "test info level message")

	logMsg, err = buf.ReadString('\n')
	assert.NoError(t, err)
	assert.Contains(t, logMsg, "ERROR:")
	assert.Contains(t, logMsg, "test error level message")
}

func TestLoggerSetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.SetLevel(logger.InfoLevel)

	logger.Info("test info level message")
	logger.Error("test error level message")

	// Contains all level log messages.
	logMsg, err := buf.ReadString('\n')
	assert.NoError(t, err)
	assert.Contains(t, logMsg, "INFO:")
	assert.Contains(t, logMsg, "test info level message")

	logMsg, err = buf.ReadString('\n')
	assert.NoError(t, err)
	assert.Contains(t, logMsg, "ERROR:")
	assert.Contains(t, logMsg, "test error level message")

	// Skipped Debug and Info log messages. Only Error log message.
	buf.Reset()
	logger.SetLevel(logger.ErrorLevel)

	logger.Info("test info level message")
	logger.Error("test error level message")

	logMsg, err = buf.ReadString('\n')
	assert.NoError(t, err)
	assert.Contains(t, logMsg, "ERROR:")
	assert.Contains(t, logMsg, "test error level message")
}
