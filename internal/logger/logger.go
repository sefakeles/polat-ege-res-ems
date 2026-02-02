package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"powerkonnekt/ems/internal/config"
)

// NewLogger creates and initializes a zap logger
func NewLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel // fallback to info level
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()

	// Set time encoder based on configuration
	switch cfg.TimeEncoder {
	case "epoch":
		encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	case "iso8601":
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // fallback
	}

	// Choose encoder based on encoding
	var encoder zapcore.Encoder
	if cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Setup output paths
	var outputs []zapcore.WriteSyncer
	for _, path := range cfg.OutputPaths {
		switch path {
		case "stdout":
			outputs = append(outputs, zapcore.AddSync(os.Stdout))
		case "stderr":
			outputs = append(outputs, zapcore.AddSync(os.Stderr))
		default:
			if file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
				outputs = append(outputs, zapcore.AddSync(file))
			}
		}
	}

	// Setup error output paths
	var errorOutputs []zapcore.WriteSyncer
	for _, path := range cfg.ErrorOutputPaths {
		switch path {
		case "stdout":
			errorOutputs = append(errorOutputs, zapcore.AddSync(os.Stdout))
		case "stderr":
			errorOutputs = append(errorOutputs, zapcore.AddSync(os.Stderr))
		default:
			if file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
				errorOutputs = append(errorOutputs, zapcore.AddSync(file))
			}
		}
	}

	// Create base core with combined outputs
	outputSyncer := zapcore.NewMultiWriteSyncer(outputs...)
	baseCore := zapcore.NewCore(encoder, outputSyncer, level)

	// Wrap with sampling core
	// Sample after the first 100 entries, then keep 1 of every 100 entries
	samplingCore := zapcore.NewSamplerWithOptions(
		baseCore,
		time.Second, // Sample period
		100,         // First N entries to keep
		100,         // Thereafter, keep 1 of every N
	)

	// Create logger with error output
	errorOutputSyncer := zapcore.NewMultiWriteSyncer(errorOutputs...)
	zapLogger := zap.New(samplingCore, zap.ErrorOutput(errorOutputSyncer))

	zapLogger.Info("Logger initialized",
		zap.String("level", cfg.Level),
		zap.String("encoding", cfg.Encoding),
		zap.String("timeEncoder", cfg.TimeEncoder),
		zap.Strings("outputPaths", cfg.OutputPaths),
		zap.Strings("errorOutputPaths", cfg.ErrorOutputPaths))

	return zapLogger, nil
}
