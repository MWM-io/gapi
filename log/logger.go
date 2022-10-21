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
// - tracing information from context (either opencensus or telemetry)
// If you don't need theses options, or just need some, you should use NewLogger and chose only the options you need.
func NewDefaultLogger(w EntryWriter) *Logger {
	return NewLogger(
		w,
		TimestampNowOpt(),
		OpencensusTraceOpt(),
		OpentelemetryTraceOpt(),
	)
}

// Log logs a message with additional EntryOptions.
func (l *Logger) Log(msg string, options ...EntryOption) {
	l.
		WithOptions(options...).
		LogEntry(NewEntry(msg))
}

// LogAny logs interface{}thing, converting the given value into an entry.
func (l *Logger) LogAny(v interface{}, options ...EntryOption) {
	l.
		WithOptions(options...).
		WithOptions(AnyOpt(v)).
		LogEntry(NewEntry(""))
}

// LogEntry logs an entry.
func (l *Logger) LogEntry(entry Entry) {
	MultiOpt(l.options...)(&entry)

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

// WithContext returns a new logger with a new context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return l.WithOptions(ContextOpt(ctx))
}

// Emergency logs an emergency message with additional EntryOptions.
func (l *Logger) Emergency(msg string, options ...EntryOption) {
	l.WithSeverity(EmergencySeverity).Log(msg, options...)
}

// Alert logs an alert message with additional EntryOptions.
func (l *Logger) Alert(msg string, options ...EntryOption) {
	l.WithSeverity(AlertSeverity).Log(msg, options...)
}

// Critical logs a critical message with additional EntryOptions.
func (l *Logger) Critical(msg string, options ...EntryOption) {
	l.WithSeverity(CriticalSeverity).Log(msg, options...)
}

// Error logs an error message with additional EntryOptions.
func (l *Logger) Error(msg string, options ...EntryOption) {
	l.WithSeverity(ErrorSeverity).Log(msg, options...)
}

// Warn logs a warning message with additional EntryOptions.
func (l *Logger) Warn(msg string, options ...EntryOption) {
	l.WithSeverity(WarnSeverity).Log(msg, options...)
}

// Info logs an info message with additional EntryOptions.
func (l *Logger) Info(msg string, options ...EntryOption) {
	l.WithSeverity(InfoSeverity).Log(msg, options...)
}

// Debug logs a debug message with additional EntryOptions.
func (l *Logger) Debug(msg string, options ...EntryOption) {
	l.WithSeverity(DebugSeverity).Log(msg, options...)
}
