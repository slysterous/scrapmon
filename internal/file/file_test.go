package file_test

import (
	"github.com/slysterous/scrapmon/internal/file"
	scrapmon"github.com/slysterous/scrapmon/internal/scrapmon"
	"os"
	"testing"
)

func TestNewManager(t *testing.T) {

}

func TestSaveFile(t *testing.T){
	t.Run("Success", func (t *testing.T){
		fm:=file.NewManager("./")
		scrapedFile:= scrapmon.ScrapedFile{
			Code: "file_test",
			Data: []byte{65, 66, 67, 226, 130, 172},
			Type: "png",
		}

		err:=fm.SaveFile(scrapedFile)
		if err!=nil{
			t.Errorf("unexpected error occured, err: %v",err)
		}

		if !fileExists("file_test.png"){
			t.Errorf("file was not created")
		}

		if nil != os.Remove("./file_test.png"){
			t.Errorf("could not remove file")
		}
	})
	t.Run("Failure", func (t *testing.T){
		fm:=file.NewManager("./")
		scrapedFile:= scrapmon.ScrapedFile{}

		err:=fm.SaveFile(scrapedFile)
		if err==nil{
			t.Errorf("expected error, got nil")
		}
	})
}

func TestPurge(t *testing.T) {
	t.Run("Success", func (t *testing.T){
		fm:=file.NewManager("./")
		scrapedFile:= scrapmon.ScrapedFile{
			Code: "file_test",
			Data: []byte{65, 66, 67, 226, 130, 172},
			Type: "png",
		}

		fm.Purge()
	})
	t.Run("Failure", func (t *testing.T){

	})
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}