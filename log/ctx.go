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
