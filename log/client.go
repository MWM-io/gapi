package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/mwm-io/gapi/errors"
)

// Log is a simple client to improve the usability of the zap logger using GAPI
type Log struct {
	// Chosen function according to severity
	f      func(string, ...zap.Field)
	fields []zap.Field
}

// With returns a new Log with additional zap.Field
func (l *Log) With(fields ...zap.Field) *Log {
	l.fields = append(l.fields, fields...)
	return l
}

// LogMsg format and log a message
func (l *Log) LogMsg(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	if len(l.fields) > 0 {
		l.f(msg, l.fields...)
		return
	}

	l.f(msg)
}

// LogError take a GAPI error, format error message and log it
func (l *Log) LogError(err error) {
	castedErr, ok := err.(errors.Error)
	if !ok {
		l.LogMsg(err.Error())
		return
	}

	if originalErr := castedErr.Unwrap(); originalErr != nil {
		l.With(zap.String("original_error", castedErr.Unwrap().Error()))
		// TODO : if original_error is a GAPI error, log with call original callstack
	}
	l.With(
		zap.String("kind", castedErr.Kind()),
		zap.Strings("callstack", castedErr.Callstack()),
		zap.String("caller", castedErr.Caller()),
		zap.String("caller_name", castedErr.CallerName()),
	).LogMsg(castedErr.Error())
}

// Debug logs a debug message with additional zap.Field
func Debug(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Debug,
	}
}

// Info logs an info message with additional zap.Field
func Info(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Info,
	}
}

// Warn logs a warning message with additional zap.Field
func Warn(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Warn,
	}
}

// Error logs an error message with additional zap.Field
func Error(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Error,
	}
}

// Critical logs a critical message with additional zap.Field
func Critical(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Panic,
	}
}

// Alert logs an alert message with additional zap.Field
func Alert(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Panic,
	}
}

// Emergency logs an emergency message with additional zap.Field
func Emergency(ctx context.Context) *Log {
	return &Log{
		f: Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Fatal,
	}
}
