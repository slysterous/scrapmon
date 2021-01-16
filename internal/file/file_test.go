package file_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/slysterous/scrapmon/internal/file"
	file_mock "github.com/slysterous/scrapmon/internal/file/mock"
	"github.com/slysterous/scrapmon/internal/scrapmon"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSaveFile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockWriter := file_mock.NewMockWriter(mockCtrl)
		mockPurger := file_mock.NewMockPurger(mockCtrl)
		manager := file.NewManager("./", mockWriter, mockPurger)

		f := scrapmon.ScrapedFile{
			Code: "test",
			Data: nil,
			Type: "png",
		}
		mockWriter.EXPECT().WriteFile("test.png", nil, os.FileMode(0644)).Return(nil).Times(1)
		err := manager.SaveFile(f)
		if err != nil {
			t.Errorf("unexpected error occured, err: %v", err)
		}
	})
	t.Run("Failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockWriter := file_mock.NewMockWriter(mockCtrl)
		mockPurger := file_mock.NewMockPurger(mockCtrl)
		manager := file.NewManager("./", mockWriter, mockPurger)

		f := scrapmon.ScrapedFile{
			Code: "test",
			Data: nil,
			Type: "png",
		}
		mockWriter.EXPECT().WriteFile("test.png", nil, os.FileMode(0644)).Return(errors.New("test error")).Times(1)
		err := manager.SaveFile(f)
		if err == nil {
			t.Error("expected error got nil")
		}
	})
}

func TestPurge(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockWriter := file_mock.NewMockWriter(mockCtrl)
		mockPurger := file_mock.NewMockPurger(mockCtrl)
		manager := file.NewManager("./", mockWriter, mockPurger)

		mockPurger.EXPECT().ReadDir("./").Return(nil, nil).Times(1)
		err := manager.Purge()
		if err != nil {
			t.Errorf("unexpected error occured, err: %v", err)
		}
	})
	t.Run("Failed to read dir", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockWriter := file_mock.NewMockWriter(mockCtrl)
		mockPurger := file_mock.NewMockPurger(mockCtrl)
		manager := file.NewManager("./", mockWriter, mockPurger)

		mockPurger.EXPECT().ReadDir("./").Return(nil, errors.New("error from read dir")).Times(1)

		err := manager.Purge()
		if err == nil {
			t.Error("expected error got nil")
		}
	})
	t.Run("Failed to remove all", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockWriter := file_mock.NewMockWriter(mockCtrl)
		mockPurger := file_mock.NewMockPurger(mockCtrl)
		manager := file.NewManager("./", mockWriter, mockPurger)

		r, _ := ioutil.ReadDir("./")

		mockPurger.EXPECT().ReadDir("./").Return(r, nil).Times(1)

		mockPurger.EXPECT().RemoveAll(path.Join([]string{"tmp", r[0].Name()}...)).
			Return(errors.New("error from remove all")).
			Times(1)

		err := manager.Purge()
		if err == nil {
			t.Error("expected error got nil")
		}
	})
}
