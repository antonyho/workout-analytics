package logger

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	InfoLevel LogLevel = iota
	ErrorLevel
)

// String represents the log level in string.
func (level LogLevel) String() string {
	switch level {
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	}

	return ""
}

var (
	info *log.Logger
	err  *log.Logger

	logLevel = ErrorLevel
)

func init() {
	info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	err = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
}

// SetLevel sets the log level for this logger package.
func SetLevel(level LogLevel) {
	logLevel = level
}

// SetOutput sets the output destinations for all log messages.
func SetOutput(out io.Writer) {
	info.SetOutput(out)
	err.SetOutput(out)
}

// Info adds log message with "INFO:" prefix at the beginning.
func Info(msgFormat string, v ...any) {
	if InfoLevel >= logLevel {
		info.Printf(msgFormat, v...)
	}
}

// Error adds log message with "ERROR:" prefix at the beginning.
func Error(msgFormat string, v ...any) {
	if ErrorLevel >= logLevel {
		err.Printf(msgFormat, v...)
	}
}
