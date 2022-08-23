package log

import (
	"context"
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

	SetGlobalLogger(NewDefaultLogger(&Writer{
		marshaler: JSONEntryMarshaler,
		writer:    os.Stdout,
	}))

	return GlobalLogger()
}

// SetGlobalLogger sets the global logger.
func SetGlobalLogger(l *Logger) {
	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()

	globalLogger = l
}

// Log log a message using the global builder.
func Log(msg string, options ...EntryOption) {
	GlobalLogger().Log(msg, options...)
}

// LogC log a message using the logger inside the context or the global logger.
func LogC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).LogC(ctx, msg, options...)
}

// LogAny logs interface{}thing using the global logger.
func LogAny(v interface{}) {
	GlobalLogger().LogAny(v)
}

// LogAnyC logs interface{}thing using the logger inside the context or the global logger.
func LogAnyC(ctx context.Context, v interface{}) {
	FromContextOrGlobal(ctx).LogAnyC(ctx, NewEntry(v))
}

// LogEntry logs an entry using the global logger.
func LogEntry(entry Entry) {
	GlobalLogger().LogEntry(entry)
}

// LogEntryC logs an entry with a context using the logger inside the context or the global logger.
func LogEntryC(ctx context.Context, entry Entry) {
	FromContextOrGlobal(ctx).LogEntryC(ctx, entry)
}

// WithLabels return a new logger from the global logger with additional labels.
func WithLabels(labels map[string]string) *Logger {
	return GlobalLogger().WithLabels(labels)
}

// WithLabelsC return a new logger from the logger inside the context or the global logger with additional labels.
func WithLabelsC(ctx context.Context, labels map[string]string) *Logger {
	return FromContextOrGlobal(ctx).WithLabels(labels)
}

// WithFields return a new logger from the global logger with additional fields.
func WithFields(fields map[string]interface{}) *Logger {
	return GlobalLogger().WithFields(fields)
}

// WithFieldsC return a new logger from the logger inside the context or the global logger with additional fields.
func WithFieldsC(ctx context.Context, fields map[string]interface{}) *Logger {
	return FromContextOrGlobal(ctx).WithFields(fields)
}

// WithSeverity return a new logger from the global logger with a default severity.
func WithSeverity(severity Severity) *Logger {
	return GlobalLogger().WithSeverity(severity)
}

// WithSeverityC return a new logger from the logger inside the context or the global logger with a default severity.
func WithSeverityC(ctx context.Context, severity Severity) *Logger {
	return FromContextOrGlobal(ctx).WithSeverity(severity)
}

// Emergency logs an emergency message with additional EntryOptions.
func Emergency(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(EmergencySeverity).LogC(context.Background(), msg, options...)
}

// Alert logs an alert message with additional EntryOptions.
func Alert(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(AlertSeverity).LogC(context.Background(), msg, options...)
}

// Critical logs a critical message with additional EntryOptions.
func Critical(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(CriticalSeverity).LogC(context.Background(), msg, options...)
}

// Error logs an error message with additional EntryOptions.
func Error(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(ErrorSeverity).LogC(context.Background(), msg, options...)
}

// Warn logs a warning message with additional EntryOptions.
func Warn(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(WarnSeverity).LogC(context.Background(), msg, options...)
}

// Info logs an info message with additional EntryOptions.
func Info(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(InfoSeverity).LogC(context.Background(), msg, options...)
}

// Debug logs a debug message with additional EntryOptions.
func Debug(msg string, options ...EntryOption) {
	GlobalLogger().WithSeverity(DebugSeverity).LogC(context.Background(), msg, options...)
}

// EmergencyC logs an emergency message with a context and additional EntryOptions.
func EmergencyC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(EmergencySeverity).LogC(ctx, msg, options...)
}

// AlertC logs an alert message with a context and  additional EntryOptions.
func AlertC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(AlertSeverity).LogC(ctx, msg, options...)
}

// CriticalC logs a critical message with a context and  additional EntryOptions.
func CriticalC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(CriticalSeverity).LogC(ctx, msg, options...)
}

// ErrorC logs an error message with a context and  additional EntryOptions.
func ErrorC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(ErrorSeverity).LogC(ctx, msg, options...)
}

// WarnC logs a warning message with a context and  additional EntryOptions.
func WarnC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(WarnSeverity).LogC(ctx, msg, options...)
}

// InfoC logs an info message with a context and  additional EntryOptions.
func InfoC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(InfoSeverity).LogC(ctx, msg, options...)
}

// DebugC logs a debug message with a context and  additional EntryOptions.
func DebugC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithSeverity(DebugSeverity).LogC(ctx, msg, options...)
}
