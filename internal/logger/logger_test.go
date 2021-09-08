package logger_test

import (
	"bufio"
	"bytes"
	"github.com/golang/mock/gomock"
	logger "github.com/slysterous/scrapmon/internal/logger"
	"testing"
)

func TestDebugf(t *testing.T){
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var buf bytes.Buffer
		wr:=bufio.NewWriter(&buf)

		lg:=logger.NewLogger(1,wr)
		lg.Debugf("omg this is a test with param: %s","parameter")

		got := buf.String()
		want := "omg this is a test with param: parameter"

		if got != want{
			t.Errorf("expected: %s, got: %s",want,got)
		}
	})
}