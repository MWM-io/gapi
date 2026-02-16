package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

// loggerKey is the key for zap.Logger values in Contexts.
// It is unexported; clients use NewContext and FromContext instead of using this key directly.
var loggerKey contextKey = "gapi-logger"

// loggerRefKey is the key for the mutable logger reference in Contexts.
var loggerRefKey contextKey = "gapi-logger-ref"

// loggerRef holds a mutable reference to a logger.
// It allows the Log middleware to see logger enrichments made by inner middlewares.
type loggerRef struct {
	Logger *zap.Logger
}

// NewContext returns a new Context that carries value Logger.
// If a logger ref exists in the context (created by NewRefContext), it also updates the ref
// so that the Log middleware can retrieve the latest enriched logger.
func NewContext(ctx context.Context, l *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if ref, ok := ctx.Value(loggerRefKey).(*loggerRef); ok {
		ref.Logger = l
	}

	return context.WithValue(ctx, loggerKey, l)
}

// NewRefContext creates a mutable logger reference in the context.
// This should be called by the Log middleware before calling NewContext,
// so that inner middlewares' calls to NewContext automatically update the ref.
func NewRefContext(ctx context.Context, l *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, loggerRefKey, &loggerRef{Logger: l})
}

// LatestLogger returns the logger from the mutable ref if available,
// otherwise falls back to the standard context logger.
// This is used by the Log middleware to get the logger enriched by inner middlewares.
func LatestLogger(ctx context.Context) *zap.Logger {
	if ref, ok := ctx.Value(loggerRefKey).(*loggerRef); ok && ref.Logger != nil {
		return ref.Logger
	}

	if l, ok := FromContext(ctx); ok {
		return l
	}

	return Logger(ctx)
}

// FromContext returns the Logger value stored in Context.
func FromContext(ctx context.Context) (*zap.Logger, bool) {
	if ctx == nil {
		return nil, false
	}

	l, ok := ctx.Value(loggerKey).(*zap.Logger)

	return l, ok
}
