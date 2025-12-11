package logger

// noopLogger is a no-operation logger that implements the Logger interface
type noopLogger struct{}

func (n *noopLogger) Debug(msg string, fields ...Field)   {}
func (n *noopLogger) Info(msg string, fields ...Field)    {}
func (n *noopLogger) Warn(msg string, fields ...Field)    {}
func (n *noopLogger) Error(msg string, fields ...Field)   {}
func (n *noopLogger) Fatal(msg string, fields ...Field)   {}
func (n *noopLogger) Debugf(template string, args ...any) {}
func (n *noopLogger) Infof(template string, args ...any)  {}
func (n *noopLogger) Warnf(template string, args ...any)  {}
func (n *noopLogger) Errorf(template string, args ...any) {}
func (n *noopLogger) Fatalf(template string, args ...any) {}
func (n *noopLogger) With(fields ...Field) Logger         { return n }
func (n *noopLogger) Sync() error                         { return nil }
