package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

// loggerKey is the key for zap.Logger values in Contexts.
// It is unexported; clients use NewContext and FromContext instead of using this key directly.
var loggerKey contextKey = "gapi-logger"

// NewContext returns a new Context that carries value Logger.
func NewContext(ctx context.Context, l *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns the Logger value stored in Context.
func FromContext(ctx context.Context) (*zap.Logger, bool) {
	if ctx == nil {
		return nil, false
	}

	l, ok := ctx.Value(loggerKey).(*zap.Logger)

	return l, ok
}
