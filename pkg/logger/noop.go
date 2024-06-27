package logger

type noopLogger struct{}

func (l *noopLogger) Debug(_ ...interface{}) {}

func (l *noopLogger) Debugf(_ string, _ ...interface{}) {}

func (l *noopLogger) Infof(_ string, _ ...interface{}) {}

func (l *noopLogger) Info(_ ...interface{}) {}

func (l *noopLogger) Errorf(_ string, _ ...interface{}) {}

func (l *noopLogger) Error(_ ...interface{}) {}

func (l *noopLogger) Warnf(_ string, _ ...interface{}) {}

func (l *noopLogger) Warn(_ ...interface{}) {
}
