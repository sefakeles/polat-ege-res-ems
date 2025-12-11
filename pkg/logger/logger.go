package logger

// Initialize sets up the global logger
func Initialize(logLevel string) error {
	config := Config{
		Level:  logLevel,
		Format: "json",
	}

	zapLog, err := NewZapLogger(config)
	if err != nil {
		return err
	}

	SetGlobalLogger(zapLog)
	return nil
}

// InitializeWithConfig sets up the global logger with custom config
func InitializeWithConfig(config Config) error {
	zapLog, err := NewZapLogger(config)
	if err != nil {
		return err
	}

	SetGlobalLogger(zapLog)
	return nil
}
