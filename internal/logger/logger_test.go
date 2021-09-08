package logger_test

import (
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	logger "github.com/slysterous/scrapmon/internal/logger"
	"strings"
	"testing"
)

func TestDebugf(t *testing.T){
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		var buf bytes.Buffer

		lg:=logger.NewLogger(1,&buf)
		lg.Debugf("omg this is a test with param: %s","parameter")

		got := buf.String()
		want := fmt.Sprintf(string(logger.ColorBlue)+"omg this is a test with param: %s"+string(logger.ColorReset),"parameter")

		if strings.Contains(got,want){
			t.Errorf("expected: %s, got: %s",want,got)
		}
	})
}