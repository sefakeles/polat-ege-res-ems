package logger

// Logger defines the logging interface
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Warnf(template string, args ...any)
	Errorf(template string, args ...any)
	Fatalf(template string, args ...any)

	With(fields ...Field) Logger
	Sync() error
}

// Field represents a structured logging field
type Field interface {
	Key() string
	Value() any
}

// Config represents logging configuration
type Config struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // "json" or "console"
}

// Global logger instance
var globalLogger Logger

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(l Logger) {
	globalLogger = l
}

// GetLogger returns the global logger instance
func GetLogger() Logger {
	if globalLogger == nil {
		// Fallback to a no-op logger if none is set
		globalLogger = &noopLogger{}
	}
	return globalLogger
}

// Convenience functions that use the global logger
func Debug(msg string, fields ...Field) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	GetLogger().Fatal(msg, fields...)
}

func Debugf(template string, args ...any) {
	GetLogger().Debugf(template, args...)
}

func Infof(template string, args ...any) {
	GetLogger().Infof(template, args...)
}

func Warnf(template string, args ...any) {
	GetLogger().Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	GetLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...any) {
	GetLogger().Fatalf(template, args...)
}

func With(fields ...Field) Logger {
	return GetLogger().With(fields...)
}

func Sync() error {
	return GetLogger().Sync()
}
