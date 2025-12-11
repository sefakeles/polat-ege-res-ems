package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements the Logger interface using zap
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new zap-based logger
func NewZapLogger(config Config) (Logger, error) {
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel // fallback to info level
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	// Choose encoder based on format
	var encoder zapcore.Encoder
	if config.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create base core
	baseCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)

	// Wrap with sampling core for production-like behavior
	// Sample after the first 100 entries, then keep 1 of every 100 entries
	samplingCore := zapcore.NewSamplerWithOptions(
		baseCore,
		time.Second, // Sample period
		100,         // First N entries to keep
		100,         // Thereafter, keep 1 of every N
	)

	// Create logger
	zapLog := zap.New(samplingCore)

	return &zapLogger{logger: zapLog}, nil
}

// convertFields converts our Field interface to zap.Field
func (l *zapLogger) convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key(), field.Value())
	}
	return zapFields
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, l.convertFields(fields)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, l.convertFields(fields)...)
}

func (l *zapLogger) Debugf(template string, args ...any) {
	l.logger.Debug(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Infof(template string, args ...any) {
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Warnf(template string, args ...any) {
	l.logger.Warn(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Errorf(template string, args ...any) {
	l.logger.Error(fmt.Sprintf(template, args...))
}

func (l *zapLogger) Fatalf(template string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(template, args...))
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(l.convertFields(fields)...),
	}
}

func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}
