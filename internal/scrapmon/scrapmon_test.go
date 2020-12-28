package scrapmon_test

import (
	"github.com/golang/mock/gomock"
	"github.com/slysterous/scrapmon/internal/scrapmon"
	scrapmon_mock "github.com/slysterous/scrapmon/internal/scrapmon/mock"
	"testing"
	"errors"
)

func TestStoragePurge(t *testing.T) {
	t.Run("Success",func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm:=scrapmon_mock.NewMockFileManager(mockCtrl)
		mockDm:=scrapmon_mock.NewMockDatabaseManager(mockCtrl)

		storage:= scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}

		mockFm.EXPECT().Purge().Return(nil).Times(1)
		mockDm.EXPECT().Purge().Return(nil).Times(1)

		err := storage.Purge()
		if err !=nil {
			t.Errorf("unexpected error occured, err: %v",err)
		}
	})
	t.Run("Failure on Database", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm:=scrapmon_mock.NewMockFileManager(mockCtrl)
		mockDm:=scrapmon_mock.NewMockDatabaseManager(mockCtrl)

		storage:= scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}

		mockDm.EXPECT().Purge().Return(errors.New("error from database manager")).Times(1)

		err := storage.Purge()
		if err ==nil {
			t.Error("expected error got nil",err)
		}
		if err.Error()!="error from database manager" {
			t.Errorf("wanted: error from database manager, got: %v",err)
		}

	})
	t.Run("Failure on FileStorage", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm:=scrapmon_mock.NewMockFileManager(mockCtrl)
		mockDm:=scrapmon_mock.NewMockDatabaseManager(mockCtrl)

		storage:= scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}

		mockDm.EXPECT().Purge().Return(nil).Times(1)
		mockFm.EXPECT().Purge().Return(errors.New("error from file manager"))

		err := storage.Purge()
		if err ==nil {
			t.Fatal("expected error got nil",err)
		}
		if err.Error()!="error from file manager" {
			t.Errorf("wanted: error from file manager, got: %v",err)
		}

	})
}