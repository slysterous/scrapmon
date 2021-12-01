package scrapmon_mock

// https://github.com/golang/mock/pull/595

// Logger is a simple std out logger
type Logger struct {
}

func NewLogger() Logger {
	return Logger{
	}
}

func (l Logger) GetLevel() uint32 {
	return 1
}
// Tracef logs a debug message.
func (l Logger) Tracef(format string, args ...interface{}) {}

// Debugf logs a debug message.
func (l Logger) Debugf(format string, args ...interface{}) {}

// Infof logs a formatted info message.
func (l Logger) Infof(format string, args ...interface{}) {}

// Info logs a message.
func (l Logger) Info(str string) {}

// Warnf logs a warning message.
func (l Logger) Warnf(format string, args ...interface{}) {}

// Errorf logs an error message.
func (l Logger) Errorf(format string, args ...interface{}) {}
