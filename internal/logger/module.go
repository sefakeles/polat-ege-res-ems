package logger

import (
	"context"
	"errors"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"powerkonnekt/ems/internal/config"
)

// Module provides logging functionality to the Fx application
var Module = fx.Module("logger",
	fx.Provide(ProvideLogger),
	fx.Invoke(RegisterLifecycle),
)

// ProvideLogger creates and provides a zap.Logger instance
func ProvideLogger(cfg *config.Config) (*zap.Logger, error) {
	return NewLogger(cfg.Logging)
}

// RegisterLifecycle registers lifecycle hooks for the logger
func RegisterLifecycle(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := logger.Sync()
			// Ignore "sync /dev/stdout: inappropriate ioctl for device" error
			// This is a known issue on macOS/Linux when syncing stdout/stderr
			if err != nil && (errors.Is(err, os.ErrInvalid) ||
				err.Error() == "sync /dev/stdout: inappropriate ioctl for device" ||
				err.Error() == "sync /dev/stderr: inappropriate ioctl for device") {
				return nil
			}
			return err
		},
	})
}

// FxLogger is an Fx option that provides an fxevent.Logger
var FxLogger = fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
	zapLogger := fxevent.ZapLogger{Logger: logger}
	zapLogger.UseLogLevel(zapcore.DebugLevel)
	return &zapLogger
})
