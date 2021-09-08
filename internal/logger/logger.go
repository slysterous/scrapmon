package logger

import (
	"fmt"
	"io"
)

//go:generate mockgen -destination mock/log.go -package log_mock . Logger

// Possible Scrap Status  values.
const (
	DebugLevel uint32 = 1
	InfoLevel uint32  = 2
	WarnLevel  uint32 = 3
	ErrorLevel uint32 = 4
)

type color string

const (
	ColorReset color = "\033[0m"
	ColorRed color = "\033[31m"
	ColorGreen color = "\033[32m"
	ColorYellow color = "\033[33m"
	ColorBlue color = "\033[34m"
	ColorWhite color = "\033[37m"
)

// Logger is a simple std out logger
type Logger struct {
	level uint32
	writer io.Writer
}

func NewLogger(level uint32,writer io.Writer) Logger {
	return Logger{
		level: level,
		writer: writer,
	}
}

func (l Logger) GetLevel() uint32 {
	return l.level
}
// Debugf logs a debug message.
func (l Logger) Debugf(format string, args ...interface{}) {
	if l.GetLevel() <= 1 {
		printWithColor(l.writer,ColorBlue,format,args)
	}
}

// Infof logs an info message.
func (l Logger) Infof(format string, args ...interface{}) {
	if l.GetLevel() <= 2 {
		printWithColor(l.writer,ColorGreen,format,args)
	}
}

// Warnf logs a warning message.
func (l Logger) Warnf(format string, args ...interface{}) {
	if l.GetLevel() <= 3 {
		printWithColor(l.writer,ColorYellow,format,args)
	}
}

// Errorf logs an error message.
func (l Logger) Errorf(format string, args ...interface{}) {
	if l.GetLevel() <= 4 {
		printWithColor(l.writer,ColorRed,format,args)
	}
}

func printWithColor(w io.Writer,color color,format string, args ...interface{}) {
	fmt.Fprintf(w,string(color)+format+string(ColorReset),args)
}