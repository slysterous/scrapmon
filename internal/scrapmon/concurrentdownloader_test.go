package scrapmon_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	log_mock "github.com/slysterous/scrapmon/internal/log/mock"
	"github.com/slysterous/scrapmon/internal/scrapmon"
	scrapmon_mock "github.com/slysterous/scrapmon/internal/scrapmon/mock"
)

func TestDownloadFiles(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		wantFiles := []scrapmon.ScrapedFile{
			{
				Code: "a",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "b",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "c",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "d",
				Data: []byte{},
				Type: "png",
			},
		}

		filesToDownload := []scrapmon.Scrap{
			{
				ID:            1,
				RefCode:       "a",
				CodeCreatedAt: time.Now(),
				FileURI:       "",
				Status:        scrapmon.StatusPending,
			},
			{
				ID:            2,
				RefCode:       "b",
				CodeCreatedAt: time.Now(),
				FileURI:       "",
				Status:        scrapmon.StatusPending,
			},
			{
				ID:            3,
				RefCode:       "c",
				CodeCreatedAt: time.Now(),
				FileURI:       "",
				Status:        scrapmon.StatusPending,
			},
			{
				ID:            4,
				RefCode:       "d",
				CodeCreatedAt: time.Now(),
				FileURI:       "",
				Status:        scrapmon.StatusPending,
			},
		}

		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		//feed codes
		for _, file := range filesToDownload {
			mockScrapper.EXPECT().ScrapeByCode(file.RefCode, "png").Return(scrapmon.ScrapedFile{
				Code: file.RefCode,
				Data: []byte{},
				Type: "png",
			}, nil).Times(1)
			mockDM.EXPECT().UpdateScrapStatusByCode(file.RefCode, scrapmon.StatusOngoing).Return(nil).Times(1)
			pendingFiles <- file
		}

		filesC, _ := cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		var downloadedFiles []scrapmon.ScrapedFile

		for downloadedFile := range filesC {
			counter++
			downloadedFiles = append(downloadedFiles, downloadedFile)

			if counter == 4 {
				cancel()
				close(pendingFiles)
				close(produceMoreCodes)
			}
		}

		if !reflect.DeepEqual(wantFiles, downloadedFiles) {
			t.Errorf("expected: %v, got: %v", wantFiles, downloadedFiles)
		}

	})
	t.Run("Success with not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		scrap := scrapmon.Scrap{
			ID:            1,
			RefCode:       "a",
			CodeCreatedAt: time.Now(),
			FileURI:       "",
			Status:        scrapmon.StatusPending,
		}

		mockScrapper.EXPECT().ScrapeByCode(scrap.RefCode, "png").Return(scrapmon.ScrapedFile{
			Code: "a",
			Data: nil,
			Type: "png",
		}, nil).Times(1)
		mockDM.EXPECT().UpdateScrapStatusByCode("a", scrapmon.StatusNotFound).Return(nil).Times(1)
		mockLogger.EXPECT().Infof("File %s was not found, requesting a new one!\n", "a")

		pendingFiles <- scrap

		cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		for range produceMoreCodes {
			counter++
			if counter == 1 {
				cancel()
				close(pendingFiles)
				close(produceMoreCodes)
			} else {
				t.Error("Expected more codes to be asked exactly once.")
			}
		}
	})
	t.Run("Failed with not found and requested a new one", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		scrap := scrapmon.Scrap{
			ID:            1,
			RefCode:       "a",
			CodeCreatedAt: time.Now(),
			FileURI:       "",
			Status:        scrapmon.StatusPending,
		}

		mockScrapper.EXPECT().ScrapeByCode(scrap.RefCode, "png").Return(scrapmon.ScrapedFile{
			Code: "a",
			Data: nil,
			Type: "png",
		}, nil).Times(1)
		mockDM.EXPECT().UpdateScrapStatusByCode("a", scrapmon.StatusNotFound).Return(nil).Times(1)
		mockLogger.EXPECT().Infof("File %s was not found, requesting a new one!\n", "a")

		pendingFiles <- scrap

		cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		for range produceMoreCodes {
			counter++
			if counter == 1 {
				cancel()
				close(pendingFiles)
				close(produceMoreCodes)
			} else {
				t.Error("Expected more codes to be asked exactly once.")
			}
		}
	})
	t.Run("Failed with not found and failed to update state to notfound", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		scrap := scrapmon.Scrap{
			ID:            1,
			RefCode:       "a",
			CodeCreatedAt: time.Now(),
			FileURI:       "",
			Status:        scrapmon.StatusPending,
		}

		mockScrapper.EXPECT().ScrapeByCode(scrap.RefCode, "png").Return(scrapmon.ScrapedFile{
			Code: "a",
			Data: nil,
			Type: "png",
		}, nil).Times(1)
		mockDM.EXPECT().UpdateScrapStatusByCode("a", scrapmon.StatusNotFound).Return(errors.New("test error")).Times(1)
		mockLogger.EXPECT().Infof("File %s was not found, requesting a new one!\n", "a")

		pendingFiles <- scrap

		_, errs := cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		for range errs {
			counter++
			if counter == 1 {
				cancel()
				close(pendingFiles)
				close(produceMoreCodes)
			} else {
				t.Error("Expected errors to be thrown exactly once.")
			}
		}

		for range produceMoreCodes {
			t.Error("Expected more codes not to be asked.")
		}
	})
	t.Run("Failed to download", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		scrap := scrapmon.Scrap{
			ID:            1,
			RefCode:       "a",
			CodeCreatedAt: time.Now(),
			FileURI:       "",
			Status:        scrapmon.StatusPending,
		}

		mockScrapper.EXPECT().ScrapeByCode(scrap.RefCode, "png").Return(scrapmon.ScrapedFile{}, errors.New("test error")).Times(1)

		pendingFiles <- scrap

		_, errC := cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		for range errC {
			counter++
			if counter == 1 {
				cancel()
				close(pendingFiles)
				close(produceMoreCodes)
			} else {
				t.Error("Expected errors to be returned exactly once.")
			}
		}
	})
	t.Run("Failed to update state", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		scrap := scrapmon.Scrap{
			ID:            1,
			RefCode:       "a",
			CodeCreatedAt: time.Now(),
			FileURI:       "",
			Status:        scrapmon.StatusPending,
		}

		mockScrapper.EXPECT().ScrapeByCode(scrap.RefCode, "png").Return(scrapmon.ScrapedFile{Data: []byte{}}, nil).Times(1)
		mockDM.EXPECT().UpdateScrapStatusByCode("a", scrapmon.StatusOngoing).Return(errors.New("test error")).Times(1)
		pendingFiles <- scrap

		_, errC := cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		for err := range errC {
			counter++
			fmt.Printf("res: %v\n", err)
			if counter == 1 {
				cancel()
				close(pendingFiles)
				close(produceMoreCodes)
			} else {
				t.Error("Expected errors to be returned exactly once.")
			}
		}
	})
}

func TestSaveFiles(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		filesToSave := []scrapmon.ScrapedFile{
			{
				Code: "a",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "b",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "c",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "d",
				Data: []byte{},
				Type: "png",
			},
		}
		counter := 0
		mockLogger := log_mock.NewMockLogger(mockCtrl)
		cd := scrapmon.ConcurrentScrapper{
			Logger: mockLogger,
		}

		mockDM := scrapmon_mock.NewMockDatabaseManager(mockCtrl)
		mockFM := scrapmon_mock.NewMockFileManager(mockCtrl)

		mockStorage := scrapmon.Storage{
			Fm: mockFM,
			Dm: mockDM,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		filesToSaveCh := make(chan scrapmon.ScrapedFile, 5)
		//feed channel
		for _, file := range filesToSave {
			mockFM.EXPECT().SaveFile(file).Return(nil).Times(1)
			scrap := scrapmon.Scrap{
				ID:            int64(counter),
				RefCode:       file.Code,
				CodeCreatedAt: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				FileURI:       "SOMEWHERE" + file.Code + ".png",
				Status:        scrapmon.StatusSuccess,
			}
			mockDM.EXPECT().UpdateScrapByCode(scrap).Return(nil).Times(1)
			filesToSaveCh <- file
		}

		scraps, _ := cd.SaveFiles(mockStorage, ctx, filesToSaveCh)

		var savedFiles []scrapmon.Scrap

		for scrap := range scraps {
			counter++
			savedFiles = append(savedFiles, scrap)

			if counter == 4 {
				cancel()
				close(filesToSaveCh)
			}
		}

	})
	t.Run("Failed to save", func(t *testing.T) {

	})
	t.Run("Failed to update state", func(t *testing.T) {

	})
}
