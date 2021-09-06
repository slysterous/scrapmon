package logger_test

import (
	"bufio"
	"fmt"
	"github.com/golang/mock/gomock"
	logger "github.com/slysterous/scrapmon/internal/logger"
	"log"
	"os"
	"testing"
)

func TestDebugf(t *testing.T){
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		scanner, reader, writer,err := mockLogger(t) // turn this off when debugging or developing as you will miss output!
		defer resetLogger(reader, writer)
		if err !=nil{
			t.Errorf("unexpected error occured, err: %v",err)
		}
		lg:=logger.NewLogger(1)
		lg.Debugf("omg this is a test with param: %s","parameter")
		got := scanner.Text() // the last line written to the scanner
		msg := fmt.Sprintf("omg this is a test with param: %s","parameter")
		if got!=msg {
			t.Errorf("expected: %s, got: %s",msg,got)
		}
	})
}

func mockLogger(t *testing.T) (*bufio.Scanner, *os.File, *os.File,error)  {
	reader, writer, err := os.Pipe()
	if err != nil {
		return nil,nil,nil,err
	}
	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer,nil
}

func resetLogger(reader *os.File, writer *os.File) {
	err := reader.Close()
	if err != nil {
		fmt.Println("error closing reader was ", err)
	}
	if err = writer.Close(); err != nil {
		fmt.Println("error closing writer was ", err)
	}
	log.SetOutput(os.Stderr)
}