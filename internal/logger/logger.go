package logger

import (
	"fmt"
	"io"
)

// Possible Scrap Status  values.
const (
	TraceLevel uint32 = iota
	DebugLevel uint32 = iota
	InfoLevel  uint32 = iota
	WarnLevel  uint32 = iota
	ErrorLevel uint32 = iota
)

type color string

const (
	ColorReset  color = "\033[0m"
	ColorRed    color = "\033[31m"
	ColorGreen  color = "\033[32m"
	ColorYellow color = "\033[33m"
	ColorBlue   color = "\033[34m"
	ColorCyan   color = "\033[36m"
	ColorWhite  color = "\033[37m"
)

// Logger is a simple std out logger
type Logger struct {
	level  uint32
	writer io.Writer
}

func NewLogger(level uint32, writer io.Writer) Logger {
	return Logger{
		level:  level,
		writer: writer,
	}
}

func (l Logger) GetLevel() uint32 {
	return l.level
}

func (l Logger) Tracef(format string, args ...interface{}) {
	if l.GetLevel() <= TraceLevel {
		printfWithColor(l.writer, ColorCyan, format, args...)
	}
}

// Debugf logs a debug message.
func (l Logger) Debugf(format string, args ...interface{}) {
	if l.GetLevel() <= DebugLevel {
		printfWithColor(l.writer, ColorBlue, format, args...)
	}
}

// Infof logs a formatted info message.
func (l Logger) Infof(format string, args ...interface{}) {
	if l.GetLevel() <= InfoLevel {
		printfWithColor(l.writer, ColorGreen, format, args...)
	}
}

// Info logs a message.
func (l Logger) Info(str string) {
	if l.GetLevel() <=  InfoLevel {
		printWithColor(l.writer,ColorGreen,str)
	}
}

// Warnf logs a warning message.
func (l Logger) Warnf(format string, args ...interface{}) {
	if l.GetLevel() <= WarnLevel {
		printfWithColor(l.writer, ColorYellow, format, args...)
	}
}

// Errorf logs an error message.
func (l Logger) Errorf(format string, args ...interface{}) {
	if l.GetLevel() <= ErrorLevel {
		printfWithColor(l.writer, ColorRed, format, args...)
	}
}

func printfWithColor(w io.Writer, color color, format string, args ...interface{}) {
	fmt.Fprintf(w, string(color)+format+string(ColorReset), args...)
}

func printWithColor(w io.Writer, color color,str string){
	fmt.Fprint(w,string(color)+str+string(ColorReset))
}