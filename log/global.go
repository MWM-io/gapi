package log

import (
	"fmt"
	"os"
	"sync"
)

var (
	globalLogger   *Logger
	globalLoggerMu sync.RWMutex
)

// GlobalLogger return the global logger.
// If the global logger hasn't been initialized, it returns a default logger printing json output to stdout.
func GlobalLogger() *Logger {
	globalLoggerMu.RLock()
	if globalLogger != nil {
		defer globalLoggerMu.RUnlock()

		return globalLogger
	}

	globalLoggerMu.RUnlock()

	SetGlobalLogger(NewDefaultLogger(NewWriter(EntryMarshalerFunc(func(entry Entry) []byte {
		var lastFrameStr string
		lastFrame, ok := entry.StackTrace.Last()
		if ok && entry.Severity <= WarnSeverity {
			lastFrameStr = fmt.Sprintf("\n                              %s %s %d", lastFrame.File, lastFrame.Function, lastFrame.Line)
		}

		return []byte(fmt.Sprintf(
			"%-9s %s | %s%s",
			entry.Severity.String(),
			entry.Timestamp.Format("15:04:05.999999999"),
			entry.Message,
			lastFrameStr,
		))
	}), os.Stdout)))

	return GlobalLogger()
}

// SetGlobalLogger sets the global logger.
func SetGlobalLogger(l *Logger) {
	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()

	globalLogger = l
}

// Log log a message using the global logger.
func Log(msg string, options ...EntryOption) {
	GlobalLogger().Log(msg, options...)
}

// LogAny logs interface{}thing using the global logger.
func LogAny(v interface{}, options ...EntryOption) {
	GlobalLogger().LogAny(v, options...)
}

// Emergency logs an emergency message with additional EntryOptions.
func Emergency(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(EmergencySeverity).Log(msg, options...)
}

// Alert logs an alert message with additional EntryOptions.
func Alert(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(AlertSeverity).Log(msg, options...)
}

// Critical logs a critical message with additional EntryOptions.
func Critical(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(CriticalSeverity).Log(msg, options...)
}

// Error logs an error message with additional EntryOptions.
func Error(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(ErrorSeverity).Log(msg, options...)
}

// Warn logs a warning message with additional EntryOptions.
func Warn(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(WarnSeverity).Log(msg, options...)
}

// Info logs an info message with additional EntryOptions.
func Info(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(InfoSeverity).Log(msg, options...)
}

// Debug logs a debug message with additional EntryOptions.
func Debug(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(DebugSeverity).Log(msg, options...)
}
