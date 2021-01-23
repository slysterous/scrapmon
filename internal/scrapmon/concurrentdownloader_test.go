package scrapmon_test

import (
	"context"
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

		wantFiles:=[]scrapmon.ScrapedFile{
			{
				Code: "a",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "b",
				Data:  []byte{},
				Type: "png",
			},
			{
				Code: "c",
				Data:  []byte{},
				Type: "png",
			},
			{
				Code: "d",
				Data:  []byte{},
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
			},nil).Times(1)
			mockDM.EXPECT().UpdateScrapStatusByCode(file.RefCode,scrapmon.StatusOngoing).Return(nil).Times(1)
			pendingFiles <- file
		}

		filesC, _ := cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		var downloadedFiles []scrapmon.ScrapedFile

		for downloadedFile:=range filesC {
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
	t.Run("Failure with produceMoreCodes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		wantFiles:=[]scrapmon.ScrapedFile{
			{
				Code: "a",
				Data: []byte{},
				Type: "png",
			},
			{
				Code: "b",
				Data:  []byte{},
				Type: "png",
			},
			{
				Code: "c",
				Data:  []byte{},
				Type: "png",
			},
			{
				Code: "d",
				Data:  []byte{},
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

		mockScrapper := scrapmon_mock.NewMockScrapper(mockCtrl)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pendingFiles := make(chan scrapmon.Scrap, 5)
		produceMoreCodes := make(chan struct{}, 5)

		scrap:=scrapmon.Scrap{
			ID:            1,
			RefCode:       "a",
			CodeCreatedAt: time.Now(),
			FileURI:       "",
			Status:        scrapmon.StatusPending,
		}

		mockScrapper.EXPECT().ScrapeByCode(scrap.RefCode, scrap.FileURI).Return(scrapmon.ScrapedFile{
			Code: "a",
			Data: nil,
			Type: "png",
		},nil).Times(1)
		mockDM.EXPECT().UpdateScrapStatusByCode("a",scrapmon.StatusNotFound).Return(nil).Times(1)
		pendingFiles <- scrap


		filesC, _ := cd.DownloadFiles(ctx, mockStorage, mockScrapper, pendingFiles, produceMoreCodes)

		var downloadedFiles []scrapmon.ScrapedFile

		for downloadedFile:=range filesC {
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

}
