package log

import (
	"context"
)

type contextKey int

// loggerKey is the key for logger.Logger values in Contexts.
// It is unexported; clients use logger.NewContext and logger.FromContext instead of using this key directly.
var loggerKey contextKey

// NewContext returns a new Context that carries value Logger.
func NewContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns the Logger value stored in ctx, if interface{}.
func FromContext(ctx context.Context) (*Logger, bool) {
	l, ok := ctx.Value(loggerKey).(*Logger)

	return l, ok
}

// FromContextOrGlobal returns the Logger value stored in ctx, or the global Logger if none are stored.
func FromContextOrGlobal(ctx context.Context) *Logger {
	l, ok := FromContext(ctx)
	if ok {
		return l
	}

	return GlobalLogger()
}

// LogC log a message using FromContextOrGlobal.
func LogC(ctx context.Context, msg string, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithContext(ctx).
		Log(msg, options...)
}

// LogAnyC logs anything using FromContextOrGlobal.
func LogAnyC(ctx context.Context, v interface{}, options ...EntryOption) {
	FromContextOrGlobal(ctx).WithContext(ctx).
		LogAny(v, options...)
}

// EmergencyC logs an emergency message with additional EntryOptions using FromContextOrGlobal.
func EmergencyC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(EmergencySeverity))
	LogC(ctx, msg, options...)
}

// AlertC logs an alert message with additional EntryOptions using FromContextOrGlobal.
func AlertC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(AlertSeverity))
	LogC(ctx, msg, options...)
}

// CriticalC logs a critical message with additional EntryOptions using FromContextOrGlobal.
func CriticalC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(CriticalSeverity))
	LogC(ctx, msg, options...)
}

// ErrorC logs an error message with additional EntryOptions using FromContextOrGlobal.
func ErrorC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(ErrorSeverity))
	LogC(ctx, msg, options...)
}

// WarnC logs a warning message with additional EntryOptions using FromContextOrGlobal.
func WarnC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(WarnSeverity))
	LogC(ctx, msg, options...)
}

// InfoC logs an info message with additional EntryOptions using FromContextOrGlobal.
func InfoC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(InfoSeverity))
	LogC(ctx, msg, options...)
}

// DebugC logs a debug message with additional EntryOptions using FromContextOrGlobal.
func DebugC(ctx context.Context, msg string, options ...EntryOption) {
	options = append(options, SeverityOpt(DebugSeverity))
	LogC(ctx, msg, options...)
}

// CtxWithOptions adds options to the logger in the context.
func CtxWithOptions(ctx context.Context, options ...EntryOption) context.Context {
	logger := FromContextOrGlobal(ctx)

	logger.options = append(logger.options, options...)

	return NewContext(ctx, logger)
}

// CtxWithSeverity set the default severity for the logger inside the context.
func CtxWithSeverity(ctx context.Context, severity Severity) context.Context {
	return CtxWithOptions(ctx, SeverityOpt(severity))
}

// CtxWithLabels adds labels for the logger inside the context.
func CtxWithLabels(ctx context.Context, labels map[string]string) context.Context {
	return CtxWithOptions(ctx, LabelsOpt(labels))
}

// CtxWithFields adds fields for the logger inside the context.
func CtxWithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	return CtxWithOptions(ctx, FieldsOpt(fields))
}

// CtxWithContext set the context for the logger inside the context.
func CtxWithContext(ctx context.Context) context.Context {
	return CtxWithOptions(ctx, ContextOpt(ctx))
}
