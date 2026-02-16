package log

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mwm-io/gapi/config"
)

var (
	globalLogger     *zap.Logger
	globalLoggerSync sync.Once

	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel(),
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}
}

func defaultConfig() *zap.Config {
	var encoding = "json"
	if config.IS_LOCAL {
		encoding = "console"
	}

	return &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:          encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}
}

// Logger return gapi global logger.
func Logger(ctx context.Context) *zap.Logger {
	if l, ok := FromContext(ctx); ok {
		return l
	}

	globalLoggerSync.Do(func() {
		if globalLogger != nil {
			return
		}

		var err error

		if globalLogger, err = defaultConfig().Build(); err != nil {
			panic(err)
		}
	})

	return globalLogger
}

//
// // Debug logs a debug message with additional zap.Field
// func Debug(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
// }
//
// // Info logs an info message with additional zap.Field
// func Info(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
// }
//
// // Warn logs a warning message with additional zap.Field
// func Warn(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
// }
//
// // Error logs an error message with additional zap.Field
// func Error(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
// }
//
// // Critical logs a critical message with additional zap.Field
// func Critical(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
// }
//
// // Alert logs an alert message with additional zap.Field
// func Alert(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
// }
//
// // Emergency logs an emergency message with additional zap.Field
// func Emergency(ctx context.Context, msg string, fields ...zap.Field) {
// 	Logger(ctx).WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
// }
