package log

// Level is the log level
type Level int

// You can use the following levels in your implementation:
// 0: Debug
// 1: Info
// 2: Warn
// 3: Error
// 4: Fatal
// 5: Panic
const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
	Panic
)

// Logger is the interface that wraps the Log and Logf methods.
type Logger interface {
	// Log logs a message at the given level. Arguments are handled in the manner of fmt.Print.
	Log(level Level, args ...interface{})
	// Logf logs a message at the given level. Arguments are handled in the manner of fmt.Printf.
	Logf(level Level, format string, args ...interface{})
}

var logger Logger

// SetLogger sets the logger to be used by the package
func SetLogger(l Logger) {
	logger = l
}

// Log is called internally by the library to log messages
func Log(level Level, args ...interface{}) {
	if logger != nil {
		logger.Log(level, args...)
	}
}

// Logf is called internally by the library to log messages
func Logf(level Level, format string, args ...interface{}) {
	if logger != nil {
		logger.Logf(level, format, args...)
	}
}
