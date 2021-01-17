package log

//go:generate mockgen -destination mock/log.go -package log_mock . Logger

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	SetLevel(level uint32)
}

// Logger logs to stdout
type Log struct {
	logger Logger
}

func (l Log) SetLevel(level uint32) {
	l.logger.SetLevel(level)
}

// Debugf logs a debug message.
func (l Log) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args)
}

// Infof logs an info message.
func (l Log) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args)
}

// Warnf logs a warning message.
func (l Log) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args)
}

// Errorf logs an error message.
func (l Log) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args)
}
