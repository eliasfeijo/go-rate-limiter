package main

import (
	"fmt"
	"os"

	"github.com/eliasfeijo/go-rate-limiter/log"
)

// Logger is a simple logger that prints to stdout
type Logger struct {
}

// Log prints a line to stdout
func (l *Logger) Log(level log.Level, args ...interface{}) {
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
func (l *Logger) Logf(level log.Level, format string, args ...interface{}) {
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
