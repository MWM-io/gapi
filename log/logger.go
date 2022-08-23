package log

import (
	"context"
)

// Logger is able to construct entries and write them using an EntryWriter.
type Logger struct {
	w       EntryWriter
	options []EntryOption
}

// NewLogger creates a new logger.
func NewLogger(w EntryWriter, options ...EntryOption) *Logger {
	return &Logger{
		w:       w,
		options: options,
	}
}

// NewDefaultLogger creates a new logger with the default stack of entry options:
// - timestamp set to time.Now()
// - stacktrace
// - tracing information from context (either opencensus or telemetry)
// If you don't need theses options, or just need some, you should use NewLogger and chose only the option you need.
func NewDefaultLogger(w EntryWriter) *Logger {
	return NewLogger(
		w,
		TimestampNowOpt(),
		DefaultStackTraceOpt(),
		OpencensusTraceOpt(),
		OpentelemetryTraceOpt(),
	)
}

// Log logs a message with additional EntryOptions.
func (l *Logger) Log(msg string, options ...EntryOption) {
	l.LogC(context.Background(), msg, options...)
}

// LogC logs a message with a context.Context and additional EntryOptions.
func (l *Logger) LogC(ctx context.Context, msg string, options ...EntryOption) {
	entry := Entry{Message: msg}
	for _, option := range options {
		option(ctx, &entry)
	}

	l.LogEntryC(ctx, entry)
}

// LogAny logs interface{}thing, converting the given value into an entry.
func (l *Logger) LogAny(v interface{}) {
	l.LogAnyC(context.Background(), v)
}

// LogAnyC logs interface{}thing with a context, converting the given value into an entry.
func (l *Logger) LogAnyC(ctx context.Context, v interface{}) {
	l.LogEntryC(ctx, NewEntry(v))
}

// LogEntry logs an entry.
func (l *Logger) LogEntry(entry Entry) {
	l.LogEntryC(context.Background(), entry)
}

// LogEntryC logs an entry with a context.
// It will set the default data using the logger's options and then write it.
func (l *Logger) LogEntryC(ctx context.Context, entry Entry) {
	for _, option := range l.options {
		option(ctx, &entry)
	}

	l.w.WriteEntry(entry)
}

// WithOptions returns a new logger with additional EntryOptions.
func (l *Logger) WithOptions(options ...EntryOption) *Logger {
	return &Logger{
		w:       l.w,
		options: append(l.options, options...),
	}
}

// WithSeverity returns a new logger with a default severity.
func (l *Logger) WithSeverity(severity Severity) *Logger {
	return l.WithOptions(SeverityOpt(severity))
}

// WithLabels returns a new logger with default labels.
func (l *Logger) WithLabels(labels map[string]string) *Logger {
	return l.WithOptions(LabelsOpt(labels))
}

// WithFields returns a new logger with default labels.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return l.WithOptions(FieldsOpt(fields))
}

// Emergency logs an emergency message with additional EntryOptions.
func (l *Logger) Emergency(msg string, options ...EntryOption) {
	l.WithSeverity(EmergencySeverity).LogC(context.Background(), msg, options...)
}

// Alert logs an alert message with additional EntryOptions.
func (l *Logger) Alert(msg string, options ...EntryOption) {
	l.WithSeverity(AlertSeverity).LogC(context.Background(), msg, options...)
}

// Critical logs a critical message with additional EntryOptions.
func (l *Logger) Critical(msg string, options ...EntryOption) {
	l.WithSeverity(CriticalSeverity).LogC(context.Background(), msg, options...)
}

// Error logs an error message with additional EntryOptions.
func (l *Logger) Error(msg string, options ...EntryOption) {
	l.WithSeverity(ErrorSeverity).LogC(context.Background(), msg, options...)
}

// Warn logs a warning message with additional EntryOptions.
func (l *Logger) Warn(msg string, options ...EntryOption) {
	l.WithSeverity(WarnSeverity).LogC(context.Background(), msg, options...)
}

// Info logs an info message with additional EntryOptions.
func (l *Logger) Info(msg string, options ...EntryOption) {
	l.WithSeverity(InfoSeverity).LogC(context.Background(), msg, options...)
}

// Debug logs a debug message with additional EntryOptions.
func (l *Logger) Debug(msg string, options ...EntryOption) {
	l.WithSeverity(DebugSeverity).LogC(context.Background(), msg, options...)
}

// EmergencyC logs an emergency message with a context and additional EntryOptions.
func (l *Logger) EmergencyC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(EmergencySeverity).LogC(ctx, msg, options...)
}

// AlertC logs an alert message with a context and  additional EntryOptions.
func (l *Logger) AlertC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(AlertSeverity).LogC(ctx, msg, options...)
}

// CriticalC logs a critical message with a context and  additional EntryOptions.
func (l *Logger) CriticalC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(CriticalSeverity).LogC(ctx, msg, options...)
}

// ErrorC logs an error message with a context and  additional EntryOptions.
func (l *Logger) ErrorC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(ErrorSeverity).LogC(ctx, msg, options...)
}

// WarnC logs a warning message with a context and  additional EntryOptions.
func (l *Logger) WarnC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(WarnSeverity).LogC(ctx, msg, options...)
}

// InfoC logs an info message with a context and  additional EntryOptions.
func (l *Logger) InfoC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(InfoSeverity).LogC(ctx, msg, options...)
}

// DebugC logs a debug message with a context and  additional EntryOptions.
func (l *Logger) DebugC(ctx context.Context, msg string, options ...EntryOption) {
	l.WithSeverity(DebugSeverity).LogC(ctx, msg, options...)
}
