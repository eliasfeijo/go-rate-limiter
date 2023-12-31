package main

import (
	"fmt"
	"os"

	"github.com/eliasfeijo/go-rate-limiter/log"
)

// Logger is a simple logger that prints to stdout
type Logger struct {
	LogLevel log.LogLevel
}

func NewLogger(logLevel log.LogLevel) *Logger {
	if logLevel < log.Debug {
		logLevel = log.Debug
	} else if logLevel > log.Panic {
		logLevel = log.Panic
	}
	return &Logger{logLevel}
}

// Log prints a line to stdout
func (l *Logger) Log(level log.LogLevel, args ...interface{}) {
	if level < l.LogLevel {
		return
	}
	switch level {
	case log.Debug:
		fmt.Println("DEBUG: ", args)
	case log.Info:
		fmt.Println("INFO: ", args)
	case log.Warn:
		fmt.Println("WARN: ", args)
	case log.Error:
		fmt.Println("ERROR: ", args)
	case log.Fatal:
		fmt.Println("FATAL: ", args)
		os.Exit(1)
	case log.Panic:
		panic(fmt.Sprint(args...))
	default:
		fmt.Println(args...)
	}
}

// Logf prints a formatted line to stdout
func (l *Logger) Logf(level log.LogLevel, format string, args ...interface{}) {
	if level < l.LogLevel {
		return
	}
	switch level {
	case 0:
		fmt.Printf("DEBUG: "+format, args...)
	case log.Info:
		fmt.Printf("INFO: "+format, args...)
	case log.Warn:
		fmt.Printf("WARN: "+format, args...)
	case log.Error:
		fmt.Printf("ERROR: "+format, args...)
	case log.Fatal:
		fmt.Printf("FATAL: "+format, args...)
		os.Exit(1)
	case log.Panic:
		panic(fmt.Sprintf(format, args...))
	default:
		fmt.Printf(format, args...)
	}
	fmt.Println()
}
