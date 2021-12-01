package logger_test

import (
	"bytes"
	"fmt"
	logger "github.com/slysterous/scrapmon/internal/logger"
	"strings"
	"testing"
)

func TestTracef(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.TraceLevel, &buf)
		lg.Tracef("omg this is a test with param: %s", "parameter")

		got := buf.String()
		want := fmt.Sprintf(string(logger.ColorCyan)+"omg this is a test with param: %s"+string(logger.ColorReset), "parameter")

		if strings.Compare(want, got) != 0 {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})
	t.Run("Don't print due to log level", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.WarnLevel, &buf)
		lg.Tracef("omg this is a test with param: %s", "parameter")

		got := buf.String()
		if got != "" {
			t.Errorf("expected empty string, got: %s", got)
		}
	})
}

func TestDebugf(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.DebugLevel, &buf)
		lg.Debugf("omg this is a test with param: %s", "parameter")

		got := buf.String()
		want := fmt.Sprintf(string(logger.ColorBlue)+"omg this is a test with param: %s"+string(logger.ColorReset), "parameter")

		if strings.Compare(want, got) != 0 {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})
	t.Run("Don't print due to log level", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.WarnLevel, &buf)
		lg.Debugf("omg this is a test with param: %s", "parameter")

		got := buf.String()
		if got != "" {
			t.Errorf("expected empty string, got: %s", got)
		}
	})
}

func TestInfof(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.InfoLevel, &buf)
		lg.Infof("omg this is a test with param: %s", "parameter")

		got := buf.String()
		want := fmt.Sprintf(string(logger.ColorGreen)+"omg this is a test with param: %s"+string(logger.ColorReset), "parameter")

		if strings.Compare(want, got) != 0 {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})
	t.Run("Don't print due to log level", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.WarnLevel, &buf)
		lg.Infof("omg this is a test with param: %s", "parameter")

		got := buf.String()
		if got != "" {
			t.Errorf("expected empty string, got: %s", got)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.InfoLevel, &buf)
		lg.Info("omg this is a test")

		got := buf.String()
		want := string(logger.ColorGreen)+"omg this is a test"+string(logger.ColorReset)

		if strings.Compare(want, got) != 0 {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})
	t.Run("Don't print due to log level", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.ErrorLevel, &buf)
		lg.Info("omg this is a test")

		got := buf.String()
		if got != "" {
			t.Errorf("expected empty string, got: %s", got)
		}
	})
}

func TestWarnf(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.WarnLevel, &buf)
		lg.Warnf("omg this is a test with param: %s", "parameter")

		got := buf.String()
		want := fmt.Sprintf(string(logger.ColorYellow)+"omg this is a test with param: %s"+string(logger.ColorReset), "parameter")

		if strings.Compare(want, got) != 0 {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})
	t.Run("Don't print due to log level", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.ErrorLevel, &buf)
		lg.Warnf("omg this is a test with param: %s", "parameter")

		got := buf.String()
		if got != "" {
			t.Errorf("expected empty string, got: %s", got)
		}
	})
}

func TestErrorf(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(logger.ErrorLevel, &buf)
		lg.Errorf("omg this is a test with param: %s", "parameter")

		got := buf.String()
		want := fmt.Sprintf(string(logger.ColorRed)+"omg this is a test with param: %s"+string(logger.ColorReset), "parameter")

		if strings.Compare(want, got) != 0 {
			t.Errorf("expected: %s, got: %s", want, got)
		}
	})
	t.Run("Don't print due to log level", func(t *testing.T) {
		var buf bytes.Buffer

		lg := logger.NewLogger(5, &buf)
		lg.Errorf("omg this is a test with param: %s", "parameter")

		got := buf.String()
		if got != "" {
			t.Errorf("expected empty string, got: %s", got)
		}
	})
}
