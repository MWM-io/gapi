package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger     *zap.Logger
	globalLoggerMU   sync.Mutex
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
	return &zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// Logger return gapi global logger.
func Logger() *zap.Logger {
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

// SetLogger override global logger
func SetLogger(l *zap.Logger) {
	globalLoggerMU.Lock()
	defer globalLoggerMU.Unlock()

	globalLogger = l
}

// Debug logs a debug message with additional zap.Field
func Debug(msg string, fields ...zap.Field) {
	Logger().Debug(msg, fields...)
}

// Info logs an info message with additional zap.Field
func Info(msg string, fields ...zap.Field) {
	Logger().Info(msg, fields...)
}

// Warn logs a warning message with additional zap.Field
func Warn(msg string, fields ...zap.Field) {
	Logger().Warn(msg, fields...)
}

// Error logs an error message with additional zap.Field
func Error(msg string, fields ...zap.Field) {
	Logger().Error(msg, fields...)
}

// Critical logs a critical message with additional zap.Field
func Critical(msg string, fields ...zap.Field) {
	Logger().DPanic(msg, fields...)
}

// Alert logs an alert message with additional zap.Field
func Alert(msg string, fields ...zap.Field) {
	Logger().Panic(msg, fields...)
}

// Emergency logs an emergency message with additional zap.Field
func Emergency(msg string, fields ...zap.Field) {
	Logger().Fatal(msg, fields...)
}
