package scrapmon_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/slysterous/scrapmon/internal/scrapmon"
	scrapmonmock "github.com/slysterous/scrapmon/internal/scrapmon/mock"
	"testing"
)

func TestStoragePurge(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm := scrapmonmock.NewMockFileManager(mockCtrl)
		mockDm := scrapmonmock.NewMockDatabaseManager(mockCtrl)

		storage := scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}

		mockFm.EXPECT().Purge().Return(nil).Times(1)
		mockDm.EXPECT().Purge().Return(nil).Times(1)

		err := storage.Purge()
		if err != nil {
			t.Errorf("unexpected error occured, err: %v", err)
		}
	})
	t.Run("Failure on Database", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm := scrapmonmock.NewMockFileManager(mockCtrl)
		mockDm := scrapmonmock.NewMockDatabaseManager(mockCtrl)

		storage := scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}

		mockDm.EXPECT().Purge().Return(errors.New("error from database manager")).Times(1)

		err := storage.Purge()
		if err == nil {
			t.Error("expected error got nil", err)
		}
		if err.Error() != "error from database manager" {
			t.Errorf("wanted: error from database manager, got: %v", err)
		}

	})
	t.Run("Failure on FileStorage", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm := scrapmonmock.NewMockFileManager(mockCtrl)
		mockDm := scrapmonmock.NewMockDatabaseManager(mockCtrl)

		storage := scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}

		mockDm.EXPECT().Purge().Return(nil).Times(1)
		mockFm.EXPECT().Purge().Return(errors.New("error from file manager"))

		err := storage.Purge()
		if err == nil {
			t.Fatal("expected error got nil", err)
		}
		if err.Error() != "error from file manager" {
			t.Errorf("wanted: error from file manager, got: %v", err)
		}

	})
}

func TestConcurrentCommandManagerPurgeCommand(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm := scrapmonmock.NewMockFileManager(mockCtrl)
		mockDm := scrapmonmock.NewMockDatabaseManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}
		mockLogger := scrapmonmock.NewLogger()
		mockScrapper := scrapmonmock.NewMockScrapper(mockCtrl)
		commandManager := scrapmon.ConcurrentCommandManager{
			Storage: mockStorage,
			CodeAuthority: scrapmon.ConcurrentCodeAuthority{
				Logger:   mockLogger,
				Scrapper: mockScrapper,
			},
		}

		mockDm.EXPECT().Purge().Return(nil).Times(1)
		mockFm.EXPECT().Purge().Return(nil).Times(1)

		err := commandManager.PurgeCommand()
		if err != nil {
			t.Errorf("expected nil, got : %v", err)
		}
	})
	t.Run("Failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockFm := scrapmonmock.NewMockFileManager(mockCtrl)
		mockDm := scrapmonmock.NewMockDatabaseManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFm,
			Dm: mockDm,
		}
		mockLogger := scrapmonmock.NewLogger()
		mockScrapper := scrapmonmock.NewMockScrapper(mockCtrl)
		commandManager := scrapmon.ConcurrentCommandManager{
			Storage: mockStorage,
			CodeAuthority: scrapmon.ConcurrentCodeAuthority{
				Logger:   mockLogger,
				Scrapper: mockScrapper,
			},
		}

		mockDm.EXPECT().Purge().Return(errors.New("test error")).Times(1)

		err := commandManager.PurgeCommand()
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
