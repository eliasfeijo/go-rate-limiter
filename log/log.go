package log

import (
	"fmt"
	"strings"
)

// LogLevel is the log LogLevel
type LogLevel uint8

// You can use the following LogLevels in your implementation:
// 0: Debug
// 1: Info
// 2: Warn
// 3: Error
// 4: Fatal
// 5: Panic
const (
	Debug LogLevel = iota
	Info
	Warn
	Error
	Fatal
	Panic
)

// Logger is the interface that wraps the Log and Logf methods.
type Logger interface {
	// Log logs a message at the given LogLevel. Arguments are handled in the manner of fmt.Print.
	Log(logLevel LogLevel, args ...interface{})
	// Logf logs a message at the given LogLevel. Arguments are handled in the manner of fmt.Printf.
	Logf(logLevel LogLevel, format string, args ...interface{})
}

var logger Logger

// SetLogger sets the logger to be used by the package
func SetLogger(l Logger) {
	logger = l
}

// Log is called internally by the library to log messages
func Log(logLevel LogLevel, args ...interface{}) {
	if logger != nil {
		logger.Log(logLevel, args...)
	}
}

// Logf is called internally by the library to log messages
func Logf(logLevel LogLevel, format string, args ...interface{}) {
	if logger != nil {
		logger.Logf(logLevel, format, args...)
	}
}

// ParseLogLevel parses a string into a LogLevel
func ParseLogLevel(logLevel string) (LogLevel, error) {
	switch strings.TrimSpace(strings.ToLower(logLevel)) {
	case "debug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error":
		return Error, nil
	case "fatal":
		return Fatal, nil
	case "panic":
		return Panic, nil
	default:
		return Debug, fmt.Errorf("Invalid log level: %s", logLevel)
	}
}
